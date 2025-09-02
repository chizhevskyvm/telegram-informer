package handlers

// CB - Call back
// MTP - Match type prefix
const (
	CBSetCreateEventState  = "set_create_event_state"
	CBGetEventToday        = "get_event_today"
	CBGetEventsActual      = "get_events_actual"
	CBDeleteAllEventsToday = "delete_all_events_today"
	CBGetEventById         = "get_event_by_id:"
	CBDeleteEventById      = "delete_event_by_id:"

	MTPStart    = "/start"
	MTCAddEvent = "" //empty
)
