package utils

import (
	log "github.com/sirupsen/logrus"
)

func Logstest() {
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
}
