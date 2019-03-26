package store

import (
	"context"
	"fmt"

	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/model"
)

type Interface interface {
	ListSent(ctx context.Context, account int64) ([]*model.Transaction, error)
	ListReceived(ctx context.Context, account int64) ([]*model.Transaction, error)
	Create(ctx context.Context, transaction *model.Newtransaction) (*model.Transaction, error)
}

type Conflict struct{}

func (c *Conflict) Error() string {
	return fmt.Sprintf("resource already exists")
}

type NotFound struct{}

func (n *NotFound) Error() string {
	return fmt.Sprintf("resource does not exist")
}
