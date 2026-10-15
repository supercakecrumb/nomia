package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteErrorRedactsInternalDetails(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(w, 500, "database host=private.invalid user=private")
	if strings.Contains(w.Body.String(), "private") {
		t.Fatal("internal detail leaked")
	}
	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["error"] != "Internal server error" {
		t.Fatal("unexpected public error")
	}
}
func TestWriteErrorEscapesValidationMessage(t *testing.T) {
	w := httptest.NewRecorder()
	message := "invalid \"quoted\" input"
	WriteError(w, 400, message)
	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["error"] != message {
		t.Fatal("validation detail changed")
	}
}
