package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"
)

// ScheduleQuery holds optional filters for the schedule tool.
type ScheduleQuery struct {
	Day         string
	Time        string
	Person      string
	Description string
}

func (s Schedule) Filter(q ScheduleQuery) []Session {
	log.Printf("filter: total=%d query day=%q time=%q person=%q description=%q",
		len(s.Sessions), q.Day, q.Time, q.Person, q.Description)

	var out []Session
	for _, sess := range s.Sessions {
		if !matchDay(sess, q.Day) {
			continue
		}
		if !matchTime(sess.Time, q.Time) {
			continue
		}
		if !matchPerson(sess.Speaker, q.Person) {
			continue
		}
		if !matchDescription(sess, q.Description) {
			continue
		}
		out = append(out, sess)
	}
	log.Printf("filter: matched=%d", len(out))
	return out
}

func matchDay(sess Session, day string) bool {
	day = strings.TrimSpace(day)
	if day == "" {
		return true
	}
	switch strings.ToLower(day) {
	case "friday", "viernes":
		return strings.EqualFold(sess.Day, "Friday") || strings.EqualFold(sess.DayLocal, "Viernes")
	case "saturday", "sabado", "sábado":
		return strings.EqualFold(sess.Day, "Saturday") || strings.EqualFold(sess.DayLocal, "Sábado")
	default:
		return strings.EqualFold(sess.Day, day) || strings.EqualFold(sess.DayLocal, day)
	}
}

func matchPerson(speaker, person string) bool {
	person = strings.TrimSpace(person)
	if person == "" {
		return true
	}
	return strings.Contains(strings.ToLower(speaker), strings.ToLower(person))
}

// matchDescription searches the query inside title + description (case-insensitive).
func matchDescription(sess Session, q string) bool {
	q = strings.TrimSpace(q)
	if q == "" {
		return true
	}
	hay := strings.ToLower(sess.Title + "\n" + sess.Description)
	return strings.Contains(hay, strings.ToLower(q))
}

var digits = regexp.MustCompile(`\d+`)

// matchTime is intentionally loose for live demos:
// "18:05", "6:05", "6:05 PM", "1805" can match "18:05" or "6:05 PM".
func matchTime(sessionTime, query string) bool {
	query = strings.TrimSpace(query)
	if query == "" {
		return true
	}
	a := normalizeClock(sessionTime)
	b := normalizeClock(query)
	if a == "" || b == "" {
		return strings.Contains(strings.ToLower(sessionTime), strings.ToLower(query))
	}
	return a == b
}

func normalizeClock(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return ""
	}
	pm := strings.Contains(s, "pm")
	am := strings.Contains(s, "am")

	nums := digits.FindAllString(s, -1)
	if len(nums) == 0 {
		return ""
	}
	h := atoi(nums[0])
	m := 0
	if len(nums) > 1 {
		m = atoi(nums[1])
	} else if len(nums[0]) >= 3 {
		// e.g. 1805
		v := nums[0]
		if len(v) == 3 {
			h = atoi(v[:1])
			m = atoi(v[1:])
		} else {
			h = atoi(v[:len(v)-2])
			m = atoi(v[len(v)-2:])
		}
	}

	if pm && h < 12 {
		h += 12
	}
	if am && h == 12 {
		h = 0
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return ""
	}
	return sprintfClock(h, m)
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		if !unicode.IsDigit(r) {
			break
		}
		n = n*10 + int(r-'0')
	}
	return n
}

func sprintfClock(h, m int) string {
	return fmt.Sprintf("%02d:%02d", h, m)
}
