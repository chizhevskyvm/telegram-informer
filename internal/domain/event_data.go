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
	return utils.FromTimeZone(e.data[keyDate])
}
func (e *EventData) SetDate(t time.Time, userTZ string) {
	loc, err := time.LoadLocation(userTZ)
	if err != nil {
		loc = time.UTC
	}

	localTime := time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc,
	)
	utcTime := localTime.UTC()

	e.data[keyDate] = utils.FormatDate(utcTime)
}

func (e *EventData) SetTime(time time.Time) {
	e.data[keyTime] = utils.FormatTime(time)
}
func (e *EventData) GetTime() (time.Time, error) {
	return utils.ParseTime(e.data[keyTime])
}
