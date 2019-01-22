package service

import (
	"fmt"

	"github.com/tetrateio/training/samples/modernbank/user/pkg/model"
)

type Store interface {
	Get(username string) (*model.User, error)
	Create(user *model.User) error
	Update(username string, user *model.User) error
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

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{store: map[string]*model.User{}}
}

type InMemoryStore struct {
	store map[string]*model.User
}

func (m *InMemoryStore) Get(username string) (*model.User, error) {
	if _, ok := m.store[username]; !ok {
		return nil, &NotFound{}
	}
	return m.store[username], nil
}

func (m *InMemoryStore) Create(user *model.User) error {
	if _, ok := m.store[user.Username]; ok {
		return fmt.Errorf("user %q already exists", user.Username)
	}
	m.store[user.Username] = user
	return nil
}

func (m *InMemoryStore) Update(username string, user *model.User) error {
	if _, ok := m.store[username]; !ok {
		return &NotFound{}
	}
	delete(m.store, username)
	m.store[user.Username] = user
	return nil
}

func (m *InMemoryStore) Delete(username string) error {
	if _, ok := m.store[username]; !ok {
		return &NotFound{}
	}
	delete(m.store, username)
	return nil
}
