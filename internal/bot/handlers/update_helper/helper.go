package updatehelper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot/models"
)

func ParseCallbackID(update *models.Update) (int, error) {
	rawID := update.CallbackQuery.Data
	parts := strings.Split(rawID, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf(rawID)
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf(rawID)
	}

	return id, nil
}
