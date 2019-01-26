package store

import (
	"fmt"
	"sync"

	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/model"
)

func NewInMemory() *InMemory {
	return &InMemory{m: &sync.RWMutex{}, store: map[string]*model.User{}}
}

type InMemory struct {
	m     *sync.RWMutex
	store map[string]*model.User
}

func (m *InMemory) Get(username string) (*model.User, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	if _, ok := m.store[username]; !ok {
		return nil, &NotFound{}
	}
	return m.store[username], nil
}

func (m *InMemory) Create(user *model.User) (*model.User, error) {
	m.m.Lock()
	defer m.m.Unlock()
	if _, ok := m.store[user.Username]; ok {
		return nil, fmt.Errorf("user %q already exists", user.Username)
	}
	m.store[user.Username] = user
	return m.store[user.Username], nil
}

func (m *InMemory) Update(username string, user *model.User) (*model.User, error) {
	m.m.Lock()
	defer m.m.Unlock()
	if _, ok := m.store[username]; !ok {
		return nil, &NotFound{}
	}
	delete(m.store, username)
	m.store[user.Username] = user
	return m.store[user.Username], nil
}

func (m *InMemory) Delete(username string) error {
	m.m.Lock()
	defer m.m.Unlock()
	if _, ok := m.store[username]; !ok {
		return &NotFound{}
	}
	delete(m.store, username)
	return nil
}
