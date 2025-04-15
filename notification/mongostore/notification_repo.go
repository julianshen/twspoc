package mongostore

import (
	"context"
	"notification/notificationrepo"
	"notification/types"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Compile-time check to ensure MongoNotificationRepository implements NotificationRepository
var _ notificationrepo.NotificationRepository = (*MongoNotificationRepository)(nil)

// MongoNotificationRepository implements the NotificationRepository interface using MongoDB.
type MongoNotificationRepository struct {
	client        *mongo.Client
	collection    *mongo.Collection
	dbName        string
	mu            sync.Mutex
	changeStreams map[string]*mongo.ChangeStream
}

// NewMongoNotificationRepository creates a new instance of MongoNotificationRepository.
// It requires a MongoDB connection URI, database name, and collection name.
func NewMongoNotificationRepository(ctx context.Context, uri, dbName, collectionName string) (*MongoNotificationRepository, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the primary
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	// Create indexes for better query performance
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "recipients.id", Value: 1}},
			Options: options.Index().SetBackground(true),
		},
		{
			Keys:    bson.D{{Key: "timestamp", Value: 1}},
			Options: options.Index().SetBackground(true),
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return nil, err
	}

	return &MongoNotificationRepository{
		client:        client,
		collection:    collection,
		dbName:        dbName,
		changeStreams: make(map[string]*mongo.ChangeStream),
	}, nil
}

// Close disconnects the MongoDB client.
func (r *MongoNotificationRepository) Close(ctx context.Context) error {
	if r.client != nil {
		return r.client.Disconnect(ctx)
	}
	return nil
}

// Create inserts a new notification into the MongoDB collection.
func (r *MongoNotificationRepository) Create(ctx context.Context, n *types.Notification) error {
	// Set timestamp if not already set
	if n.Timestamp.IsZero() {
		n.Timestamp = time.Now()
	}

	// If ID is empty, MongoDB will generate one
	if n.ID == "" {
		// Generate a new ObjectID and convert it to string
		objID := primitive.NewObjectID()
		n.ID = objID.Hex()
	}

	_, err := r.collection.InsertOne(ctx, n)
	return err
}

// Get retrieves a notification by its ID from MongoDB.
func (r *MongoNotificationRepository) Get(ctx context.Context, id string) (*types.Notification, error) {
	var notification types.Notification

	// Try to convert the ID to an ObjectID if it's in that format
	var filter bson.M
	objID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		// If it's a valid ObjectID, search by _id
		filter = bson.M{"_id": objID}
	} else {
		// Otherwise, search by the ID field
		filter = bson.M{"id": id}
	}

	// First try to find by _id or id
	err = r.collection.FindOne(ctx, filter).Decode(&notification)
	if err != nil {
		// If not found, try to find by the ID field directly
		if err == mongo.ErrNoDocuments {
			filter = bson.M{"id": id}
			err = r.collection.FindOne(ctx, filter).Decode(&notification)
			if err != nil {
				return nil, mongo.ErrNoDocuments
			}
		} else {
			return nil, err
		}
	}

	return &notification, nil
}

// Update modifies an existing notification in MongoDB.
func (r *MongoNotificationRepository) Update(ctx context.Context, n *types.Notification) error {
	// Try to convert the ID to an ObjectID if it's in that format
	var filter bson.M
	objID, err := primitive.ObjectIDFromHex(n.ID)
	if err == nil {
		// If it's a valid ObjectID, search by _id
		filter = bson.M{"_id": objID}
	} else {
		// Otherwise, search by the ID field
		filter = bson.M{"id": n.ID}
	}

	update := bson.M{"$set": n}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete removes a notification by its ID from MongoDB.
func (r *MongoNotificationRepository) Delete(ctx context.Context, id string) error {
	// Try to convert the ID to an ObjectID if it's in that format
	var filter bson.M
	objID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		// If it's a valid ObjectID, search by _id
		filter = bson.M{"_id": objID}
	} else {
		// Otherwise, search by the ID field
		filter = bson.M{"id": id}
	}

	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}

// ListByUser retrieves all notifications for a specific user ID from MongoDB.
func (r *MongoNotificationRepository) ListByUser(ctx context.Context, userID string) ([]*types.Notification, error) {
	// Find notifications where the user is a recipient and not in deletedFor
	filter := bson.M{
		"recipients.id": userID,
		"$or": []bson.M{
			{"deletedFor": bson.M{"$exists": false}},
			{"deletedFor": bson.M{"$nin": []string{userID}}},
		},
	}

	// Sort by timestamp descending (newest first)
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []*types.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

// Subscribe creates a MongoDB change stream to listen for new notifications for a user.
func (r *MongoNotificationRepository) Subscribe(ctx context.Context, userID string) (<-chan *types.Notification, error) {
	notificationChan := make(chan *types.Notification, 100)

	// Pipeline to filter changes for the specific user
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "operationType", Value: "insert"},
			{Key: "fullDocument.recipients.id", Value: userID},
		}}},
	}

	// Options for the change stream
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	// Create a change stream
	changeStream, err := r.collection.Watch(ctx, pipeline, opts)
	if err != nil {
		close(notificationChan)
		return nil, err
	}

	// Start a goroutine to listen for changes
	go func() {
		defer changeStream.Close(ctx)
		defer close(notificationChan)

		for changeStream.Next(ctx) {
			var changeDoc struct {
				FullDocument types.Notification `bson:"fullDocument"`
			}
			if err := changeStream.Decode(&changeDoc); err != nil {
				continue
			}

			// Send the notification to the channel
			select {
			case notificationChan <- &changeDoc.FullDocument:
			case <-ctx.Done():
				return
			}
		}
	}()

	return notificationChan, nil
}
