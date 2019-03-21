package store

import (
	"fmt"

	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/model"
)

type Interface interface {
	List(owner string) ([]*model.Account, error)
	Get(owner string, number int64) (*model.Account, error)
	Create(owner string, accountType string) (*model.Account, error)
	Delete(owner string, number int64) error
	UpdateBalance(number int64, deltaAmount float64) error
}

type Conflict struct{}

func (c *Conflict) Error() string {
	return fmt.Sprintf("resource already exists")
}

type NotFound struct{}

func (n *NotFound) Error() string {
	return fmt.Sprintf("resource does not exist")
}
