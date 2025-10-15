package httphandlers

import (
	"fmt"
	"strconv"
)

func parseID(idStr string) (int, error) {
	if idStr == "" {
		return 0, fmt.Errorf("id не может быть пустым")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("некорректный формат id")
	}

	if id <= 0 {
		return 0, fmt.Errorf("id должен быть положительным")
	}

	return id, nil
}
