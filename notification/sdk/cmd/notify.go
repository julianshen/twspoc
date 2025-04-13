// sdk/cmd/notify.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"notification/sdk"
	"notification/types"
)

func main() {
	var (
		baseURL    string
		title      string
		message    string
		priority   string
		recipients string
		labels     string
		appName    string
		appIcon    string
		groupId    string
	)

	flag.StringVar(&baseURL, "url", "http://localhost:3000", "Notification service base URL")
	flag.StringVar(&title, "title", "", "Notification title (required)")
	flag.StringVar(&message, "message", "", "Notification message (required)")
	flag.StringVar(&priority, "priority", "normal", "Priority (low|normal|high|critical)")
	flag.StringVar(&recipients, "recipients", "", "Comma-separated recipient user IDs (required)")
	flag.StringVar(&labels, "labels", "", "Comma-separated labels")
	flag.StringVar(&appName, "app", "", "App name")
	flag.StringVar(&appIcon, "icon", "", "App icon URL")
	flag.StringVar(&groupId, "group", "", "Group ID")
	flag.Parse()

	if title == "" || message == "" || recipients == "" {
		fmt.Fprintln(os.Stderr, "title, message, and recipients are required")
		flag.Usage()
		os.Exit(1)
	}

	recipList := strings.Split(recipients, ",")
	var recips []types.Recipient
	for _, r := range recipList {
		r = strings.TrimSpace(r)
		if r != "" {
			recips = append(recips, types.Recipient{Type: "user", ID: r})
		}
	}
	if len(recips) == 0 {
		fmt.Fprintln(os.Stderr, "At least one recipient is required")
		os.Exit(1)
	}

	var labelList []string
	if labels != "" {
		for _, l := range strings.Split(labels, ",") {
			labelList = append(labelList, strings.TrimSpace(l))
		}
	}

	now := time.Now()
	notification := types.Notification{
		ID:         fmt.Sprintf("notif-%d", now.UnixNano()),
		Timestamp:  now,
		Title:      title,
		Message:    message,
		Priority:   priority,
		Read:       false,
		Recipients: recips,
		Labels:     labelList,
		AppName:    appName,
		AppIcon:    appIcon,
		GroupID:    groupId,
	}

	client := sdk.NewClient(baseURL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.SendNotification(ctx, notification); err != nil {
		log.Fatalf("Failed to send notification: %v", err)
	}
	fmt.Println("Notification sent successfully")
}
