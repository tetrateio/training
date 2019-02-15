package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/store"
)

var (
	// Enforce that MongoDB matches the Store Interface
	_ store.Interface = NewMongoDB()

	defaultAddress    = "mongodb://user-mongodb:27017"
	defaultDatabase   = "users"
	defaultCollection = "users"

	ctx = context.Background()
)

func NewMongoDB() *MongoDB {
	client, _ := mongo.NewClient(defaultAddress)
	for i := 1; i < 360; i += 5 {
		time.Sleep(5 * time.Second)
		log.Printf("attempting to connect to mongodb at %v", defaultAddress)
		if err := client.Connect(ctx); err != nil {
			log.Printf("unable to connect to mongodb: %v", err)
		}
		err := client.Ping(ctx, nil)
		if err == nil {
			break
		}
		log.Printf("unable to ping mongodb: %v", err)
	}

	return &MongoDB{client: client}
}

type MongoDB struct {
	client *mongo.Client
}

func (m *MongoDB) Get(username string) (*model.User, error) {
	var user model.User
	res := m.defaultCollection().FindOne(ctx, bson.M{"username": username})
	if res.Err().Error() == mongo.ErrNoDocuments.Error() {
		return nil, &store.NotFound{}
	} else if res.Err() != nil {
		return nil, fmt.Errorf("unable to get user in database: %v", res.Err())
	}
	return &user, res.Decode(&user)
}

func (m *MongoDB) Create(user *model.User) (*model.User, error) {
	_, err := m.defaultCollection().InsertOne(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("unable to create user in database: %v", err)
	}
	return user, nil
}

func (m *MongoDB) Update(username string, user *model.User) (*model.User, error) {
	res, err := m.defaultCollection().UpdateOne(ctx, bson.M{"username": username}, *user)
	if res.UpsertedCount == 0 {
		return nil, &store.NotFound{}
	}
	if err != nil {
		return nil, fmt.Errorf("unable to update user in database: %v", err)
	}
	return user, nil
}

func (m *MongoDB) Delete(username string) error {
	res, err := m.defaultCollection().DeleteOne(ctx, bson.M{"username": username})
	if res.DeletedCount == 0 {
		return &store.NotFound{}
	}
	if err != nil {
		return fmt.Errorf("unable to delete user in database: %v", err)
	}
	return nil
}

func (m *MongoDB) defaultCollection() *mongo.Collection {
	return m.client.Database(defaultDatabase).Collection(defaultCollection)
}
