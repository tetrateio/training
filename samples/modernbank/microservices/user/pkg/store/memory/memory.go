package memory

import (
	"context"
	"sync"

	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/store"
)

// Enforce that InMemory matches the Store Interface
var _ store.Interface = &InMemory{}

func NewInMemory() *InMemory {
	return &InMemory{m: &sync.RWMutex{}, store: map[string]model.User{}}
}

type InMemory struct {
	m     *sync.RWMutex
	store map[string]model.User
}

func (m *InMemory) Get(ctx context.Context, username string) (*model.User, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	tmp, ok := m.store[username]
	if !ok {
		return nil, &store.NotFound{}
	}
	return &tmp, nil
}

func (m *InMemory) Create(ctx context.Context, user *model.User) (*model.User, error) {
	m.m.Lock()
	defer m.m.Unlock()
	if _, ok := m.store[*user.Username]; ok {
		return nil, &store.Conflict{}
	}
	m.store[*user.Username] = *user
	tmp := m.store[*user.Username]
	return &tmp, nil
}

func (m *InMemory) Update(ctx context.Context, username string, user *model.User) (*model.User, error) {
	m.m.Lock()
	defer m.m.Unlock()
	if _, ok := m.store[username]; !ok {
		return nil, &store.NotFound{}
	}
	delete(m.store, username)
	m.store[*user.Username] = *user
	tmp := m.store[*user.Username]
	return &tmp, nil
}

func (m *InMemory) Delete(ctx context.Context, username string) error {
	m.m.Lock()
	defer m.m.Unlock()
	if _, ok := m.store[username]; !ok {
		return &store.NotFound{}
	}
	delete(m.store, username)
	return nil
}
