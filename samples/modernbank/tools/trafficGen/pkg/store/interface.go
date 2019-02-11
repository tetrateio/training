package store

type Interface interface {
	AddUser(username string)
	AddAccount(username string, account int64)
	GetRandomUser() string
	GetRandomAccount() int64
	UserCount() int64
	AccountCount() int64
}
