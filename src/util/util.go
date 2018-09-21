package util

import (
	"github.com/getsentry/raven-go"
	"fmt"
	"log"
)

func ErrorHandler(err error){

	raven.CaptureErrorAndWait(err, nil)

	log.Panic(err)

	raven.CapturePanic(func() {

		fmt.Println(err.Error())

	}, nil)
}