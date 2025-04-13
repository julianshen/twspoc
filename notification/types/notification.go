// types/notification.go
package types

import (
	"time"
)

// Recipient represents a notification recipient
type Recipient struct {
	Type string `bson:"type" json:"type"`
	ID   string `bson:"id" json:"id"`
}

// Attachment represents a notification attachment
type Attachment struct {
	Type string `bson:"type" json:"type"`
	ID   string `bson:"id" json:"id"`
	URL  string `bson:"url" json:"url"`
}

// ActionButton represents an action button in a notification
type ActionButton struct {
	Label  string `bson:"label" json:"label"`
	Action string `bson:"action" json:"action"`
	URL    string `bson:"url" json:"url"`
}

// Notification represents a notification message
type Notification struct {
	ID            string         `json:"id" rethinkdb:"id,omitempty"`
	Timestamp     time.Time      `json:"timestamp" rethinkdb:"timestamp"`
	Title         string         `json:"title" rethinkdb:"title"`
	Message       string         `json:"message" rethinkdb:"message"`
	Priority      string         `json:"priority" rethinkdb:"priority"`
	Read          bool           `json:"read" rethinkdb:"read"`
	Recipients    []Recipient    `json:"recipients" rethinkdb:"recipients"`
	Labels        []string       `json:"labels,omitempty" rethinkdb:"labels,omitempty"`
	Attachments   []Attachment   `json:"attachments,omitempty" rethinkdb:"attachments,omitempty"`
	AppName       string         `json:"appName,omitempty" rethinkdb:"appName,omitempty"`
	AppIcon       string         `json:"appIcon,omitempty" rethinkdb:"appIcon,omitempty"`
	Expiry        *time.Time     `json:"expiry,omitempty" rethinkdb:"expiry,omitempty"`
	ActionButtons []ActionButton `json:"actionButtons,omitempty" rethinkdb:"actionButtons,omitempty"`
	GroupID       string         `json:"groupId,omitempty" rethinkdb:"groupId,omitempty"`
	DeletedFor    []string       `json:"deletedFor,omitempty" rethinkdb:"deletedFor,omitempty"`
}

// SearchParams represents search parameters for notifications
type SearchParams struct {
	Keyword   string   `json:"keyword"`
	Title     string   `json:"title"`
	Message   string   `json:"message"`
	Labels    []string `json:"labels"`
	AppName   string   `json:"appName"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	UserID    string   `json:"userId"`
}
