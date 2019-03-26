package mongodb

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-openapi/swag"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/store"
)

var (
	// Enforce that MongoDB matches the Store Interface
	_ store.Interface = MongoDB{}

	defaultAddress    = "mongodb://account-mongodb:27017"
	defaultDatabase   = "accounts"
	defaultCollection = "accounts"

	randomAccountNumber = func() int64 {
		// up to 15 digit account numbers
		return rand.Int63n(999999999999999)
	}
)

func NewMongoDB() MongoDB {
	rand.Seed(time.Now().UnixNano())
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

func (m MongoDB) List(ctx context.Context, owner string) ([]*model.Account, error) {
	accounts := []*model.Account{}
	res, err := m.defaultCollection().Find(ctx, bson.M{"owner": owner})
	if err != nil {
		return nil, fmt.Errorf("unable to get accounts in database: %v", err)
	}
	defer res.Close(ctx)
	for res.Next(ctx) {
		var account model.Account
		if err := res.Decode(&account); err != nil {
			return nil, fmt.Errorf("unable to decode response: %v", err)
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

func (m MongoDB) Get(ctx context.Context, owner string, number int64) (*model.Account, error) {
	var account model.Account
	if err := m.defaultCollection().FindOne(ctx, bson.M{"owner": owner, "number": number}).Decode(&account); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return nil, &store.NotFound{}
		}
		return nil, fmt.Errorf("unable to get user in database: %v", err)
	}
	return &account, nil
}

func (m MongoDB) Create(ctx context.Context, owner string, accountType string) (*model.Account, error) {
	newAccountNumber, err := m.unAssignedAccountNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("error finding a vacant account number: %v", err)
	}
	newAccount := &model.Account{
		Balance: swag.Float64(100),
		Owner:   swag.String(owner),
		Number:  swag.Int64(newAccountNumber),
		Type:    swag.String(accountType),
	}
	_, err = m.defaultCollection().InsertOne(ctx, newAccount)
	if err != nil {
		return nil, fmt.Errorf("error creating account in database: %v", err)
	}
	return newAccount, nil
}

// Not concurrency safe but close enough for a demo app
// Clashes are highly unlikely
func (m MongoDB) unAssignedAccountNumber(ctx context.Context) (int64, error) {
	var err error
	candidate, count := int64(0), int64(0)
	for i := 0; i < 10; i++ {
		candidate = randomAccountNumber()
		count, err = m.defaultCollection().CountDocuments(ctx, bson.M{"number": candidate})
		if count == 0 {
			return candidate, nil
		}
	}
	return 0, err
}

func (m MongoDB) Delete(ctx context.Context, owner string, number int64) error {
	res := m.defaultCollection().FindOneAndDelete(ctx, bson.M{"owner": owner, "number": number})
	if res.Err().Error() == mongo.ErrNoDocuments.Error() {
		return &store.NotFound{}
	} else if res.Err() != nil {
		return fmt.Errorf("unable to delete account in database: %v", res.Err())
	}
	return nil
}

func (m MongoDB) UpdateBalance(ctx context.Context, number int64, deltaAmount float64) error {
	res, err := m.defaultCollection().UpdateOne(ctx, bson.M{"number": number}, bson.D{{"$inc", bson.D{{"balance", deltaAmount}}}})
	if res != nil && res.ModifiedCount != 1 {
		return &store.NotFound{}
	} else if err != nil {
		return fmt.Errorf("unable to update account in database: %v", err)
	}
	return nil
}

func (m MongoDB) defaultCollection() *mongo.Collection {
	return m.client.Database(defaultDatabase).Collection(defaultCollection)
}
