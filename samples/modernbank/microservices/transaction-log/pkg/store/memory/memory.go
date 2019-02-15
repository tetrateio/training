package memory

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/go-openapi/swag"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/store"
)

// Enforce that InMemory matches the Store Interface
var _ store.Interface = &InMemory{}

func NewInMemory() *InMemory {
	return &InMemory{
		m:        &sync.RWMutex{},
		sent:     map[int64]transactions{},
		received: map[int64]transactions{},
	}
}

type InMemory struct {
	m        *sync.RWMutex
	sent     map[int64]transactions
	received map[int64]transactions
}

type transactions struct {
	m *sync.RWMutex
	// Transactions are immutable
	// Because this is a demo app we will just duplicate in memory.
	transactions map[string]model.Transaction
}

func (t transactions) add(id string, transaction *model.Transaction) {
	t.m.Lock()
	defer t.m.Unlock()
	t.transactions[id] = *transaction
}

func (t transactions) get(id string) (*model.Transaction, bool) {
	t.m.RLock()
	defer t.m.RUnlock()
	tmp, found := t.transactions[id]
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
		return nil, &store.NotFound{}
	}
	return m.sent[account].list(), nil
}

func (m *InMemory) ListReceived(account int64) ([]*model.Transaction, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	_, ok := m.received[account]
	if !ok {
		return nil, &store.NotFound{}
	}
	return m.received[account].list(), nil
}

func (m *InMemory) GetSent(account int64, id string) (*model.Transaction, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	res, ok := m.sent[account]
	if !ok {
		return nil, &store.NotFound{}
	}
	transaction, ok := res.get(id)
	if !ok {
		return nil, &store.NotFound{}
	}
	return transaction, nil
}

func (m *InMemory) GetReceived(account int64, id string) (*model.Transaction, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	res, ok := m.received[account]
	if !ok {
		return nil, &store.NotFound{}
	}
	transaction, ok := res.get(id)
	if !ok {
		return nil, &store.NotFound{}
	}
	return transaction, nil
}

func (m *InMemory) Create(transaction *model.Newtransaction) (*model.Transaction, error) {
	m.m.Lock()
	defer m.m.Unlock()

	if _, ok := m.sent[*transaction.Sender]; !ok {
		m.sent[*transaction.Sender] = transactions{m: &sync.RWMutex{}, transactions: map[string]model.Transaction{}}
	}
	if _, ok := m.received[*transaction.Receiver]; !ok {
		m.received[*transaction.Receiver] = transactions{m: &sync.RWMutex{}, transactions: map[string]model.Transaction{}}
	}

	newTransactionID := m.genID(transaction)
	newTransaction := &model.Transaction{transaction.Amount, transaction.Receiver, transaction.Sender, swag.String(newTransactionID)}
	m.sent[*transaction.Sender].add(newTransactionID, newTransaction)
	m.received[*transaction.Receiver].add(newTransactionID, newTransaction)
	return newTransaction, nil
}

func (m *InMemory) genID(transaction *model.Newtransaction) string {
	h := md5.New()
	timeSalt := strconv.FormatInt(time.Now().UnixNano(), 10)
	io.WriteString(h, timeSalt)
	bytes, _ := json.Marshal(transaction)
	h.Write(bytes)
	return fmt.Sprintf("%x", h.Sum(nil))
}
