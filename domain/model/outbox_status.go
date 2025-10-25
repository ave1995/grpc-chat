package model

type OutboxStatus int

const (
	Pending   OutboxStatus = 1
	Processed OutboxStatus = 2
)
