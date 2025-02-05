package main

import (
	"context"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	tracer.Start(tracer.WithAgentAddr("localhost:8866"))
	defer tracer.Stop()

	span, ctx := tracer.StartSpanFromContext(context.Background(), "main")

	child, _ := tracer.StartSpanFromContext(ctx, "child")
	time.Sleep(time.Second)
	child.Finish()

	child2, _ := tracer.StartSpanFromContext(ctx, "child2")
	time.Sleep(time.Second)
	child2.Finish()

	span.Finish()
	tracer.Flush()

	time.Sleep(time.Second * 10)
}
