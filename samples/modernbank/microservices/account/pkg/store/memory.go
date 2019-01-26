package store

import (
	"sync"
	"sync/atomic"

	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/model"
)

var _ Interface = NewInMemory()

func NewInMemory() *InMemory {
	return &InMemory{m: &sync.RWMutex{}, currentAccountNumber: 0}
}

type InMemory struct {
	m                    *sync.RWMutex
	ownerAccounts        map[string]*accounts
	currentAccountNumber int64
}

type accounts struct {
	m        *sync.RWMutex
	accounts map[int64]model.Account
}

func (a *accounts) add(number int64, account *model.Account) {
	a.m.Lock()
	defer a.m.Unlock()
	a.accounts[number] = *account
}

func (a *accounts) delete(number int64) error {
	a.m.Lock()
	defer a.m.Unlock()
	if _, found := a.accounts[number]; !found {
		return &NotFound{}
	}
	delete(a.accounts, number)
	return nil
}

func (a *accounts) get(number int64) (*model.Account, bool) {
	a.m.RLock()
	defer a.m.RUnlock()
	tmp, found := a.accounts[number]
	return &tmp, found
}

func (a *accounts) list() []*model.Account {
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
		return nil, &NotFound{}
	}
	return m.ownerAccounts[owner].list(), nil
}

func (m *InMemory) Get(owner string, number int64) (*model.Account, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	res, ok := m.ownerAccounts[owner]
	if !ok {
		return nil, &NotFound{}
	}
	account, ok := res.get(number)
	if !ok {
		return nil, &NotFound{}
	}
	return account, nil
}

func (m *InMemory) Create(owner string) (*model.Account, error) {
	m.m.Lock()
	defer m.m.Unlock()
	_, ok := m.ownerAccounts[owner]
	if !ok {
		m.ownerAccounts[owner] = &accounts{m: &sync.RWMutex{}, accounts: map[int64]model.Account{}}
	}
	newAccountNumber := m.unAssignedAccountNumber()
	newAccount := &model.Account{
		Balance: 0,
		Owner:   owner,
		Number:  newAccountNumber,
	}
	m.ownerAccounts[owner].add(newAccountNumber, newAccount)
	return newAccount, nil
}

func (m *InMemory) Delete(owner string, number int64) error {
	m.m.Lock()
	defer m.m.Unlock()
	_, ok := m.ownerAccounts[owner]
	if !ok {
		return &NotFound{}
	}
	return m.ownerAccounts[owner].delete(number)
}

func (m *InMemory) unAssignedAccountNumber() int64 {
	return atomic.AddInt64(&m.currentAccountNumber, 1)
}
