package main

import (
	"context"
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
	post_upload weaver.Ref[PostUpload]
	postnot     weaver.Listener
}

// serve is called by weaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {
	logger := app.Logger(ctx)
	logger.Info("post-notification listener available", "address", app.postnot)

	post_upload := app.post_upload.Get()

	// Serve the /post_notification endpoint.
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
