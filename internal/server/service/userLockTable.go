package service

import (
	"sync"
)

type UserLockTable struct {
	lockByUser map[string]*sync.RWMutex
	lockMux    sync.RWMutex
}

func NewUserLockTable() *UserLockTable {
	return &UserLockTable{
		lockByUser: make(map[string]*sync.RWMutex),
	}
}

func (t *UserLockTable) getLockForUser(username string) *sync.RWMutex {
	if t.lockByUser[username] == nil {
		t.lockByUser[username] = &sync.RWMutex{}
	}

	return t.lockByUser[username]
}

func (t *UserLockTable) Lock(username string) {
	t.lockMux.Lock()
	defer t.lockMux.Unlock()

	lock := t.getLockForUser(username)
	lock.Lock()
}

func (t *UserLockTable) Unlock(username string) {
	t.lockMux.Lock()
	defer t.lockMux.Unlock()

	lock := t.getLockForUser(username)
	lock.Unlock()
}
