package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	defaultSlug = "google-gdg-valencia-presents-devfest-valencia-2026-1"
	userAgent   = "devfestValencia2026-demo-example/1.0 (+https://github.com/aljesusg/devfestValencia2026)"
)

// apiBase is the GDG Community (Bevy) API origin. It is a variable so tests can point
// FetchEvent at an httptest server instead of the live API.
var apiBase = "https://gdg.community.dev"

// EventPayload is the subset we need from Bevy's event / event_slim JSON.
type EventPayload struct {
	ID       int             `json:"id"`
	Title    string          `json:"title"`
	URL      string          `json:"url"`
	Agenda   json.RawMessage `json:"agenda"`
	Timezone string          `json:"timezone_abbreviation"`
	Venue    string          `json:"venue_name"`
	Address  string          `json:"venue_address"`
	City     string          `json:"venue_city"`
}

// FetchEvent loads the event JSON from the public GDG Community API.
// Prefer event_slim by slug; fall back to /api/event/{id}/ when eventID > 0.
func FetchEvent(slug string, eventID int) (*EventPayload, error) {
	client := &http.Client{Timeout: 20 * time.Second}

	var url string
	switch {
	case eventID > 0:
		url = fmt.Sprintf("%s/api/event/%d/", apiBase, eventID)
		log.Printf("fetch: using event id=%d", eventID)
	default:
		if slug == "" {
			slug = defaultSlug
		}
		url = fmt.Sprintf("%s/api/event_slim/%s/", apiBase, slug)
		log.Printf("fetch: using event_slim slug=%q", slug)
	}

	log.Printf("fetch: GET %s", url)
	started := time.Now()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("fetch: build request failed: %v", err)
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("fetch: request failed after %s: %v", time.Since(started).Round(time.Millisecond), err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("fetch: read body failed after %s: %v", time.Since(started).Round(time.Millisecond), err)
		return nil, err
	}
	log.Printf("fetch: response status=%s bytes=%d elapsed=%s",
		resp.Status, len(body), time.Since(started).Round(time.Millisecond))

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("HTTP %s from %s", resp.Status, url)
		log.Printf("fetch: %v", err)
		return nil, err
	}

	var payload EventPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("fetch: decode json failed: %v", err)
		return nil, fmt.Errorf("decode event json: %w", err)
	}
	if len(payload.Agenda) == 0 {
		err := fmt.Errorf("event payload has empty agenda")
		log.Printf("fetch: %v (id=%d title=%q)", err, payload.ID, payload.Title)
		return nil, err
	}
	log.Printf("fetch: ok id=%d title=%q agenda_bytes=%d timezone=%q",
		payload.ID, payload.Title, len(payload.Agenda), payload.Timezone)
	return &payload, nil
}
