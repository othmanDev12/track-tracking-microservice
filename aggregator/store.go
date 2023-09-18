package main

import "github.com/track-tracking/types"

type MememoryStore struct {
}

func NewMemoryStore() *MememoryStore {
	return &MememoryStore{}
}

func (s *MememoryStore) Insert(distance types.Distance) error {
	return nil
}
