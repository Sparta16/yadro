package event

import (
	"errors"
	"github.com/Sparta16/yadro/app/pkg/client"
	"github.com/Sparta16/yadro/app/pkg/table"
	"strconv"
	"strings"
	"time"
)

var (
	Price          int
	TableCount     int
	AvailableTable int
	CloseTime      time.Time
	OpenTime       time.Time
)

type Event struct {
	Time    time.Time
	Command string
	Client  string
	Table   int
}

func New(time time.Time, command, client string, table int) Event {
	return Event{
		Time:    time,
		Command: command,
		Client:  client,
		Table:   table,
	}
}

func (e Event) Event(
	clients map[string]*client.Client,
	table []table.Table,
	crowd *[]string,
	output *strings.Builder,

) error {
	var err error
	switch e.Command {
	case "1":
		err = e.NewClient(clients, output)
		if err != nil {
			return err
		}
	case "2":
		err = e.Sit(clients, table, output)
		if err != nil {
			return err
		}
	case "3":
		crowd, err = e.Wait(clients, crowd, output)
		if err != nil {
			return err
		}
	case "4":
		crowd, err = e.Leave(clients, table, crowd, output)
		if err != nil {
			return err
		}
	default:
		err = errors.New("Not format")
		return err //выйти из приложения
	}

	return nil
}

func (e Event) NewClient(
	clients map[string]*client.Client,
	output *strings.Builder,
) error {
	if e.Table != 0 {
		return errors.New("Not format")
	}
	if (OpenTime.Compare(e.Time) == 1) || (CloseTime.Compare(e.Time) == -1) {
		output.WriteString(e.Time.Format("15:04") + " 13 NotOpenYet" + "\r\n")
		return nil
	}

	if _, exists := clients[e.Client]; exists {
		output.WriteString(e.Time.Format("15:04") + " 13 YouShallNotPass" + "\r\n")
		return nil
	}

	clients[e.Client] = client.New()
	return nil
}

func (e Event) Sit(
	clients map[string]*client.Client,
	table []table.Table,
	output *strings.Builder,
) error {
	if e.Table == 0 {
		return errors.New("Not format")
	}
	e.Table--
	if table[e.Table].ClientName != "" {
		output.WriteString(e.Time.Format("15:04") + " 13 PlaceIsBusy" + "\r\n")
		return nil
	}

	if _, exists := clients[e.Client]; !exists {
		output.WriteString(e.Time.Format("15:04") + " 13 ClientUnknown" + "\r\n")
		return nil
	}

	clients[e.Client] = &client.Client{StartTime: e.Time, Table: e.Table, OnSitting: true}
	table[e.Table].ChangeClientName(e.Client)
	AvailableTable--
	return nil
}

func (e Event) Wait(
	clients map[string]*client.Client,
	crowd *[]string,
	output *strings.Builder,
) (*[]string, error) {
	if e.Table != 0 {
		return nil, errors.New("Not format")
	}

	if AvailableTable > 0 {
		output.WriteString(e.Time.Format("15:04") + " 13 ICanWaitNoLonger!" + "\r\n")
		return crowd, nil
	}

	if len(*crowd) == TableCount {
		output.WriteString(e.Time.Format("15:04") + " 11 " + e.Client + "\r\n")
		delete(clients, e.Client)
		return crowd, nil
	}

	*crowd = append(*crowd, e.Client)
	return crowd, nil
}

func (e Event) Leave(
	clients map[string]*client.Client,
	table []table.Table,
	crowd *[]string,
	output *strings.Builder,
) (*[]string, error) {
	if e.Table != 0 {
		return nil, errors.New("Not format")
	}
	if _, exists := clients[e.Client]; !exists {
		output.WriteString(e.Time.Format("15:04") + " 13 ClientUnknown" + "\r\n")
		return crowd, nil
	}

	c := clients[e.Client]
	freeTable, spentTime := c.Leave(e.Time)
	table[freeTable].ChangeValue(spentTime, Price)
	delete(clients, e.Client)

	if len(*crowd) > 0 {
		newTableClient := (*crowd)[0]
		*crowd = (*crowd)[1:]
		clients[newTableClient] = &client.Client{StartTime: e.Time, Table: freeTable, OnSitting: true}
		table[freeTable].ClientName = newTableClient
		output.WriteString(e.Time.Format("15:04") + " 12 " + newTableClient + " " + strconv.Itoa(freeTable+1) + "\r\n")
	} else {
		AvailableTable++
	}

	return crowd, nil
}
