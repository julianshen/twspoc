// sdk/sse.go
package sdk

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"notification/types"
	"strings"
)

// NotificationEvent represents a notification event received from the SSE stream
type NotificationEvent struct {
	Notification types.Notification
	Error        error
}

// SubscribeToNotifications subscribes to notification updates for a user
// It returns a channel that will receive notification events
// The channel will be closed when the context is canceled or an error occurs
func (c *Client) SubscribeToNotifications(ctx context.Context, userID string) (<-chan NotificationEvent, error) {
	url := fmt.Sprintf("%s/api/notifications/subscribe?userId=%s", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	eventCh := make(chan NotificationEvent)

	go func() {
		defer resp.Body.Close()
		defer close(eventCh)

		scanner := bufio.NewScanner(resp.Body)

		var notification types.Notification

		for scanner.Scan() {
			line := scanner.Text()

			// Skip empty lines
			if line == "" {
				continue
			}

			// Check if this is a data line
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")

				// Parse the notification
				if err := json.Unmarshal([]byte(data), &notification); err != nil {
					select {
					case eventCh <- NotificationEvent{Error: fmt.Errorf("failed to parse notification: %w", err)}:
					case <-ctx.Done():
						return
					}
					continue
				}

				// Send the notification to the channel
				select {
				case eventCh <- NotificationEvent{Notification: notification}:
				case <-ctx.Done():
					return
				}
			}
		}

		// Check if the scanner stopped due to an error
		if err := scanner.Err(); err != nil {
			select {
			case eventCh <- NotificationEvent{Error: fmt.Errorf("SSE stream error: %w", err)}:
			case <-ctx.Done():
			}
		}
	}()

	return eventCh, nil
}
