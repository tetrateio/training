package store

import (
	"fmt"

	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/model"
)

type Interface interface {
	ListSent(account int64) ([]*model.Transaction, error)
	ListReceived(account int64) ([]*model.Transaction, error)
	GetSent(account int64, id int64) (*model.Transaction, error)
	GetReceived(account int64, id int64) (*model.Transaction, error)
	Create(transaction *model.Newtransaction) (*model.Transaction, error)
}

type Conflict struct{}

func (c *Conflict) Error() string {
	return fmt.Sprintf("resource already exists")
}

type NotFound struct{}

func (n *NotFound) Error() string {
	return fmt.Sprintf("resource does not exist")
}
