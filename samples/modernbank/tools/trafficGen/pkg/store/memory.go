package store

import "math/rand"

var _ Interface = NewInMemory()

func NewInMemory() *InMemory {
	return &InMemory{
		users:        map[string][]int64{},
		accountCount: 0,
	}
}

type InMemory struct {
	users        map[string][]int64
	accountCount int64
}

func (m *InMemory) AddUser(username string) {
	m.users[username] = []int64{}
}

func (m *InMemory) AddAccount(username string, account int64) {
	m.users[username] = append(m.users[username], account)
	m.accountCount++
}

func (m *InMemory) GetRandomUser() string {
	// Map iteration is random so returning first key is good enough
	for username := range m.users {
		return username
	}
	return ""
}
func (m *InMemory) GetRandomAccount() int64 {
	// Map iteration is random so returning first key is good enough
	for _, accounts := range m.users {
		if len(accounts) > 0 {
			return accounts[rand.Intn(len(accounts))]
		}
	}
	return 0
}

func (m *InMemory) UserCount() int64 {
	return int64(len(m.users))
}

func (m *InMemory) AccountCount() int64 {
	return m.accountCount
}
