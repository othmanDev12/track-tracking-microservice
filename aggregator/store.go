package main

import "github.com/track-tracking/types"

type MememoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MememoryStore {
	return &MememoryStore{
		data: make(map[int]float64),
	}
}

func (s *MememoryStore) Insert(distance types.Distance) error {
	s.data[distance.OBUId] += distance.Value
	return nil
}
