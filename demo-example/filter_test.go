package main

import "testing"

func sample() Schedule {
	return Schedule{Sessions: []Session{
		{ID: "s5", Day: "Friday", DayLocal: "Viernes", Time: "18:05", Speaker: "Alberto Gutiérrez", Title: "Taller Práctico: Construye tu primer servidor MCP paso a paso", Description: "Hands-on MCP workshop"},
		{ID: "s4", Day: "Friday", DayLocal: "Viernes", Time: "17:30", Speaker: "Jaime Blasco", Title: "Chef Gemini: Tu pinche de cocina en tiempo real con Flutter y Live AI"},
		{ID: "s10", Day: "Saturday", DayLocal: "Sábado", Time: "10:00", Speaker: "Javier Gallego", Title: "De Junior a Tech Lead"},
	}}
}

func TestFilterPerson(t *testing.T) {
	got := sample().Filter(ScheduleQuery{Person: "alberto"})
	if len(got) != 1 || got[0].ID != "s5" {
		t.Fatalf("got %#v", got)
	}
}

func TestFilterDayAndDescription(t *testing.T) {
	got := sample().Filter(ScheduleQuery{Day: "viernes", Description: "MCP"})
	if len(got) != 1 || got[0].ID != "s5" {
		t.Fatalf("got %#v", got)
	}
}

func TestFilterTime(t *testing.T) {
	got := sample().Filter(ScheduleQuery{Time: "6:05 PM"})
	if len(got) != 1 || got[0].ID != "s5" {
		t.Fatalf("got %#v", got)
	}
}
