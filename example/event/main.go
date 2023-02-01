package main

import (
	"context"
	"golang.org/x/exp/event"
	"golang.org/x/exp/event/adapter/logfmt"
	"golang.org/x/exp/event/eventtest"
	"os"
)

func main() {
	ctx := event.WithExporter(context.Background(), event.NewExporter(logfmt.NewHandler(os.Stdout), eventtest.ExporterOptions()))
	event.Log(ctx, "my event", event.Int64("myInt", 6))
	event.Log(ctx, "error event", event.String("myString", "some string value"))
}
