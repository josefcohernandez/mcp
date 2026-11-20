package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// AgendaRoot matches Bevy's agenda object (multiday days → items).
type AgendaRoot struct {
	Multiday bool        `json:"multiday"`
	Empty    bool        `json:"empty"`
	Days     []AgendaDay `json:"days"`
}

type AgendaDay struct {
	Title  string       `json:"title"`
	Agenda []AgendaItem `json:"agenda"`
}

type AgendaItem struct {
	Time         string `json:"time"`
	Activity     string `json:"activity"`
	Description  string `json:"description"`
	AudienceType string `json:"audience_type"`
}

// Schedule is the normalized shape used by the workshop mock + CLI output.
type Schedule struct {
	Event    string    `json:"event"`
	Location string    `json:"location"`
	Source   string    `json:"source"`
	Timezone string    `json:"timezone"`
	Sessions []Session `json:"sessions"`
}

type Session struct {
	ID           string `json:"id"`
	Day          string `json:"day"`
	DayLocal     string `json:"day_local"`
	Time         string `json:"time"`
	Track        string `json:"track,omitempty"`
	Title        string `json:"title"`
	Speaker      string `json:"speaker"`
	Room         string `json:"room,omitempty"`
	Description  string `json:"description,omitempty"`
	AudienceType string `json:"audience_type,omitempty"`
}

func ParseSchedule(payload *EventPayload) (Schedule, error) {
	log.Printf("parse: decoding agenda (%d bytes) for event id=%d", len(payload.Agenda), payload.ID)

	var root AgendaRoot
	if err := json.Unmarshal(payload.Agenda, &root); err != nil {
		// Some responses embed agenda as a JSON string.
		log.Printf("parse: agenda is not an object, trying embedded JSON string")
		var asString string
		if err2 := json.Unmarshal(payload.Agenda, &asString); err2 != nil {
			log.Printf("parse: agenda decode failed: %v", err)
			return Schedule{}, fmt.Errorf("agenda: %w", err)
		}
		if err := json.Unmarshal([]byte(asString), &root); err != nil {
			log.Printf("parse: agenda string decode failed: %v", err)
			return Schedule{}, fmt.Errorf("agenda string: %w", err)
		}
	}
	if root.Empty || len(root.Days) == 0 {
		err := fmt.Errorf("agenda is empty")
		log.Printf("parse: %v (empty=%v days=%d)", err, root.Empty, len(root.Days))
		return Schedule{}, err
	}

	loc := strings.TrimSpace(strings.Join([]string{
		payload.Venue,
		payload.Address,
		payload.City,
	}, ", "))
	loc = strings.ReplaceAll(loc, ", ,", ",")
	loc = strings.Trim(loc, ", ")

	out := Schedule{
		Event:    payload.Title,
		Location: loc,
		Source:   payload.URL,
		Timezone: payload.Timezone,
	}

	n := 1
	for _, day := range root.Days {
		dayEN := normalizeDay(day.Title)
		dayCount := 0
		for _, item := range day.Agenda {
			speaker, title := splitSpeakerTitle(item.Activity)
			out.Sessions = append(out.Sessions, Session{
				ID:           fmt.Sprintf("s%d", n),
				Day:          dayEN,
				DayLocal:     day.Title,
				Time:         strings.TrimSpace(item.Time),
				Title:        title,
				Speaker:      speaker,
				Description:  strings.TrimSpace(item.Description),
				AudienceType: item.AudienceType,
			})
			n++
			dayCount++
		}
		log.Printf("parse: day %q → %q sessions=%d", day.Title, dayEN, dayCount)
	}
	log.Printf("parse: ok event=%q location=%q total_sessions=%d", out.Event, out.Location, len(out.Sessions))
	return out, nil
}

func normalizeDay(title string) string {
	switch strings.ToLower(strings.TrimSpace(title)) {
	case "viernes", "friday":
		return "Friday"
	case "sábado", "sabado", "saturday":
		return "Saturday"
	default:
		return title
	}
}

// splitSpeakerTitle parses "Name - Talk title" (first " - " only).
func splitSpeakerTitle(activity string) (speaker, title string) {
	activity = strings.TrimSpace(activity)
	if activity == "" {
		return "", ""
	}
	const sep = " - "
	i := strings.Index(activity, sep)
	if i <= 0 {
		return "", activity
	}
	speaker = strings.TrimSpace(activity[:i])
	title = strings.TrimSpace(activity[i+len(sep):])
	// Guard against false positives like times or very short "speakers".
	if speaker == "" || title == "" || len(speaker) > 80 {
		return "", activity
	}
	return speaker, title
}
