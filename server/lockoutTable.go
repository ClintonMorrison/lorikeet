package main

import "time"

const errorWindow = time.Hour * 5
const maxErrorsInWindow = 5

type LockoutTable struct {
	errorTimesByIP map[string][]time.Time
	errorTimesByUser map[string][]time.Time
}

func NewLockoutTable() *LockoutTable {
	return &LockoutTable{
		errorTimesByIP: make(map[string][]time.Time, 0),
		errorTimesByUser: make(map[string][]time.Time, 0)}
}

func (l *LockoutTable) shouldAllow(ip string, username string) bool {
	l.purgeErrors(ip, username)

	byIP := l.errorTimesByIP[ip]
	byUser := l.errorTimesByIP[username]

	if len(byIP) > maxErrorsInWindow || len(byUser) > maxErrorsInWindow {
		return false
	}

	return true
}

func (l *LockoutTable) logFailure(ip string, username string) {
	t := time.Now()

	if len(ip) > 0 {
		l.errorTimesByIP[ip] = append(l.errorTimesByIP[ip], t)
	}

	if len(username) > 0 {
		l.errorTimesByUser[username] = append(l.errorTimesByIP[username], t)
	}
}

func (l *LockoutTable) purgeErrors(ip string, username string) {
	earlistTime := time.Now().Add(errorWindow * -1)

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