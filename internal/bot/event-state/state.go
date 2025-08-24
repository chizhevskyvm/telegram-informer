package eventsstate

import "fmt"

const (
	userStateKey          = "user_state:%d"
	userStateKeyWithStage = "user_state:%d:%s"
	userStateKeyWithInner = "user_state:%d:%s"
	userStateKeyWithData  = "%s:%s:data"

	StateAddEventTitle = "add_event:%d:title"
	StateAddEventDate  = "add_event:%d:date"
	StateAddEventTime  = "add_event:%d:time"
	StateAddEventDone  = "add_event:%d:done"
)

func GetUserStateKey(id int64) string {
	return fmt.Sprintf(userStateKey, id)
}
func GetUserStageState(stageState string, id int64) string {
	return fmt.Sprintf(userStateKeyWithStage, id, stageState)
}

func GetUserStateDataKey(prefix string, value string) string {
	return fmt.Sprintf(userStateKeyWithData, prefix, value)
}

func GetUserInnerState(innerStateTemplate string, id int64) string {
	return fmt.Sprintf(userStateKeyWithInner, id, innerStateTemplate)
}

func IsAddEventState(state string, userId int64) bool {
	return state != GetUserStageState(StateAddEventTitle, userId) &&
		state != GetUserStageState(StateAddEventDate, userId) &&
		state != GetUserStageState(StateAddEventTime, userId) &&
		state != GetUserStageState(StateAddEventDone, userId)
}
