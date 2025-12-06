package steps

import (
	"strings"

	"github.com/google/uuid"
)

func CheckSessionID(current string) bool {

	if !strings.HasPrefix(current, "PXID-") {
		return false
	}

	raw := current[5:]
	_, err := uuid.Parse(raw)
	return err == nil
}
