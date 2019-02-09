package user

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/time/rate"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/go-openapi/swag"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/client"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/client/users"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/tools/trafficGen/pkg/store"
)

type Creator struct {
	client  *client.User
	store   store.Interface
	limiter *rate.Limiter
	count   int64
}

func NewCreator(host string, userStore store.Interface, limit rate.Limit) *Creator {
	transportConfig := client.DefaultTransportConfig().WithHost(host)
	userClient := client.NewHTTPClientWithConfig(nil, transportConfig)
	return &Creator{
		client:  userClient,
		store:   userStore,
		limiter: rate.NewLimiter(limit, 1),
		count:   0,
	}
}

func (c *Creator) Run(ctx context.Context) {
	rand.Seed(time.Now().UnixNano())
	c.prepopulate(50)
	for {
		select {
		case <-ctx.Done():
			log.Printf("User creation terminating, successfully created %v users.", c.count)
			return
		default:
			c.limiter.Wait(ctx)
			c.createUser()
		}
	}
}
func (c *Creator) prepopulate(population int64) {
	for c.store.UserCount() < population {
		c.createUser()
	}
}

func (c *Creator) createUser() {
	u := genUser()
	params := users.NewCreateUserParams().WithBody(u)
	created, err := c.client.Users.CreateUser(params)
	if err != nil {
		log.Printf("Error creating user %q: %v", *u.Username, err)
		return
	}
	log.Printf("Successfully created user %q, adding to internal store.", *created.Payload.Username)
	c.count++
	c.store.AddUser(*created.Payload.Username)
}

func genUser() *model.User {
	adjective, adverb, animal := petname.Adjective(), petname.Adverb(), petname.Name()
	return &model.User{
		Email:     swag.String(fmt.Sprintf("%v@%v.%v", adverb, adjective, animal)),
		FirstName: &adjective,
		LastName:  &animal,
		Password:  swag.String(genRandString(25)),
		Username:  swag.String(fmt.Sprintf("%v%v%v", adjective, animal, strconv.Itoa((rand.Intn(99999))))),
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!_@#%^&*(){}]["

func genRandString(length int) string {
	p := make([]byte, length)
	for i := range p {
		p[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(p)
}
