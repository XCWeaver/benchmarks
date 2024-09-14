package main

import (
	"context"
	"log"
	"net/http"

	"github.com/TiagoMalhadas/xcweaver"
)

func main() {
	if err := xcweaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	xcweaver.Implements[xcweaver.Main]
	post_upload       xcweaver.Ref[Post_upload]
	post_notification xcweaver.Listener
}

// serve is called by xcweaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {
	logger := app.Logger(ctx)
	logger.Info("post-notification listener available", "address", app.post_notification)

	post_upload := app.post_upload.Get()

	// Serve the /post_notification endpoint.
	http.HandleFunc("/post_notification", func(w http.ResponseWriter, r *http.Request) {
		requests.Inc()
		post := r.URL.Query().Get("post")
		err := post_upload.Post(ctx, post, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return http.Serve(app.post_notification, nil)
}
