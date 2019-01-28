package store

import (
	"fmt"

	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/model"
)

type Interface interface {
	Get(username string) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	Update(username string, user *model.User) (*model.User, error)
	Delete(username string) error
}

type Conflict struct{}

func (c *Conflict) Error() string {
	return fmt.Sprintf("resource already exists")
}

type NotFound struct{}

func (n *NotFound) Error() string {
	return fmt.Sprintf("resource does not exist")
}
