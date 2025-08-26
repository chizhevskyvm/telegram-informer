package domain

import (
	"telegram-informer/common/utils"
	"time"
)

const (
	keyTitle = "title_value"
	keyDate  = "date_value"
	keyTime  = "time_value"
)

type EventData struct {
	data map[string]string
}

func NewEventData(data map[string]string) *EventData {
	if data == nil {
		data = make(map[string]string)
	}
	return &EventData{data: data}
}

func (e *EventData) Raw() map[string]string {
	return e.data
}

func (e *EventData) GetTitle() string {
	return e.data[keyTitle]
}
func (e *EventData) SetTitle(title string) {
	e.data[keyTitle] = title
}

func (e *EventData) GetDate() (time.Time, error) {
	return utils.ParseDateLocal(e.data[keyDate])
}
func (e *EventData) SetDate(time time.Time) {
	e.data[keyDate] = utils.FormatDate(time)
}

func (e *EventData) SetTime(time time.Time) {
	e.data[keyTime] = utils.FormatTime(time)
}
func (e *EventData) GetTime() (time.Time, error) {
	return utils.ParseTime(e.data[keyTime])
}
