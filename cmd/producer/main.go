package main

import (
	"Jobqueue/internal/redisclient"
	"Jobqueue/internal/task"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func enqueueHandler(rdb *redisclient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		var t task.Task

		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "decoder error", http.StatusBadRequest)
			return
		}
		if t.Type == "" {
			http.Error(w, "type is required", http.StatusBadRequest)
			return
		}
		if t.Retries < 0 {
			http.Error(w, "retries > 0 ", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := rdb.PushTask(ctx, t); err != nil {
			http.Error(w, "failed to enqueue", http.StatusInternalServerError)
			log.Println("Redis err", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "queued",
			"type":   t.Type,
		})
		log.Println("task enqueued:", t.Type)
	}
}

func main() {
	rdb := redisclient.New()
	http.HandleFunc("/enqueue", enqueueHandler(rdb))
	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal(err)
	}

}
