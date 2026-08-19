package main

import (
	"encoding/json"
	"testing"
)

func TestFindDNSRecordByName(t *testing.T) {
	records := []Result{
		{ID: "1", Name: "example.com", Content: "1.1.1.1"},
		{ID: "2", Name: "www.example.com", Content: "2.2.2.2"},
	}

	found, ok := findDNSRecordByName(records, "EXAMPLE.com")
	if !ok {
		t.Fatal("expected record to be found case-insensitively")
	}
	if found.Content != "1.1.1.1" {
		t.Fatalf("expected 1.1.1.1, got %q", found.Content)
	}

	if _, ok := findDNSRecordByName(records, "missing.example.com"); ok {
		t.Fatal("expected missing record to be rejected")
	}
}

func TestFindDNSRecordByNameEmpty(t *testing.T) {
	if _, ok := findDNSRecordByName(nil, "example.com"); ok {
		t.Fatal("expected empty record set to be rejected")
	}
}

func TestParseResultListAcceptsSingleObject(t *testing.T) {
	payload := json.RawMessage(`{"id":"abc","name":"example.com","content":"78.110.169.189"}`)

	records, err := parseResultList(payload)
	if err != nil {
		t.Fatalf("expected object payload to parse, got error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected one record, got %d", len(records))
	}
	if records[0].Content != "78.110.169.189" {
		t.Fatalf("expected content 78.110.169.189, got %q", records[0].Content)
	}
}
