package util

import (
	"github.com/getsentry/raven-go"
	"fmt"
	"net/http"
	"net/url"
)

type HandlerParams struct {
	writ http.ResponseWriter
	req *http.Request
	ps url.Values
}

func ErrorHandler(err error){

	raven.CaptureErrorAndWait(err, nil)

	raven.CapturePanic(func() {

		fmt.Println(err.Error())

	}, nil)
}
