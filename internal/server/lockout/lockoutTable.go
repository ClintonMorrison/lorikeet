package lockout

import (
	"sync"
	"time"
)

const errorWindow = time.Hour * 5
const maxErrorsInWindow = 5

type Table struct {
	errorTimesByIP   map[string][]time.Time
	errorTimesByUser map[string][]time.Time
	mux              sync.RWMutex
}

func NewTable() *Table {
	return &Table{
		errorTimesByIP:   make(map[string][]time.Time, 0),
		errorTimesByUser: make(map[string][]time.Time, 0)}
}

func (l *Table) ShouldAllow(ip string, username string) bool {
	l.purgeErrors(ip, username)

	l.mux.RLock()
	defer l.mux.RUnlock()

	byIP := l.errorTimesByIP[ip]
	byUser := l.errorTimesByUser[username]

	if len(byIP) > maxErrorsInWindow || len(byUser) > maxErrorsInWindow {
		return false
	}

	return true
}

func (l *Table) LogFailure(ip string, username string) {
	t := time.Now()

	l.mux.Lock()
	defer l.mux.Unlock()

	if len(ip) > 0 {
		l.errorTimesByIP[ip] = append(l.errorTimesByIP[ip], t)
	}

	if len(username) > 0 {
		l.errorTimesByUser[username] = append(l.errorTimesByUser[username], t)
	}
}

func (l *Table) purgeErrors(ip string, username string) {
	earlistTime := time.Now().Add(errorWindow * -1)

	l.mux.Lock()
	defer l.mux.Unlock()

	// Errors by IP
	ipErrorTimes := l.errorTimesByIP[ip]
	filteredIpErrorTimes := make([]time.Time, 0)

	for _, t := range ipErrorTimes {
		if earlistTime.Before(t) {
			filteredIpErrorTimes = append(filteredIpErrorTimes, t)
		}
	}

	l.errorTimesByIP[ip] = filteredIpErrorTimes

	// Errors by username
	userErrorTimes := l.errorTimesByUser[username]
	filteredUserErrorTimes := make([]time.Time, 0)

	for _, t := range userErrorTimes {
		if earlistTime.Before(t) {
			filteredUserErrorTimes = append(filteredUserErrorTimes, t)
		}
	}

	l.errorTimesByUser[username] = filteredUserErrorTimes
}
