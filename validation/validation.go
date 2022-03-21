package validation

import (
	"github.com/sirupsen/logrus"
)

func Headers(headers []string) {
	if len(headers)%2 != 0 {
		logrus.Panic("Headers must be provided in multiple of 2 where each pair represents <key, value>")
	}
}
