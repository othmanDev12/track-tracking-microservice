package main

import (
	"github.com/sirupsen/logrus"
	"github.com/track-tracking/types"
	"time"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddlware(next Aggregator) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("Aggregate distance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}
