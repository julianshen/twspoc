// sdk/sse_test.go
package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSubscribeToNotifications tests the SubscribeToNotifications method
func TestSubscribeToNotifications(t *testing.T) {
	// Create test notifications
	notification1 := createTestNotification()
	notification2 := createTestNotification()
	notification2.Title = "Another Notification"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/notifications/subscribe", r.URL.Path)

		// Check query parameters
		assert.Equal(t, "user123", r.URL.Query().Get("userId"))

		// Check headers
		assert.Equal(t, "text/event-stream", r.Header.Get("Accept"))
		assert.Equal(t, "no-cache", r.Header.Get("Cache-Control"))
		assert.Equal(t, "keep-alive", r.Header.Get("Connection"))

		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		// Flush headers
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		// Send notifications
		notificationJSON1, _ := json.Marshal(notification1)
		notificationJSON2, _ := json.Marshal(notification2)

		// Send first notification
		fmt.Fprintf(w, "data: %s\n\n", notificationJSON1)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		// Wait a bit
		time.Sleep(100 * time.Millisecond)

		// Send second notification
		fmt.Fprintf(w, "data: %s\n\n", notificationJSON2)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		// Keep the connection open
		select {
		case <-r.Context().Done():
			return
		case <-time.After(500 * time.Millisecond):
			return
		}
	}))
	defer server.Close()

	// Create a client
	client := NewClient(server.URL)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to notifications
	eventCh, err := client.SubscribeToNotifications(ctx, "user123")
	assert.NoError(t, err)

	// Receive first notification
	event1 := <-eventCh
	assert.NoError(t, event1.Error)
	assert.Equal(t, "Test Notification", event1.Notification.Title)

	// Receive second notification
	event2 := <-eventCh
	assert.NoError(t, event2.Error)
	assert.Equal(t, "Another Notification", event2.Notification.Title)
}

// TestSubscribeToNotificationsError tests error handling in SubscribeToNotifications
func TestSubscribeToNotificationsError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Create a client
	client := NewClient(server.URL)

	// Subscribe to notifications
	_, err := client.SubscribeToNotifications(context.Background(), "user123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}

// TestSubscribeToNotificationsParseError tests error handling for malformed SSE data
func TestSubscribeToNotificationsParseError(t *testing.T) {
	// Create a test server that sends malformed data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		// Flush headers
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		// Send malformed data
		fmt.Fprintf(w, "data: {\"invalid\":\"json\",}\n\n")
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		// Keep the connection open
		select {
		case <-r.Context().Done():
			return
		case <-time.After(500 * time.Millisecond):
			return
		}
	}))
	defer server.Close()

	// Create a client
	client := NewClient(server.URL)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to notifications
	eventCh, err := client.SubscribeToNotifications(ctx, "user123")
	assert.NoError(t, err)

	// Receive error event
	event := <-eventCh
	assert.Error(t, event.Error)
	assert.Contains(t, event.Error.Error(), "failed to parse notification")
}

// TestSubscribeToNotificationsCancel tests cancellation of the SSE stream
func TestSubscribeToNotificationsCancel(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		// Flush headers
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		// Keep the connection open
		<-r.Context().Done()
	}))
	defer server.Close()

	// Create a client
	client := NewClient(server.URL)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Subscribe to notifications
	eventCh, err := client.SubscribeToNotifications(ctx, "user123")
	assert.NoError(t, err)

	// Cancel the context
	cancel()

	// Wait for the channel to close
	time.Sleep(100 * time.Millisecond)

	// Try to receive from the channel (should be closed)
	_, ok := <-eventCh
	assert.False(t, ok, "Channel should be closed")
}
