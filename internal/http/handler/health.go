package handler

import (
	"fmt"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OrderPulse is healthy")
}
