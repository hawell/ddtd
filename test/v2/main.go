package main

import (
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

func main() {
	// Start the tracer and defer the Stop method.
	tracer.Start(tracer.WithAgentAddr("localhost:8866"))
	defer tracer.Stop()

	// Start a root span.
	span := tracer.StartSpan("main")
	span.SetTag("ptag1", 11)
	span.SetTag("ptag2", "onetwo")

	child := span.StartChild("child1")
	child.SetTag("tag1", 1)
	child.SetTag("tag2", "two")
	time.Sleep(time.Second)
	child.Finish()

	child2 := span.StartChild("child2")
	child2.SetTag("tag21", 21)
	child2.SetTag("tag22", "twotwo")
	time.Sleep(time.Second)
	child2.Finish()

	span.Finish()
	tracer.Flush()

	time.Sleep(time.Second * 10)
}
