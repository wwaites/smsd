package main

import (
	"strings"
)

// FIXME This won't work well for NANPA...
type Dialplan string

func (d Dialplan) Normalise(number string) string {
	if strings.HasPrefix(number, "00") {
		number = strings.TrimPrefix(number, "00")
	}
	if strings.HasPrefix(number, "0") {
		number = string(d) + strings.TrimLeft(number, "0")
	}
	return number
}
