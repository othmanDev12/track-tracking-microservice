package main

import (
	"fmt"
	"github.com/track-tracking/types"
)

type Aggregator interface {
	AggregateDistance(distance types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewAggregateService(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

type Storer interface {
	Insert(distance types.Distance) error
}

func (s *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("aggregate distance and inserting into store")
	return s.store.Insert(distance)

}
