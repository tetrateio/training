package account

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/client"
	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/client/accounts"
	"github.com/tetrateio/training/samples/modernbank/tools/trafficGen/pkg/store"
	"golang.org/x/time/rate"
)

type Creator struct {
	client  *client.Account
	store   store.Interface
	limiter *rate.Limiter
	count   int64
}

func NewCreator(host string, userStore store.Interface, limit rate.Limit) *Creator {
	transportConfig := client.DefaultTransportConfig().WithHost(host)
	accountClient := client.NewHTTPClientWithConfig(nil, transportConfig)
	return &Creator{
		client:  accountClient,
		store:   userStore,
		limiter: rate.NewLimiter(limit, 1),
		count:   0,
	}
}

func (c *Creator) Run(ctx context.Context) {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Second * 5) // allow some users to be created...
	c.prepopulate(100)
	for {
		select {
		case <-ctx.Done():
			log.Printf("Account creation terminating, successfully created %v accounts.", c.count)
			return
		default:
			c.limiter.Wait(ctx)
			c.createAccount()
		}
	}
}

func (c *Creator) prepopulate(population int64) {
	for c.store.AccountCount() < population {
		c.createAccount()
	}
}

func (c *Creator) createAccount() {
	username := c.store.GetRandomUser()
	params := accounts.NewCreateAccountParams().WithOwner(username)
	created, err := c.client.Accounts.CreateAccount(params)
	if err != nil {
		log.Printf("Error creating account for user %q: %v", username, err)
		return
	}
	log.Printf("Successfully created account %v for user %q, adding to internal store.", *created.Payload.Number, *created.Payload.Owner)
	c.count++
	c.store.AddAccount(*created.Payload.Owner, *created.Payload.Number)
}
