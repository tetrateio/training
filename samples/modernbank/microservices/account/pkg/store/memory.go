package store

import (
	"sync"
	"sync/atomic"

	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/model"
)

var _ Interface = NewInMemory()

func NewInMemory() *InMemory {
	return &InMemory{currentAccountNumber: 0}
}

type InMemory struct {
	ownerAccounts        sync.Map
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
	res, ok := m.ownerAccounts.Load(owner)
	if !ok {
		return nil, &NotFound{}
	}
	return res.(*accounts).list(), nil
}

func (m *InMemory) Get(owner string, number int64) (*model.Account, error) {
	accountRes, ok := m.ownerAccounts.Load(owner)
	if !ok {
		return nil, &NotFound{}
	}
	account, ok := accountRes.(*accounts).get(number)
	if !ok {
		return nil, &NotFound{}
	}
	return account, nil
}

func (m *InMemory) Create(owner string) (*model.Account, error) {
	newAccountNumber := m.unAssignedAccountNumber()
	newAccount := &model.Account{
		Balance: 0,
		Owner:   owner,
		Number:  newAccountNumber,
	}
	accountRes, ok := m.ownerAccounts.Load(owner)
	var newAccounts accounts
	if !ok {
		newAccounts = accounts{m: &sync.RWMutex{}, accounts: map[int64]model.Account{}}
	} else {
		newAccounts = accountRes.(accounts)
	}
	newAccounts.add(newAccountNumber, newAccount)
	m.ownerAccounts.Store(owner, newAccounts)
	return newAccount, nil
}

func (m *InMemory) Delete(owner string, number int64) error {
	accountRes, ok := m.ownerAccounts.Load(owner)
	if !ok {
		return &NotFound{}
	}
	return accountRes.(*accounts).delete(number)
}

func (m *InMemory) unAssignedAccountNumber() int64 {
	return atomic.AddInt64(&m.currentAccountNumber, 1)
}
