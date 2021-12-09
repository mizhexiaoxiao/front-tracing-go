package common

type QueryArgs struct {
	Api        string `query:"api" validate:"required"`
	StartTime  string `query:"startTime" validate:"required,len=13"`
	FinishTime string `query:"finishTime" validate:"required,len=13"`
}

type BodySlice []BodyArgs
type BodyArgs struct {
	Api          string `body:"api" validate:"required"`
	StartTime    string `body:"startTime" validate:"required,len=13"`
	FinishTime   string `body:"finishTime" validate:"required,len=13"`
	TraceID      string `body:"traceid" validate:"required,len=32"`
	SpanID       string `body:"spanid" validate:"required,len=16"`
	ParentSpanID string `body:"parentspanid"`
	Sampled      string `body:"sampled" validate:"required"`
	RequestID    string `body:"requestid" validate:"required"`
}
