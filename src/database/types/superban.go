package types

import (
	"time"
)

type Superban struct {
	UserId            string
	MotivationText    string
	UserDate          time.Time
	UserFirstName     string
	IdOperator        string
	UsernameOperator  string
	FirstNameOperator string
}
