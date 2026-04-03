package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEventsEndpoint_ValidEvent(t *testing.T) {
	body := `{"type":"purchase","user_id":"u123","data":{"amount":49.99}}`
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	Events(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rec.Code)
	}

	var resp EventResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Status != "accepted" {
		t.Errorf("expected status 'accepted', got %q", resp.Status)
	}
}

func TestEventsEndpoint_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := httptest.NewRecorder()

	Events(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestEventsEndpoint_MissingType(t *testing.T) {
	body := `{"user_id":"u123"}`
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	Events(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestEventsEndpoint_UnknownType(t *testing.T) {
	body := `{"type":"unknown_event"}`
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	Events(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestEventsEndpoint_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString("{broken"))
	rec := httptest.NewRecorder()

	Events(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
