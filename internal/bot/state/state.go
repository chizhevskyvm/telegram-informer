package state

import "fmt"

type State string

const (
	AddEvent = "add_event"

	userStateKey          = "user_state:%d"    // user_state:<id>
	userStateKeyWithStage = "user_state:%s:%d" // user_state:<stage>:<id>
	userStateKeyWithData  = "%s:%s:data"       // <prefix>:<key>:data
)

const (
	addEventTitle State = AddEvent + ":title"
	addEventDate  State = AddEvent + ":date"
	addEventTime  State = AddEvent + ":time"
	addEventDone  State = AddEvent + ":done"
)

func UserStateKey(userID int64) string {
	return fmt.Sprintf(userStateKey, userID)
}

func StageKey(stage State, userID int64) string {
	return fmt.Sprintf(userStateKeyWithStage, stage, userID)
}

func DataKey(prefix, value string) string {
	return fmt.Sprintf(userStateKeyWithData, prefix, value)
}

func CreateEventState(userID int64) State { return State(StageKey(AddEvent, userID)) }

func AddEventTitleState(userID int64) State { return State(StageKey(addEventTitle, userID)) }
func AddEventDateState(userID int64) State  { return State(StageKey(addEventDate, userID)) }
func AddEventTimeState(userID int64) State  { return State(StageKey(addEventTime, userID)) }
func AddEventDoneState(userID int64) State  { return State(StageKey(addEventDone, userID)) }

func IsAddEventState(stateKey State, userID int64) bool {
	return stateKey == AddEventTitleState(userID) ||
		stateKey == AddEventDateState(userID) ||
		stateKey == AddEventTimeState(userID) ||
		stateKey == AddEventDoneState(userID)
}
