// api/handler_test.go
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"notification/notificationrepo"
	"notification/types"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function to create a test notification
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

// TestCreateNotification tests creating a notification using a mock repository and gin.
func TestCreateNotification(t *testing.T) {
	mockRepo := new(notificationrepo.MockNotificationRepository)
	router := gin.New()
	router.POST("/api/notifications", GinCreateNotificationHandler(mockRepo))

	notification := createTestNotification()
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*types.Notification")).Return(nil)

	body, _ := json.Marshal(notification)
	req, _ := http.NewRequest("POST", "/api/notifications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockRepo.AssertExpectations(t)
}

// TestGetNotifications tests listing notifications for a user using a mock repository and gin.
func TestGetNotifications(t *testing.T) {
	mockRepo := new(notificationrepo.MockNotificationRepository)
	router := gin.New()
	router.GET("/api/notifications", GinListNotificationsHandler(mockRepo))

	n1 := createTestNotification()
	n2 := createTestNotification()
	notifications := []*types.Notification{&n1, &n2}
	mockRepo.On("ListByUser", mock.Anything, "user123").Return(notifications, nil)

	req, _ := http.NewRequest("GET", "/api/notifications?userId=user123", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []types.Notification
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	mockRepo.AssertExpectations(t)
}

// func TestMarkNotificationAsRead(t *testing.T) {
// 	mockRepo := new(MockNotificationRepo)
// 	handler := MakeHandler(mockRepo)

// 	notificationID := primitive.NewObjectID().Hex()
// 	mockRepo.On("MarkRead", mock.Anything, notificationID).Return(nil)

// 	// Create a request
// 	req, _ := http.NewRequest("POST", "/api/notifications/"+notificationID+"/read", nil)
// 	rr := httptest.NewRecorder()

// 	// Call the handler
// 	handler(rr, req)

// 	// Check the status code
// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	mockRepo.AssertExpectations(t)
// }

// func TestDeleteNotification(t *testing.T) {
// 	mockRepo := new(MockNotificationRepo)
// 	handler := MakeHandler(mockRepo)

// 	notificationID := primitive.NewObjectID().Hex()
// 	userID := "user123"
// 	mockRepo.On("Delete", mock.Anything, notificationID, userID).Return(nil)

// 	// Create a request
// 	req, _ := http.NewRequest("DELETE", "/api/notifications/"+notificationID+"?userId="+userID, nil)
// 	rr := httptest.NewRecorder()

// 	// Call the handler
// 	handler(rr, req)

// 	// Check the status code
// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	mockRepo.AssertExpectations(t)
// }

// func TestSearchNotifications(t *testing.T) {
// 	mockRepo := new(MockNotificationRepo)
// 	handler := SearchHandler(mockRepo)

// 	notifications := []types.Notification{createTestNotification()}
// 	searchParams := struct {
// 		Keyword   string   `json:"keyword"`
// 		Title     string   `json:"title"`
// 		Message   string   `json:"message"`
// 		Labels    []string `json:"labels"`
// 		AppName   string   `json:"appName"`
// 		StartDate string   `json:"startDate"`
// 		EndDate   string   `json:"endDate"`
// 		UserID    string   `json:"userId"`
// 	}{
// 		Keyword: "test",
// 		UserID:  "user123",
// 	}

// 	mockRepo.On("Search", mock.Anything, mock.Anything).Return(notifications, nil)

// 	// Create a request
// 	body, _ := json.Marshal(searchParams)
// 	req, _ := http.NewRequest("POST", "/api/notifications/search", bytes.NewBuffer(body))
// 	rr := httptest.NewRecorder()

// 	// Call the handler
// 	handler(rr, req)

// 	// Check the status code
// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	// Parse the response
// 	var response []types.Notification
// 	err := json.Unmarshal(rr.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Len(t, response, 1)
// 	mockRepo.AssertExpectations(t)
// }

// Test filtering notifications
func TestFilterNotifications(t *testing.T) {
	// Create test notifications
	now := time.Now()
	n1 := types.Notification{
		ID:        "test-id-2",
		Timestamp: now,
		Title:     "Test 1",
		Priority:  "high",
		Read:      true,
		GroupID:   "group1",
		Labels:    []string{"label1", "label2"},
	}
	n2 := types.Notification{
		ID:        "test-id-3",
		Timestamp: now.Add(-1 * time.Hour),
		Title:     "Test 2",
		Priority:  "low",
		Read:      false,
		GroupID:   "group2",
		Labels:    []string{"label3"},
	}
	// Create the array with n2 first to test the sorting
	notifications := []types.Notification{n2, n1}

	// Test filtering by read status
	t.Run("FilterByReadStatus", func(t *testing.T) {
		query := map[string][]string{
			"read": {"true"},
		}
		filtered := applyFilters(notifications, query)
		assert.Len(t, filtered, 1)
		assert.Equal(t, "Test 1", filtered[0].Title)
	})

	// Test filtering by priority
	t.Run("FilterByPriority", func(t *testing.T) {
		query := map[string][]string{
			"priority": {"high"},
		}
		filtered := applyFilters(notifications, query)
		assert.Len(t, filtered, 1)
		assert.Equal(t, "Test 1", filtered[0].Title)
	})

	// Test filtering by labels
	t.Run("FilterByLabels", func(t *testing.T) {
		query := map[string][]string{
			"labels": {"label1"},
		}
		filtered := applyFilters(notifications, query)
		assert.Len(t, filtered, 1)
		assert.Equal(t, "Test 1", filtered[0].Title)
	})

	// Test filtering by groupId
	t.Run("FilterByGroupId", func(t *testing.T) {
		query := map[string][]string{
			"groupId": {"group2"},
		}
		filtered := applyFilters(notifications, query)
		assert.Len(t, filtered, 1)
		assert.Equal(t, "Test 2", filtered[0].Title)
	})

	// Test sorting by timestamp (descending order - newest first)
	t.Run("SortByTimestamp", func(t *testing.T) {
		query := map[string][]string{
			"sort":  {"timestamp"},
			"order": {"desc"}, // Explicitly set descending order to get newest first
		}
		sorted := applyFilters(notifications, query)
		assert.Equal(t, "Test 1", sorted[0].Title) // n1 has a newer timestamp
	})

	// Test sorting by priority
	t.Run("SortByPriority", func(t *testing.T) {
		query := map[string][]string{
			"sort":  {"priority"},
			"order": {"desc"},
		}
		sorted := applyFilters(notifications, query)
		assert.Equal(t, "Test 1", sorted[0].Title) // high priority first
	})
}
