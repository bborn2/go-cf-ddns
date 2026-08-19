package main

import "testing"

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
