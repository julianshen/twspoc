// sdk/client_test.go
package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"notification/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// createTestNotification creates a test notification
func createTestNotification() types.Notification {
	now := time.Now()
	expiry := now.Add(24 * time.Hour)
	return types.Notification{
		ID:        "test-id-1",
		Timestamp: now,
		Title:     "Test Notification",
		Message:   "This is a test notification",
		Priority:  "normal",
		Read:      false,
		Recipients: []types.Recipient{
			{Type: "user", ID: "user123"},
		},
		Labels:      []string{"test", "notification"},
		AppName:     "TestApp",
		AppIcon:     "https://example.com/icon.png",
		Expiry:      &expiry,
		GroupID:     "test-group",
		Attachments: []types.Attachment{},
		ActionButtons: []types.ActionButton{
			{Label: "View", Action: "view", URL: "https://example.com/view"},
		},
	}
}

// TestSendNotification tests the SendNotification method
func TestSendNotification(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/notifications", r.URL.Path)

		// Check content type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Decode the request body
		var notification types.Notification
		err := json.NewDecoder(r.Body).Decode(&notification)
		assert.NoError(t, err)

		// Check notification fields
		assert.Equal(t, "Test Notification", notification.Title)
		assert.Equal(t, "This is a test notification", notification.Message)
		assert.Equal(t, "normal", notification.Priority)
		assert.Equal(t, false, notification.Read)
		assert.Len(t, notification.Recipients, 1)
		assert.Equal(t, "user", notification.Recipients[0].Type)
		assert.Equal(t, "user123", notification.Recipients[0].ID)

		// Return success
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	// Create a client
	client := NewClient(server.URL)

	// Send a notification
	notification := createTestNotification()
	err := client.SendNotification(context.Background(), notification)
	assert.NoError(t, err)
}

// TestGetNotifications tests the GetNotifications method
func TestGetNotifications(t *testing.T) {
	// Create test notifications
	notifications := []types.Notification{
		createTestNotification(),
		createTestNotification(),
	}
	notifications[1].Title = "Another Notification"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/notifications", r.URL.Path)

		// Check query parameters
		assert.Equal(t, "user123", r.URL.Query().Get("userId"))

		// Return notifications
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notifications)
	}))
	defer server.Close()

	// Create a client
	client := NewClient(server.URL)

	// Get notifications
	result, err := client.GetNotifications(context.Background(), "user123")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Test Notification", result[0].Title)
	assert.Equal(t, "Another Notification", result[1].Title)
}

// TestMarkAsRead tests the MarkAsRead method
func TestMarkAsRead(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/notifications/123456789012345678901234/read", r.URL.Path)

		// Return success
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a client
	client := NewClient(server.URL)

	// Mark notification as read
	err := client.MarkAsRead(context.Background(), "123456789012345678901234")
	assert.NoError(t, err)
}

// TestDeleteNotification tests the DeleteNotification method
func TestDeleteNotification(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/notifications/123456789012345678901234", r.URL.Path)

		// Return success
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a client
	client := NewClient(server.URL)

	// Delete notification
	err := client.DeleteNotification(context.Background(), "123456789012345678901234")
	assert.NoError(t, err)
}

// TestSearchNotifications tests the SearchNotifications method
func TestSearchNotifications(t *testing.T) {
	// Create test notifications
	notifications := []types.Notification{
		createTestNotification(),
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/notifications/search", r.URL.Path)

		// Check content type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Decode the request body
		var params SearchParams
		err := json.NewDecoder(r.Body).Decode(&params)
		assert.NoError(t, err)

		// Check search parameters
		assert.Equal(t, "test", params.Keyword)
		assert.Equal(t, "user123", params.UserID)

		// Return notifications
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notifications)
	}))
	defer server.Close()

	// Create a client
	client := NewClient(server.URL)

	// Search notifications
	params := SearchParams{
		Keyword: "test",
		UserID:  "user123",
	}
	result, err := client.SearchNotifications(context.Background(), params)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Test Notification", result[0].Title)
}

// TestClientOptions tests the client options
func TestClientOptions(t *testing.T) {
	// Test WithTimeout
	client := NewClient("http://example.com", WithTimeout(5*time.Second))
	assert.Equal(t, 5*time.Second, client.httpClient.Timeout)

	// Test WithHTTPClient
	customClient := &http.Client{
		Timeout: 20 * time.Second,
	}
	client = NewClient("http://example.com", WithHTTPClient(customClient))
	assert.Equal(t, customClient, client.httpClient)
}
