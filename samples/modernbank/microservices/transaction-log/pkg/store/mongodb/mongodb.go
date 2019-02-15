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
	_ store.Interface = NewMongoDB()

	defaultAddress    = "mongodb://transaction-log-mongodb:27017"
	defaultDatabase   = "transactions"
	defaultCollection = "transactions"

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

func (m *MongoDB) ListSent(account int64) ([]*model.Transaction, error) {
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
	if len(transactions) == 0 {
		return transactions, &store.NotFound{}
	}
	return transactions, nil
}

func (m *MongoDB) ListReceived(account int64) ([]*model.Transaction, error) {
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
	if len(transactions) == 0 {
		return transactions, &store.NotFound{}
	}
	return transactions, nil
}

func (m *MongoDB) GetSent(account int64, id string) (*model.Transaction, error) {
	var transaction model.Transaction
	res := m.defaultCollection().FindOne(ctx, bson.M{"_id": id, "sender": account})
	if res.Err().Error() == mongo.ErrNoDocuments.Error() {
		return nil, &store.NotFound{}
	} else if res.Err() != nil {
		return nil, fmt.Errorf("unable to get transaction in database: %v", res.Err())
	}
	return &transaction, res.Decode(&transaction)
}

func (m *MongoDB) GetReceived(account int64, id string) (*model.Transaction, error) {
	var transaction model.Transaction
	res := m.defaultCollection().FindOne(ctx, bson.M{"_id": id, "receiver": account})
	if res.Err().Error() == mongo.ErrNoDocuments.Error() {
		return nil, &store.NotFound{}
	} else if res.Err() != nil {
		return nil, fmt.Errorf("unable to get transaction in database: %v", res.Err())
	}
	return &transaction, res.Decode(&transaction)
}

func (m *MongoDB) Create(transaction *model.Newtransaction) (*model.Transaction, error) {
	res, err := m.defaultCollection().InsertOne(ctx, transaction)
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

func (m *MongoDB) defaultCollection() *mongo.Collection {
	return m.client.Database(defaultDatabase).Collection(defaultCollection)
}
