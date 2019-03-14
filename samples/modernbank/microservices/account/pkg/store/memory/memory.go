package memory

import (
	"sync"
	"sync/atomic"

	"github.com/go-openapi/swag"

	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/store"
)

// Enforce that InMemory matches the Store Interface
var _ store.Interface = &InMemory{}

func NewInMemory() *InMemory {
	return &InMemory{
		m:                    &sync.RWMutex{},
		currentAccountNumber: 0,
		ownerAccounts:        map[string]accounts{},
		accountsOwner:        map[int64]string{},
	}
}

type InMemory struct {
	m                    *sync.RWMutex
	ownerAccounts        map[string]accounts
	accountsOwner        map[int64]string
	currentAccountNumber int64
}

type accounts struct {
	m        *sync.RWMutex
	accounts map[int64]model.Account
}

func (a accounts) add(number int64, account *model.Account) {
	a.m.Lock()
	defer a.m.Unlock()
	a.accounts[number] = *account
}

func (a accounts) delete(number int64) error {
	a.m.Lock()
	defer a.m.Unlock()
	if _, found := a.accounts[number]; !found {
		return &store.NotFound{}
	}
	delete(a.accounts, number)
	return nil
}

func (a accounts) get(number int64) (*model.Account, bool) {
	a.m.RLock()
	defer a.m.RUnlock()
	tmp, found := a.accounts[number]
	return &tmp, found
}

func (a accounts) list() []*model.Account {
	a.m.RLock()
	defer a.m.RUnlock()
	res := make([]*model.Account, len(a.accounts))
	for _, val := range a.accounts {
		tmp := val
		res = append(res, &tmp)
	}
	return res
}

func (m *InMemory) List(owner string) ([]*model.Account, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	_, ok := m.ownerAccounts[owner]
	if !ok {
		return nil, &store.NotFound{}
	}
	return m.ownerAccounts[owner].list(), nil
}

func (m *InMemory) Get(owner string, number int64) (*model.Account, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	res, ok := m.ownerAccounts[owner]
	if !ok {
		return nil, &store.NotFound{}
	}
	account, ok := res.get(number)
	if !ok {
		return nil, &store.NotFound{}
	}
	return account, nil
}

func (m *InMemory) Create(owner string) (*model.Account, error) {
	m.m.Lock()
	defer m.m.Unlock()
	_, ok := m.ownerAccounts[owner]
	if !ok {
		m.ownerAccounts[owner] = accounts{m: &sync.RWMutex{}, accounts: map[int64]model.Account{}}
	}
	newAccountNumber := m.unAssignedAccountNumber()
	newAccount := &model.Account{
		Balance: swag.Float64(100),
		Owner:   swag.String(owner),
		Number:  swag.Int64(newAccountNumber),
	}
	m.ownerAccounts[owner].add(newAccountNumber, newAccount)
	m.accountsOwner[newAccountNumber] = owner
	return newAccount, nil
}

func (m *InMemory) Delete(owner string, number int64) error {
	m.m.Lock()
	defer m.m.Unlock()
	_, ok := m.ownerAccounts[owner]
	if !ok {
		return &store.NotFound{}
	}
	delete(m.accountsOwner, number)
	return m.ownerAccounts[owner].delete(number)
}

func (m *InMemory) UpdateBalance(number int64, deltaAmount float64) error {
	owner, ok := m.accountsOwner[number]
	if !ok {
		return &store.NotFound{}
	}
	account := m.ownerAccounts[owner].accounts[number]
	newBalance := *account.Balance + deltaAmount
	account.Balance = &newBalance
	m.ownerAccounts[owner].accounts[number] = account
	return nil
}

func (m *InMemory) unAssignedAccountNumber() int64 {
	return atomic.AddInt64(&m.currentAccountNumber, 1)
}
