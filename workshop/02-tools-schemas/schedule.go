package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Session struct {
	ID          string `json:"id"`
	Day         string `json:"day"`
	DayLocal    string `json:"day_local,omitempty"`
	Time        string `json:"time"`
	Track       string `json:"track,omitempty"`
	Title       string `json:"title"`
	Speaker     string `json:"speaker"`
	Room        string `json:"room,omitempty"`
	Description string `json:"description,omitempty"`
}

type Schedule struct {
	Event    string    `json:"event"`
	Location string    `json:"location"`
	Source   string    `json:"source,omitempty"`
	Timezone string    `json:"timezone,omitempty"`
	Sessions []Session `json:"sessions"`
}

func LoadSchedule(path string) (*Schedule, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s Schedule
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Schedule) Filter(day, track string) []Session {
	var out []Session
	for _, sess := range s.Sessions {
		if day != "" && !strings.EqualFold(sess.Day, day) && !strings.EqualFold(sess.DayLocal, day) {
			continue
		}
		if track != "" && !strings.EqualFold(sess.Track, track) {
			continue
		}
		out = append(out, sess)
	}
	return out
}

func (s *Schedule) ByID(id string) (Session, error) {
	for _, sess := range s.Sessions {
		if sess.ID == id {
			return sess, nil
		}
	}
	return Session{}, fmt.Errorf("session %q not found", id)
}

func (s *Schedule) BySpeaker(name string) []Session {
	needle := strings.ToLower(strings.TrimSpace(name))
	var out []Session
	for _, sess := range s.Sessions {
		if strings.Contains(strings.ToLower(sess.Speaker), needle) {
			out = append(out, sess)
		}
	}
	return out
}
