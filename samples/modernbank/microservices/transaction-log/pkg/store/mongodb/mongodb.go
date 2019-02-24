package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-openapi/swag"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

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
	client, _ := mongo.NewClient(defaultAddress)
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

func (m MongoDB) ListSent(account int64) ([]*model.Transaction, error) {
	transactions := []*model.Transaction{}
	res, err := m.defaultCollection().Find(context.Background(), bson.M{"sender": account})
	if err != nil {
		return nil, fmt.Errorf("unable to get transactions in database: %v", err)
	}
	defer res.Close(context.Background())
	for res.Next(context.Background()) {
		var transaction model.Transaction
		if err := res.Decode(&transaction); err != nil {
			return nil, fmt.Errorf("unable to decode response: %v", err)
		}
		transactions = append(transactions, &transaction)
	}
	if len(transactions) == 0 {
		return transactions, &store.NotFound{}
	}
	return transactions, nil
}

func (m MongoDB) ListReceived(account int64) ([]*model.Transaction, error) {
	transactions := []*model.Transaction{}
	res, err := m.defaultCollection().Find(context.Background(), bson.M{"receiver": account})
	if err != nil {
		return nil, fmt.Errorf("unable to get transactions in database: %v", err)
	}
	defer res.Close(context.Background())
	for res.Next(context.Background()) {
		var transaction model.Transaction
		if err := res.Decode(&transaction); err != nil {
			return nil, fmt.Errorf("unable to decode response: %v", err)
		}
		transactions = append(transactions, &transaction)
	}
	if len(transactions) == 0 {
		return transactions, &store.NotFound{}
	}
	return transactions, nil
}

func (m MongoDB) GetSent(account int64, id string) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := m.defaultCollection().FindOne(context.Background(), bson.M{"_id": id, "sender": account}).Decode(&transaction); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return nil, &store.NotFound{}
		}
		return nil, fmt.Errorf("unable to get transaction in database: %v", err)
	}
	return &transaction, nil
}

func (m MongoDB) GetReceived(account int64, id string) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := m.defaultCollection().FindOne(context.Background(), bson.M{"_id": id, "receiver": account}).Decode(&transaction); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return nil, &store.NotFound{}
		}
		return nil, fmt.Errorf("unable to get transaction in database: %v", err)
	}
	return &transaction, nil
}

func (m MongoDB) Create(transaction *model.Newtransaction) (*model.Transaction, error) {
	res, err := m.defaultCollection().InsertOne(context.Background(), transaction)
	if err != nil {
		return nil, fmt.Errorf("unable to create transaction in database: %v", err)
	}
	return &model.Transaction{
		ID:       swag.String(res.InsertedID.(primitive.ObjectID).Hex()),
		Amount:   transaction.Amount,
		Sender:   transaction.Sender,
		Receiver: transaction.Receiver,
	}, nil
}

func (m MongoDB) defaultCollection() *mongo.Collection {
	return m.client.Database(defaultDatabase).Collection(defaultCollection)
}
