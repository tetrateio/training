package store

import (
	"math"

	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/model"
)

func NewInMemory() *InMemory {
	return &InMemory{store: map[int64]*model.Account{}}
}

type InMemory struct {
	store map[int64]*model.Account
}

func (m *InMemory) Get(number int64) (*model.Account, error) {
	if _, ok := m.store[number]; !ok {
		return nil, &NotFound{}
	}
	return m.store[number], nil
}

func (m *InMemory) Create(owner string) (*model.Account, error) {
	accountNumber := m.unAssignedAccountNumber()
	m.store[accountNumber] = &model.Account{
		Balance: 0,
		Owner:   owner,
		Number:  accountNumber,
	}
	return m.store[accountNumber], nil
}

func (m *InMemory) Update(number int64, account *model.Account) (*model.Account, error) {
	if _, ok := m.store[number]; !ok {
		return nil, &NotFound{}
	}
	delete(m.store, number)
	m.store[account.Number] = account
	return m.store[account.Number], nil
}

func (m *InMemory) Delete(number int64) error {
	if _, ok := m.store[number]; !ok {
		return &NotFound{}
	}
	delete(m.store, number)
	return nil
}

func (m *InMemory) unAssignedAccountNumber() int64 {
	for i := int64(0); i < int64(math.MaxInt64); i++ {
		if _, ok := m.store[i]; !ok {
			return i
		}
	}
	// Obviously don't actually panic in a real life scenario...
	panic("we have run out of account numbers")
}
