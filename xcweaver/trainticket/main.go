//go:generate xcweaver generate ./pkg/wrk2 . ./pkg/services ./pkg/model ./pkg/metrics

package main

import (
	"context"
	"log"

	"trainticket/pkg/wrk2"

	"github.com/TiagoMalhadas/xcweaver"
)

// this is an entry file for socialnetwork application
// the source code of services is in the "pkg" folder
func main() {
	if err := xcweaver.Run(context.Background(), wrk2.Serve); err != nil {
		log.Fatal(err)
	}
}
