package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type LogEntry struct {
	Severity  string `json:"severity"`
	ErrorType string `json:"error_type,omitempty"`
	Component string `json:"component,omitempty"`
	StoreID   string `json:"store_id,omitempty"`
	EventType string `json:"event_type,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Retryable bool   `json:"retryable,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp string `json:"timestamp"`
}

func LogError(
	severity string,
	errorType string,
	component string,
	storeID string,
	eventType string,
	requestID string,
	retryable bool,
	message string,
) {
	entry := LogEntry{
		Severity:  severity,
		ErrorType: errorType,
		Component: component,
		StoreID:   storeID,
		EventType: eventType,
		RequestID: requestID,
		Retryable: retryable,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
	}

	b, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintln(os.Stderr, `{"severity":"CRITICAL","message":"failed to marshal log entry"}`)
		return
	}

	if severity == "CRITICAL" {
		fmt.Fprintln(os.Stderr, string(b))
	} else {
		fmt.Fprintln(os.Stdout, string(b))
	}
}
