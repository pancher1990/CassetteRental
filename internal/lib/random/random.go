package random

import (
	"github.com/google/uuid"
)

func GenerateGUID() string {
	id := uuid.New()
	return id.String()
}
