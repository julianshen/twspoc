# Notification Service SDK

This SDK provides a convenient way to interact with the Notification Service API. It includes methods for sending notifications, retrieving notifications, marking notifications as read, deleting notifications, searching notifications, and subscribing to real-time notification updates.

## Installation

```bash
go get github.com/yourusername/notification
```

## Usage

### Creating a Client

```go
import (
    "notification/sdk"
    "time"
)

// Create a client with default options
client := sdk.NewClient("http://localhost:3000")

// Create a client with custom timeout
client := sdk.NewClient("http://localhost:3000", sdk.WithTimeout(5*time.Second))

// Create a client with custom HTTP client
httpClient := &http.Client{
    Timeout: 10 * time.Second,
    // Other customizations...
}
client := sdk.NewClient("http://localhost:3000", sdk.WithHTTPClient(httpClient))
```

### Sending a Notification

```go
notification := types.Notification{
    Title:     "Test Notification",
    Message:   "This is a test notification",
    Priority:  "normal",
    Read:      false,
    Recipients: []types.Recipient{
        {Type: "user", ID: "user123"},
    },
    // Other fields...
}

err := client.SendNotification(ctx, notification)
if err != nil {
    // Handle error
}
```

### Getting Notifications for a User

```go
notifications, err := client.GetNotifications(ctx, "user123")
if err != nil {
    // Handle error
}

for _, notification := range notifications {
    fmt.Printf("Notification: %s - %s\n", notification.Title, notification.Message)
}
```

### Marking a Notification as Read

```go
err := client.MarkAsRead(ctx, notificationID)
if err != nil {
    // Handle error
}
```

### Deleting a Notification

```go
err := client.DeleteNotification(ctx, notificationID)
if err != nil {
    // Handle error
}
```

### Searching Notifications

```go
searchParams := sdk.SearchParams{
    Keyword:   "test",
    Title:     "",
    Message:   "",
    Labels:    []string{"important"},
    AppName:   "MyApp",
    StartDate: "2023-01-01T00:00:00Z",
    EndDate:   "2023-12-31T23:59:59Z",
    UserID:    "user123",
}

results, err := client.SearchNotifications(ctx, searchParams)
if err != nil {
    // Handle error
}

for _, notification := range results {
    fmt.Printf("Found: %s - %s\n", notification.Title, notification.Message)
}
```

### Subscribing to Notification Updates

```go
eventCh, err := client.SubscribeToNotifications(ctx, "user123")
if err != nil {
    // Handle error
}

for event := range eventCh {
    if event.Error != nil {
        fmt.Printf("Error: %v\n", event.Error)
        continue
    }
    
    notification := event.Notification
    fmt.Printf("Received: %s - %s\n", notification.Title, notification.Message)
}
```

## Complete Example

See [sdk/example/main.go](example/main.go) for a complete example of how to use the SDK.

## Error Handling

All methods return errors that should be checked and handled appropriately. The errors include detailed information about what went wrong, including HTTP status codes and response bodies when applicable.

## Concurrency

The SDK is safe for concurrent use by multiple goroutines.

## Cancellation

All methods accept a context.Context parameter that can be used to cancel operations. This is particularly useful for long-running operations like subscribing to notification updates.

## Customization

The SDK can be customized using the provided options:

- `WithTimeout`: Sets the timeout for HTTP requests
- `WithHTTPClient`: Sets a custom HTTP client for advanced customization

## License

This SDK is licensed under the MIT License. See the LICENSE file for details.
