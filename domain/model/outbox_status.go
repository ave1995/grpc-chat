package model

type OutboxStatus int

const (
	Pending   OutboxStatus = 1
	Processed OutboxStatus = 2
	Failed    OutboxStatus = 3
)
