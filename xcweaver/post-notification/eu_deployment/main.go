package main

import (
	"context"
	"log"
	"net/http"

	"github.com/XCWeaver/xcweaver"
)

func main() {
	if err := xcweaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	xcweaver.Implements[xcweaver.Main]
	post_upload xcweaver.Ref[PostUpload]
	postnot     xcweaver.Listener
}

// serve is called by xcweaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {
	logger := app.Logger(ctx)
	logger.Info("post-notification listener available", "address", app.postnot)

	post_upload := app.post_upload.Get()

	// Serve the /postnot endpoint.
	http.HandleFunc("/post_notification", func(w http.ResponseWriter, r *http.Request) {
		requests.Inc()
		post := r.URL.Query().Get("post")
		postId, err := post_upload.Post(ctx, post, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(postId))
	})
	return http.Serve(app.postnot, nil)
}
