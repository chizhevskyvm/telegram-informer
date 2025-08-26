package handlers

// CB - Call back
// MTP - Match type prefix
const (
	CBSetCreateEventState  = "add-event"
	CBGetEventToday        = "today-events"
	CBAllEvents            = "all-events"
	CBDeleteAllEventsToday = "cancel-today"
	CBGetEventById         = "get-by-id:"
	CBDeleteEventById      = "delete-by-id:"

	MTPStart    = "/start"
	MTCAddEvent = "" //empty
)
