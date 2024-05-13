package table

import (
	"math"
	"time"
)

type Table struct {
	ClientName string
	Profit     int
	WorkTime   time.Time
}

func (table *Table) ChangeValue(spentTime time.Duration, price int) {
	table.ClientName = ""
	table.WorkTime = table.WorkTime.Add(spentTime)
	table.Profit += price * int(math.Ceil(spentTime.Hours()))
}

func (table *Table) ChangeClientName(clientName string) {
	table.ClientName = clientName
}
