package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: No matching record was found")
