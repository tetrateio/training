package store

import (
	"context"
	"fmt"

	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/model"
)

type Interface interface {
	List(ctx context.Context, owner string) ([]*model.Account, error)
	Get(ctx context.Context, owner string, number int64) (*model.Account, error)
	Create(ctx context.Context, owner string, accountType string) (*model.Account, error)
	Delete(ctx context.Context, owner string, number int64) error
	UpdateBalance(ctx context.Context, number int64, deltaAmount float64) error
}

type Conflict struct{}

func (c *Conflict) Error() string {
	return fmt.Sprintf("resource already exists")
}

type NotFound struct{}

func (n *NotFound) Error() string {
	return fmt.Sprintf("resource does not exist")
}
