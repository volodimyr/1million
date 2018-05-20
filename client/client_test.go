package main

import (
	"encoding/json"
	"testing"
)

type Event struct {
	Token string `json:"token"`
}

func TestMakeEvent(t *testing.T) {
	var event Event
	str := MakeEvent()
	if len(str) == 0 {
		t.Error("Event cannot be empty")
	}
	err := json.Unmarshal([]byte(str), &event)
	if err != nil {
		t.Error("Cannot unmarshall event", err)
	}
	if len(event.Token) == 0 {
		t.Error("Token cannot be empty")
	}
}
