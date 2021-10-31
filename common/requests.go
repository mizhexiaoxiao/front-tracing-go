package common

type QueryArgs struct {
	Api        string `query:"api" validate:"required"`
	StartTime  string `query:"startTime" validate:"required,len=19"`
	FinishTime string `query:"finishTime" validate:"required,len=19"`
}
