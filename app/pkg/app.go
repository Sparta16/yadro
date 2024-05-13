package app

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"yadro/app/pkg/client"
	"yadro/app/pkg/event"
	table2 "yadro/app/pkg/table"
)

func Club() string {
	clients := make(map[string]*client.Client)
	var table []table2.Table
	var crowd []string
	var output strings.Builder

	fileName := os.Args[1]

	input, err := os.ReadFile(fileName)
	if err != nil {
		return "Не удалось прочитать файл"
	}

	str := strings.Split(string(input), "\r\n")

	event.TableCount, err = strconv.Atoi(str[0])
	if err != nil || event.TableCount <= 0 {
		return str[0]
	}

	event.AvailableTable = event.TableCount
	for i := 0; i < event.TableCount; i++ {
		table = append(table, table2.Table{})
	}

	times := strings.Split(str[1], " ")

	if len(times) != 2 {
		return str[1]
	}

	event.OpenTime, err = time.Parse("15:04", times[0])
	if err != nil {
		return str[1]
	}
	output.WriteString(event.OpenTime.Format("15:04") + "\r\n")

	event.CloseTime, err = time.Parse("15:04", times[1])
	if err != nil {
		return str[1]
	}

	event.Price, err = strconv.Atoi(str[2])
	if err != nil || event.Price <= 0 {
		return str[2]
	}

	for i := 3; i < len(str); i++ {
		eventString := str[i]
		eventSplit := strings.Split(str[i], " ")

		if len(eventSplit) < 2 || len(eventSplit) > 4 {
			return eventString
		}

		commandTime, err := time.Parse("15:04", eventSplit[0])
		if err != nil {
			return eventString
		}

		output.WriteString(eventString + "\r\n")

		command := eventSplit[1]
		clientName := eventSplit[2]

		tableNumber := 0
		if len(eventSplit) == 4 {
			tableNumber, err = strconv.Atoi(eventSplit[3])
			if err != nil {
				return eventString
			}
			if tableNumber > event.TableCount {
				return eventString
			}
		}

		e := event.New(commandTime, command, clientName, tableNumber)
		err = e.Event(clients, table, &crowd, &output)
		if err != nil {
			return eventString
		}
	}
	remainingClients := make([]string, 0, len(clients))
	for k := range clients {
		remainingClients = append(remainingClients, k)
	}
	sort.Strings(remainingClients)

	for _, clientName := range remainingClients {
		c := clients[clientName]
		freeTable, spentTime := c.Leave(event.CloseTime)
		if c.OnSitting {
			table[freeTable].ChangeValue(spentTime, event.Price)
		}
		output.WriteString(event.CloseTime.Format("15:04") + " 11 " + clientName + "\r\n")
	}

	output.WriteString(event.CloseTime.Format("15:04") + "\r\n")
	for i, value := range table {
		output.WriteString(strconv.Itoa(i+1) + " " + strconv.Itoa(value.Profit) + " " + value.WorkTime.Format("15:04") + "\r\n")
	}

	return output.String()
}
