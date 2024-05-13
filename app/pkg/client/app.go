package client

import (
	"time"
)

type Client struct {
	StartTime time.Time
	Table     int
	OnSitting bool
}

func New() *Client {
	return &Client{
		StartTime: time.Date(0000, time.January, 1, 0, 0, 0, 0, time.UTC),
		Table:     0,
	}
}

func (client *Client) Leave(commandTime time.Time) (int, time.Duration) {
	freeTable := client.Table
	spentTime := commandTime.Sub(client.StartTime)
	return freeTable, spentTime
}
