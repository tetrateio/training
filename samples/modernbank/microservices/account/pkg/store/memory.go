package store

import (
	"math"

	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/model"
)

var _ Interface = NewInMemory()

func NewInMemory() *InMemory {
	return &InMemory{
		ownerAccounts: map[string]map[int64]*model.Account{},
		accounts:      map[int64]*model.Account{},
	}
}

type InMemory struct {
	ownerAccounts map[string]map[int64]*model.Account
	accounts      map[int64]*model.Account
}

func (m *InMemory) List(owner string) ([]*model.Account, error) {
	if m.ownerAccounts[owner] == nil {
		return nil, &NotFound{}
	}
	res := make([]*model.Account, len(m.ownerAccounts[owner]))
	for _, val := range m.ownerAccounts[owner] {
		res = append(res, val)
	}
	return res, nil
}

func (m *InMemory) Get(owner string, number int64) (*model.Account, error) {
	if _, ok := m.ownerAccounts[owner][number]; !ok {
		return nil, &NotFound{}
	}
	return m.ownerAccounts[owner][number], nil
}

func (m *InMemory) Create(owner string) (*model.Account, error) {
	if m.ownerAccounts[owner] == nil {
		m.ownerAccounts[owner] = map[int64]*model.Account{}
	}
	accountNumber := m.unAssignedAccountNumber()
	m.ownerAccounts[owner][accountNumber] = &model.Account{
		Balance: 0,
		Owner:   owner,
		Number:  accountNumber,
	}
	m.accounts[accountNumber] = m.ownerAccounts[owner][accountNumber]
	return m.ownerAccounts[owner][accountNumber], nil
}

func (m *InMemory) Delete(owner string, number int64) error {
	if _, ok := m.ownerAccounts[owner][number]; !ok {
		return &NotFound{}
	}
	delete(m.accounts, number)
	delete(m.ownerAccounts[owner], number)
	return nil
}

func (m *InMemory) unAssignedAccountNumber() int64 {
	for i := int64(0); i < int64(math.MaxInt64); i++ {
		if _, ok := m.accounts[i]; !ok {
			return i
		}
	}
	// Obviously don't actually panic in a real life scenario...
	panic("we have run out of account numbers")
}
