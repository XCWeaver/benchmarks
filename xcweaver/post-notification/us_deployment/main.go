package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/XCWeaver/xcweaver"
)

func main() {
	if err := xcweaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	xcweaver.Implements[xcweaver.Main]
	notifier     xcweaver.Ref[Notifier]
	post_storage xcweaver.Ref[PostStorageUS]
	postnot      xcweaver.Listener
}

// serve is called by xcweaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {

	logger := app.Logger(ctx)
	logger.Info("post-notification listener available", "address", app.postnot)

	go app.notifier.Get().ReadNotification(ctx)

	post_storage := app.post_storage.Get()

	// Serve the /consistency_window endpoint.
	http.HandleFunc("/consistency_window", func(w http.ResponseWriter, r *http.Request) {
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
