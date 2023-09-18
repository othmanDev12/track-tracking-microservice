package main

import (
	"fmt"
	"github.com/track-tracking/types"
)

type Aggregator interface {
	AggregateDistance(distance types.Distance) error
}

type AggregateService struct {
	Storer Storer
}

type Storer interface {
	Insert(distance types.Distance) error
}

func NewAggregateService(storer Storer) *AggregateService {
	return &AggregateService{
		Storer: storer,
	}
}

func (s *AggregateService) AggregateDistance(distance types.Distance) error {
	fmt.Println("aggregate distance ............. ")
	return s.Storer.Insert(distance)

}
