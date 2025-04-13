// api/handler.go
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notification/notificationrepo"
	"notification/types"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var subscribers = struct {
	sync.Mutex
	clients map[string][]chan types.Notification
}{clients: make(map[string][]chan types.Notification)}

// --- Gin Handlers ---

func GinCreateNotificationHandler(repo notificationrepo.NotificationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var n types.Notification
		if err := c.ShouldBindJSON(&n); err != nil {
			log.Printf("[GinCreateNotificationHandler] Invalid request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[GinCreateNotificationHandler] Creating notification: %+v", n)
		if err := repo.Create(c.Request.Context(), &n); err != nil {
			log.Printf("[GinCreateNotificationHandler] Create error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[GinCreateNotificationHandler] Notification created successfully: %+v", n)
		Broadcast(n)
		c.Status(http.StatusCreated)
	}
}

func GinListNotificationsHandler(repo notificationrepo.NotificationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("userId")
		if userId == "" {
			log.Printf("[GinListNotificationsHandler] Missing userId")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing userId"})
			return
		}
		log.Printf("[GinListNotificationsHandler] Listing notifications for userId: %s", userId)
		notifs, err := repo.ListByUser(c.Request.Context(), userId)
		if err != nil {
			log.Printf("[GinListNotificationsHandler] ListByUser error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Convert []*types.Notification to []types.Notification for applyFilters
		plainNotifs := make([]types.Notification, 0, len(notifs))
		for _, n := range notifs {
			plainNotifs = append(plainNotifs, *n)
		}
		filteredNotifs := applyFilters(plainNotifs, c.Request.URL.Query())
		log.Printf("[GinListNotificationsHandler] Returning %d notifications for userId: %s", len(filteredNotifs), userId)
		c.JSON(http.StatusOK, filteredNotifs)
	}
}

func GinMarkAsReadHandler(repo notificationrepo.NotificationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			log.Printf("[GinMarkAsReadHandler] Missing notification id")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing notification id"})
			return
		}
		log.Printf("[GinMarkAsReadHandler] Marking notification as read: %s", id)
		n, err := repo.Get(c.Request.Context(), id)
		if err != nil {
			log.Printf("[GinMarkAsReadHandler] Get error: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}
		n.Read = true
		if err := repo.Update(c.Request.Context(), n); err != nil {
			log.Printf("[GinMarkAsReadHandler] Update error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[GinMarkAsReadHandler] Notification marked as read: %s", id)
		c.Status(http.StatusOK)
	}
}

func GinDeleteNotificationHandler(repo notificationrepo.NotificationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			log.Printf("[GinDeleteNotificationHandler] Missing notification id")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing notification id"})
			return
		}
		log.Printf("[GinDeleteNotificationHandler] Deleting notification: %s", id)
		if err := repo.Delete(c.Request.Context(), id); err != nil {
			log.Printf("[GinDeleteNotificationHandler] Delete error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[GinDeleteNotificationHandler] Notification deleted: %s", id)
		c.Status(http.StatusOK)
	}
}

/* Search handler is not implemented in the new NotificationRepository interface.
func GinSearchHandler(repo notificationrepo.NotificationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "search not implemented"})
	}
}
*/

func GinSSEHandler(repo notificationrepo.NotificationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("userId")
		if userID == "" {
			log.Printf("[GinSSEHandler] Missing userId")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing userId"})
			return
		}

		log.Printf("[GinSSEHandler] New SSE subscription for userId: %s", userID)
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			log.Printf("[GinSSEHandler] Streaming unsupported")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming unsupported"})
			return
		}
		flusher.Flush()

		notifCh := make(chan types.Notification)
		subscribers.Lock()
		subscribers.clients[userID] = append(subscribers.clients[userID], notifCh)
		subscribers.Unlock()

		ctx := c.Request.Context()
		for {
			select {
			case <-ctx.Done():
				log.Printf("[GinSSEHandler] SSE context done for userId: %s", userID)
				return
			case notif := <-notifCh:
				jsonData, _ := json.Marshal(notif)
				log.Printf("[GinSSEHandler] Sending notification to userId %s: %+v", userID, notif)
				fmt.Fprintf(c.Writer, "data: %s\n\n", jsonData)
				flusher.Flush()
			}
		}
	}
}

// --- Utility and Broadcast ---

func applyFilters(notifs []types.Notification, query map[string][]string) []types.Notification {
	result := notifs

	// Filter by read status
	if readParam, ok := query["read"]; ok && len(readParam) > 0 {
		readStatus := readParam[0] == "true"
		filtered := []types.Notification{}
		for _, n := range result {
			if n.Read == readStatus {
				filtered = append(filtered, n)
			}
		}
		result = filtered
	}

	// Filter by priority
	if priorityParam, ok := query["priority"]; ok && len(priorityParam) > 0 {
		priority := priorityParam[0]
		filtered := []types.Notification{}
		for _, n := range result {
			if n.Priority == priority {
				filtered = append(filtered, n)
			}
		}
		result = filtered
	}

	// Filter by labels
	if labelsParam, ok := query["labels"]; ok && len(labelsParam) > 0 {
		labels := strings.Split(labelsParam[0], ",")
		filtered := []types.Notification{}
		for _, n := range result {
			// Check if notification has any of the specified labels
			hasLabel := false
			for _, label := range labels {
				for _, nLabel := range n.Labels {
					if nLabel == label {
						hasLabel = true
						break
					}
				}
				if hasLabel {
					break
				}
			}
			if hasLabel {
				filtered = append(filtered, n)
			}
		}
		result = filtered
	}

	// Filter by groupId
	if groupIdParam, ok := query["groupId"]; ok && len(groupIdParam) > 0 {
		groupId := groupIdParam[0]
		filtered := []types.Notification{}
		for _, n := range result {
			if n.GroupID == groupId {
				filtered = append(filtered, n)
			}
		}
		result = filtered
	}

	// Filter by timestamp (after a certain time)
	if timestampParam, ok := query["timestamp"]; ok && len(timestampParam) > 0 {
		timestamp, err := time.Parse(time.RFC3339, timestampParam[0])
		if err == nil {
			filtered := []types.Notification{}
			for _, n := range result {
				if n.Timestamp.After(timestamp) || n.Timestamp.Equal(timestamp) {
					filtered = append(filtered, n)
				}
			}
			result = filtered
		}
	}

	// Apply sorting
	if sortParam, ok := query["sort"]; ok && len(sortParam) > 0 {
		sortField := sortParam[0]
		sortOrder := "asc"

		// Check if sort order is specified
		if orderParam, ok := query["order"]; ok && len(orderParam) > 0 {
			if orderParam[0] == "desc" {
				sortOrder = "desc"
			}
		}

		switch sortField {
		case "timestamp":
			sortByTimestamp(result, sortOrder == "desc")
		case "priority":
			sortByPriority(result, sortOrder == "desc")
		case "read":
			sortByReadStatus(result, sortOrder == "desc")
		case "groupId":
			sortByGroupID(result, sortOrder == "desc")
		}
	} else {
		// Default sort by timestamp (newest first)
		sortByTimestamp(result, true)
	}

	return result
}

func sortByTimestamp(notifs []types.Notification, descending bool) {
	if descending {
		sort.Slice(notifs, func(i, j int) bool {
			return notifs[i].Timestamp.After(notifs[j].Timestamp)
		})
	} else {
		sort.Slice(notifs, func(i, j int) bool {
			return notifs[i].Timestamp.Before(notifs[j].Timestamp)
		})
	}
}

func sortByPriority(notifs []types.Notification, descending bool) {
	priorityMap := map[string]int{
		"low":      1,
		"normal":   2,
		"high":     3,
		"critical": 4,
	}

	if descending {
		sort.Slice(notifs, func(i, j int) bool {
			return priorityMap[notifs[i].Priority] > priorityMap[notifs[j].Priority]
		})
	} else {
		sort.Slice(notifs, func(i, j int) bool {
			return priorityMap[notifs[i].Priority] < priorityMap[notifs[j].Priority]
		})
	}
}

func sortByReadStatus(notifs []types.Notification, descending bool) {
	if descending {
		sort.Slice(notifs, func(i, j int) bool {
			return notifs[i].Read && !notifs[j].Read
		})
	} else {
		sort.Slice(notifs, func(i, j int) bool {
			return !notifs[i].Read && notifs[j].Read
		})
	}
}

func sortByGroupID(notifs []types.Notification, descending bool) {
	if descending {
		sort.Slice(notifs, func(i, j int) bool {
			return notifs[i].GroupID > notifs[j].GroupID
		})
	} else {
		sort.Slice(notifs, func(i, j int) bool {
			return notifs[i].GroupID < notifs[j].GroupID
		})
	}
}

func Broadcast(n types.Notification) {
	for _, r := range n.Recipients {
		subscribers.Lock()
		for _, ch := range subscribers.clients[r.ID] {
			select {
			case ch <- n:
				log.Printf("[Broadcast] Notification sent to userId %s: %+v", r.ID, n)
			default:
				log.Printf("[Broadcast] Channel full for userId %s, notification dropped: %+v", r.ID, n)
			}
		}
		subscribers.Unlock()
	}
}
