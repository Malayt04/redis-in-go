package store

import "time"

type Entry struct {
	value     string
	expiresAt time.Time
}

type Store struct {
	data map[string]*Entry
}

func New() *Store {
	return &Store{
		data: make(map[string]*Entry),
	}
}

func (s *Store) Set(key, value string) {
	s.data[key] = &Entry{
		value:     value,
		expiresAt: time.Time{},
	}
}

func (s *Store) SetWithExpiry(key, value string, expiry time.Duration) {
	s.data[key] = &Entry{
		value:     value,
		expiresAt: time.Now().Add(expiry),
	}
}

func (s *Store) Get(key string) (string, bool) {
	entry, ok := s.data[key]
	if !ok {
		return "", false
	}

	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		delete(s.data, key)
		return "", false
	}

	return entry.value, true
}

func (s *Store) Exists(key string) bool {
	entry, ok := s.data[key]
	if !ok {
		return false
	}

	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		delete(s.data, key)
		return false
	}

	return true
}

func (s *Store) Delete(key string) {
	delete(s.data, key)
}

func (s *Store) Expire(key string, expiry time.Duration) bool {
	entry, ok := s.data[key]
	if !ok {
		return false
	}

	entry.expiresAt = time.Now().Add(expiry)
	return true
}

func (s *Store) TTL(key string) int {
	entry, ok := s.data[key]
	if !ok {
		return -2
	}

	if entry.expiresAt.IsZero() {
		return -1
	}

	remaining := time.Until(entry.expiresAt)
	if remaining <= 0 {
		delete(s.data, key)
		return -2
	}

	return int(remaining.Seconds())
}