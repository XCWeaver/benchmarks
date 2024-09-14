package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ServiceWeaver/weaver"
)

func main() {
	if err := weaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	notifier          weaver.Ref[Notifier]
	post_storage      weaver.Ref[Post_storage_america]
	post_notification weaver.Listener
}

// serve is called by weaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {

	logger := app.Logger(ctx)
	logger.Info("post-notification listener available", "address", app.post_notification)

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
	return http.Serve(app.post_notification, nil)
}
