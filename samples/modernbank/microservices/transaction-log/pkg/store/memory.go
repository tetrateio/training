package store

import (
	"sync"
	"sync/atomic"

	"github.com/go-openapi/swag"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/model"
)

// Enforce that InMemory matches the Store Interface
var _ Interface = NewInMemory()

func NewInMemory() *InMemory {
	return &InMemory{
		m:        &sync.RWMutex{},
		sent:     map[int64]transactions{},
		received: map[int64]transactions{},
	}
}

type InMemory struct {
	m                        *sync.RWMutex
	sent                     map[int64]transactions
	received                 map[int64]transactions
	currentTransactionNumber int64
}

type transactions struct {
	m *sync.RWMutex
	// Transactions are immutable
	// Because this is a demo app we will just duplicate in memory.
	transactions map[int64]model.Transaction
}

func (t transactions) add(number int64, transaction *model.Transaction) {
	t.m.Lock()
	defer t.m.Unlock()
	t.transactions[number] = *transaction
}

func (t transactions) get(number int64) (*model.Transaction, bool) {
	t.m.RLock()
	defer t.m.RUnlock()
	tmp, found := t.transactions[number]
	return &tmp, found
}

func (t transactions) list() []*model.Transaction {
	t.m.RLock()
	defer t.m.RUnlock()
	res := make([]*model.Transaction, len(t.transactions))
	for _, val := range t.transactions {
		tmp := val
		res = append(res, &tmp)
	}
	return res
}

func (m *InMemory) ListSent(account int64) ([]*model.Transaction, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	_, ok := m.sent[account]
	if !ok {
		return nil, &NotFound{}
	}
	return m.sent[account].list(), nil
}

func (m *InMemory) ListReceived(account int64) ([]*model.Transaction, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	_, ok := m.received[account]
	if !ok {
		return nil, &NotFound{}
	}
	return m.received[account].list(), nil
}

func (m *InMemory) GetSent(account int64, number int64) (*model.Transaction, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	res, ok := m.sent[account]
	if !ok {
		return nil, &NotFound{}
	}
	transaction, ok := res.get(number)
	if !ok {
		return nil, &NotFound{}
	}
	return transaction, nil
}

func (m *InMemory) GetReceived(account int64, number int64) (*model.Transaction, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	res, ok := m.received[account]
	if !ok {
		return nil, &NotFound{}
	}
	transaction, ok := res.get(number)
	if !ok {
		return nil, &NotFound{}
	}
	return transaction, nil
}

func (m *InMemory) Create(transaction *model.Newtransaction) (*model.Transaction, error) {
	m.m.Lock()
	defer m.m.Unlock()

	if _, ok := m.sent[*transaction.Sender]; !ok {
		m.sent[*transaction.Sender] = transactions{m: &sync.RWMutex{}, transactions: map[int64]model.Transaction{}}
	}
	if _, ok := m.received[*transaction.Receiver]; !ok {
		m.received[*transaction.Receiver] = transactions{m: &sync.RWMutex{}, transactions: map[int64]model.Transaction{}}
	}

	newTransactionNumber := m.unAssignedTransactionNumber()
	newTransaction := &model.Transaction{transaction.Amount, transaction.Receiver, transaction.Sender, swag.Int64(newTransactionNumber)}
	m.sent[*transaction.Sender].add(newTransactionNumber, newTransaction)
	m.received[*transaction.Receiver].add(newTransactionNumber, newTransaction)
	return newTransaction, nil
}

func (m *InMemory) unAssignedTransactionNumber() int64 {
	return atomic.AddInt64(&m.currentTransactionNumber, 1)
}
