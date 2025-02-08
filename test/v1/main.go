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
	span.SetTag("ptag1", 11)
	span.SetTag("ptag2", "onetwo")

	child, _ := tracer.StartSpanFromContext(ctx, "child")
	child.SetTag("tag1", 1)
	child.SetTag("tag2", "two")
	time.Sleep(time.Second)
	child.Finish()

	child2, _ := tracer.StartSpanFromContext(ctx, "child2")
	child2.SetTag("tag21", 21)
	child2.SetTag("tag22", "twotwo")
	time.Sleep(time.Second)
	child2.Finish()

	span.Finish()
	tracer.Flush()

	time.Sleep(time.Second * 10)
}
