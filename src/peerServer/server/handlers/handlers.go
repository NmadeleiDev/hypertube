package handlers

import (
	"net/http"
)

func PeerRequestsHandler(w http.ResponseWriter, r *http.Request) {
	SendSuccessResponse(w)
}
