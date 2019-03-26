package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/store"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Enforce that MongoDB matches the Store Interface
	_ store.Interface = MongoDB{}

	defaultAddress    = "mongodb://user-mongodb:27017"
	defaultDatabase   = "users"
	defaultCollection = "users"
)

func NewMongoDB() MongoDB {
	client, _ := mongo.NewClient(options.Client().ApplyURI(defaultAddress))
	// Keep retrying every 5 seconds until the mongo backend is up or 6 minutes have passed.
	for i := 1; i < 360; i += 5 {
		time.Sleep(5 * time.Second)
		log.Printf("attempting to connect to mongodb at %v", defaultAddress)
		if err := client.Connect(context.Background()); err != nil {
			log.Printf("unable to connect to mongodb: %v", err)
		}
		if err := client.Ping(context.Background(), nil); err != nil {
			log.Printf("unable to ping mongodb: %v", err)
		} else {
			break
		}
	}

	return MongoDB{client: client}
}

type MongoDB struct {
	client *mongo.Client
}

func (m MongoDB) Get(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := m.defaultCollection().FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return nil, &store.NotFound{}
		}
		return nil, fmt.Errorf("unable to get user in database: %v", err)
	}
	return &user, nil
}

func (m MongoDB) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if _, err := m.defaultCollection().InsertOne(ctx, *user); err != nil {
		return nil, fmt.Errorf("unable to create user in database: %v", err)
	}
	return user, nil
}

func (m MongoDB) Update(ctx context.Context, username string, user *model.User) (*model.User, error) {
	res, err := m.defaultCollection().UpdateOne(ctx, bson.M{"username": username}, *user)
	if res != nil && res.UpsertedCount == 0 {
		return nil, &store.NotFound{}
	}
	if err != nil {
		return nil, fmt.Errorf("unable to update user in database: %v", err)
	}
	return user, nil
}

func (m MongoDB) Delete(ctx context.Context, username string) error {
	res, err := m.defaultCollection().DeleteOne(ctx, bson.M{"username": username})
	if res != nil && res.DeletedCount == 0 {
		return &store.NotFound{}
	}
	if err != nil {
		return fmt.Errorf("unable to delete user in database: %v", err)
	}
	return nil
}

func (m MongoDB) defaultCollection() *mongo.Collection {
	return m.client.Database(defaultDatabase).Collection(defaultCollection)
}
