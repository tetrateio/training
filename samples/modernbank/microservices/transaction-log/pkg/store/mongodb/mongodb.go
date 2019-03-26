package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-openapi/swag"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/store"
)

var (
	// Enforce that MongoDB matches the Store Interface
	_ store.Interface = MongoDB{}

	defaultAddress    = "mongodb://transaction-log-mongodb:27017"
	defaultDatabase   = "transactions"
	defaultCollection = "transactions"
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

func (m MongoDB) ListSent(ctx context.Context, account int64) ([]*model.Transaction, error) {
	transactions := []*model.Transaction{}
	res, err := m.defaultCollection().Find(ctx, bson.M{"sender": account})
	if err != nil {
		return nil, fmt.Errorf("unable to get transactions in database: %v", err)
	}
	defer res.Close(ctx)
	for res.Next(ctx) {
		var transaction model.Transaction
		if err := res.Decode(&transaction); err != nil {
			return nil, fmt.Errorf("unable to decode response: %v", err)
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func (m MongoDB) ListReceived(ctx context.Context, account int64) ([]*model.Transaction, error) {
	transactions := []*model.Transaction{}
	res, err := m.defaultCollection().Find(ctx, bson.M{"receiver": account})
	if err != nil {
		return nil, fmt.Errorf("unable to get transactions in database: %v", err)
	}
	defer res.Close(ctx)
	for res.Next(ctx) {
		var transaction model.Transaction
		if err := res.Decode(&transaction); err != nil {
			return nil, fmt.Errorf("unable to decode response: %v", err)
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func (m MongoDB) Create(ctx context.Context, transaction *model.Newtransaction) (*model.Transaction, error) {
	new := &model.Transaction{
		Amount:    transaction.Amount,
		Sender:    transaction.Sender,
		Receiver:  transaction.Receiver,
		Timestamp: swag.Int64(time.Now().Unix()),
	}
	res, err := m.defaultCollection().InsertOne(ctx, new)
	if err != nil {
		return nil, fmt.Errorf("unable to create transaction in database: %v", err)
	}
	new.ID = swag.String(res.InsertedID.(primitive.ObjectID).Hex())
	return new, nil
}

func (m MongoDB) defaultCollection() *mongo.Collection {
	return m.client.Database(defaultDatabase).Collection(defaultCollection)
}
