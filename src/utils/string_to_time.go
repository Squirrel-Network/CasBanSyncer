package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func StringToTime(duration, typeTime string) (time.Time, error) {
	typeTime = strings.ToLower(typeTime)
	if strings.HasSuffix(typeTime, "s") {
		typeTime = typeTime[:len(typeTime)-1]
	}
	intDur, err := strconv.Atoi(duration)
	if typeTime == "than" {
		return time.Now(), nil
	} else if err == nil {
		switch typeTime {
		case "minute":
			return time.Now().Add(time.Minute * time.Duration(-intDur)), nil
		case "hour":
			return time.Now().Add(time.Hour * time.Duration(-intDur)), nil
		case "day":
			return time.Now().AddDate(0, 0, -intDur), nil
		}
	}
	return time.Time{}, errors.New(fmt.Sprintf("Type of time \"%s\" not found with duration of \"%s\"", typeTime, duration))
}
