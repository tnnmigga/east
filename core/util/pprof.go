package util

import (
	"net/http"
	_ "net/http/pprof"
)

func PProf() {
	http.ListenAndServe("0.0.0.0:8080", nil)
}
