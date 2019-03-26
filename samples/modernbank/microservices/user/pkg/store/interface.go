package store

import (
	"context"
	"fmt"

	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/model"
)

type Interface interface {
	Get(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, username string, user *model.User) (*model.User, error)
	Delete(ctx context.Context, username string) error
}

type Conflict struct{}

func (c *Conflict) Error() string {
	return fmt.Sprintf("resource already exists")
}

type NotFound struct{}

func (n *NotFound) Error() string {
	return fmt.Sprintf("resource does not exist")
}
