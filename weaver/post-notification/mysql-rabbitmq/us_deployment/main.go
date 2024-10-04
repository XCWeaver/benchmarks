package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ServiceWeaver/weaver"
)

func main() {
	if err := weaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	notifier     weaver.Ref[Notifier]
	post_storage weaver.Ref[PostStorageUs]
	postnot      weaver.Listener
}

// serve is called by weaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {

	logger := app.Logger(ctx)
	logger.Info("post-notification listener available", "address", app.postnot)

	go app.notifier.Get().ReadNotification(ctx)

	post_storage := app.post_storage.Get()

	// Serve the /post_notification endpoint.
	http.HandleFunc("/consistency_window", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("request")
		values, err := post_storage.GetConsistencyWindowValues(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			jsonData, err := json.Marshal(values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		}
	})
	// Serve the /inconsistencies endpoint.
	http.HandleFunc("/inconsistencies", func(w http.ResponseWriter, r *http.Request) {
		result, err := post_storage.GetInconsistencies(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "text/plain")
			resultStr := strconv.Itoa(result)

			w.Write([]byte(resultStr))
		}
	})
	// Serve the /reset endpoint.
	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		post_storage.Reset(ctx)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("done!"))
	})
	return http.Serve(app.postnot, nil)
}
