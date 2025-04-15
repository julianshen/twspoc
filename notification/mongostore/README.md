# MongoDB Notification Repository

This package provides a MongoDB implementation of the `NotificationRepository` interface for the notification service.

## Features

- Full implementation of the `NotificationRepository` interface using MongoDB
- Support for all CRUD operations (Create, Read, Update, Delete)
- Efficient querying with indexes on commonly used fields
- Real-time notification delivery using MongoDB Change Streams
- Proper handling of MongoDB ObjectIDs and string IDs

## Usage

To use the MongoDB notification repository, set the `DB_TYPE` environment variable to `mongo` when running the notification service:

```bash
DB_TYPE=mongo MONGODB_URI=mongodb://localhost:27017 DB_NAME=notifdb go run main.go
```

Or when using Docker Compose:

```yaml
environment:
  - DB_TYPE=mongo
  - MONGODB_URI=mongodb://mongodb:27017
  - DB_NAME=notifdb
```

## Implementation Details

### Data Storage

Notifications are stored in a MongoDB collection with the following structure:

- Each notification is stored as a document
- The notification ID can be either a MongoDB ObjectID or a custom string ID
- Indexes are created on `recipients.id` and `timestamp` fields for efficient querying

### Real-time Notifications

The `Subscribe` method uses MongoDB Change Streams to provide real-time notification delivery:

1. A change stream is created with a pipeline that filters for insert operations
2. The pipeline filters for documents where the recipient ID matches the subscribed user
3. When a matching document is inserted, it's sent to the notification channel
4. The change stream is properly closed when the context is canceled

### Error Handling

The implementation includes proper error handling for:

- Connection failures
- Query errors
- Document not found scenarios
- Change stream errors

## Dependencies

- `go.mongodb.org/mongo-driver/mongo` - Official MongoDB Go driver
- `go.mongodb.org/mongo-driver/bson` - BSON encoding/decoding
- `go.mongodb.org/mongo-driver/bson/primitive` - MongoDB primitive types
