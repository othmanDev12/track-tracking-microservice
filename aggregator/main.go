package main

import (
	"encoding/json"
	"flag"
	"github.com/track-tracking/types"
	"log"
	"net/http"
)

func main() {
	listenAddress := flag.String("listenAddress", ":3000", "the listen address of the Http Server")
	flag.Parse()
	store := NewMemoryStore()
	svc := NewAggregateService(store)
	makeHttpTransport(*listenAddress, svc)
}

func makeHttpTransport(listenAddress string, aggregator Aggregator) {
	http.HandleFunc("/aggregate", handleAggregate(aggregator))
	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleAggregate(aggregator Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			err := writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			if err != nil {
				return
			}
			return
		}
		if err := aggregator.AggregateDistance(distance); err != nil {
			err := writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			if err != nil {
				return
			}
			return
		}
	}
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
