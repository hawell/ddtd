package types

import "context"

type Trace = struct {
	Type     string            `bson:"type" json:"type"`
	Resource string            `bson:"resource" json:"resource"`
	Service  string            `bson:"service" json:"service"`
	Start    int64             `bson:"start" json:"start"`
	Duration int64             `bson:"duration" json:"duration"`
	Metrics  map[string]int    `bson:"metrics" json:"metrics"`
	Meta     map[string]string `bson:"meta" json:"meta"`
	SpanID   int64             `bson:"span_id" json:"span_id"`
	TraceID  int64             `bson:"trace_id" json:"trace_id"`
	ParentID int64             `bson:"parent_id" json:"parent_id"`
	Error    int               `bson:"error" json:"error"`
	Name     string            `bson:"name" json:"name"`
	Level    int               `json:"level"`
}

type Storage interface {
	AddTrace(ctx context.Context, b []byte) error
	GetTraces(ctx context.Context) ([]Trace, error)
	Clear(ctx context.Context)
}
