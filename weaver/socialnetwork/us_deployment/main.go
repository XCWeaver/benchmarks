//go:generate weaver generate . ./pkg/services ./pkg/model ./pkg/trace ./pkg/metrics

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"us_deployment/pkg/services"

	"github.com/ServiceWeaver/weaver"
)

func main() {
	if err := weaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	writeHomeTimelineService weaver.Ref[services.WriteHomeTimelineService]
	_                        weaver.Ref[services.UpdateHomeTimelineService]
	socialnetwork            weaver.Listener
}

// serve is called by weaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {

	logger := app.Logger(ctx)
	logger.Info("socialnetwork listener available", "address", app.socialnetwork)

	writeHomeTimelineService := app.writeHomeTimelineService.Get()

	go writeHomeTimelineService.ReadNotifications(ctx)

	// Serve the /post_notification endpoint.
	http.HandleFunc("/consistency_window", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("request")
		values, err := writeHomeTimelineService.GetConsistencyWindowValues(ctx)
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
	return http.Serve(app.socialnetwork, nil)

}
