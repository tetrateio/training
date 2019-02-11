package transaction

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/client"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/client/transactions"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/tools/trafficGen/pkg/store"
	"golang.org/x/time/rate"
)

type Creator struct {
	client  *client.Transaction
	store   store.Interface
	limiter *rate.Limiter
	count   int64
}

func NewCreator(host string, userStore store.Interface, limit rate.Limit) *Creator {
	transportConfig := client.DefaultTransportConfig().WithHost(host)
	transactionClient := client.NewHTTPClientWithConfig(nil, transportConfig)
	return &Creator{
		client:  transactionClient,
		store:   userStore,
		limiter: rate.NewLimiter(limit, 1),
		count:   0,
	}
}

func (c *Creator) Run(ctx context.Context) {
	time.Sleep(time.Second * 10) // allow some accounts to be created...
	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case <-ctx.Done():
			log.Printf("Transaction creation terminating, successfully created %v transactions.", c.count)
			return
		default:
			c.limiter.Wait(ctx)
			c.createTransaction()
		}
	}
}

func (c *Creator) createTransaction() {
	t := c.genTransaction()
	params := transactions.NewCreateTransactionParams().WithBody(t)
	created, err := c.client.Transactions.CreateTransaction(params)
	if err != nil {
		log.Printf("Error creating transaction from %v to %v for %v: %v", *t.Sender, *t.Receiver, *t.Amount, err)
		return
	}
	res := created.Payload
	log.Printf("Successfully created transaction %v from %v to %v for %v.", *res.ID, *res.Sender, *res.Receiver, *res.Amount)
	c.count++
}

func (c *Creator) genTransaction() *model.Newtransaction {
	sender, receiver := c.store.GetRandomAccount(), c.store.GetRandomAccount()
	randFloat := rand.Float64() * 1000
	amount := float64(int(randFloat*100)) / 100 // 2 d.p. precision!
	return &model.Newtransaction{
		Amount:   &amount,
		Sender:   &sender,
		Receiver: &receiver,
	}
}
