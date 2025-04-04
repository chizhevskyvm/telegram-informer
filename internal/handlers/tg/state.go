package tg

import "fmt"

const (
	userStateKey          = "user_state:%d"
	userStateKeyWithStage = "user_state:%d:%s"
	userStateKeyWithInner = "user_state:%d:%s"
	userStateKeyWithData  = "%s:%s:data"
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
