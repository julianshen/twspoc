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
	ID            string         `json:"id" bson:"_id,omitempty" rethinkdb:"id,omitempty"`
	Timestamp     time.Time      `json:"timestamp" bson:"timestamp" rethinkdb:"timestamp"`
	Title         string         `json:"title" bson:"title" rethinkdb:"title"`
	Message       string         `json:"message" bson:"message" rethinkdb:"message"`
	Priority      string         `json:"priority" bson:"priority" rethinkdb:"priority"`
	Read          bool           `json:"read" bson:"read" rethinkdb:"read"`
	Recipients    []Recipient    `json:"recipients" bson:"recipients" rethinkdb:"recipients"`
	Labels        []string       `json:"labels,omitempty" bson:"labels,omitempty" rethinkdb:"labels,omitempty"`
	Attachments   []Attachment   `json:"attachments,omitempty" bson:"attachments,omitempty" rethinkdb:"attachments,omitempty"`
	AppName       string         `json:"appName,omitempty" bson:"appName,omitempty" rethinkdb:"appName,omitempty"`
	AppIcon       string         `json:"appIcon,omitempty" bson:"appIcon,omitempty" rethinkdb:"appIcon,omitempty"`
	Expiry        *time.Time     `json:"expiry,omitempty" bson:"expiry,omitempty" rethinkdb:"expiry,omitempty"`
	ActionButtons []ActionButton `json:"actionButtons,omitempty" bson:"actionButtons,omitempty" rethinkdb:"actionButtons,omitempty"`
	GroupID       string         `json:"groupId,omitempty" bson:"groupId,omitempty" rethinkdb:"groupId,omitempty"`
	DeletedFor    []string       `json:"deletedFor,omitempty" bson:"deletedFor,omitempty" rethinkdb:"deletedFor,omitempty"`
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
