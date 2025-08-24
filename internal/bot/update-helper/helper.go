package updatehelper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot/models"
)

func GetId(update *models.Update) (int, error) {
	rawID := update.CallbackQuery.Data // "get-by-id:13"
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
