package models

import "errors"

// ErrNoRecord is returned when there is no matching record in the database
var ErrNoRecord = errors.New("models: no matching record found")
