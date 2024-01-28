package internalhttp

import (
	"net/http"

	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/common"
)

func (s *Server) defaultRoute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello-world, action-id=" + r.Context().Value(common.CtxActionID).(string)))
	})
}
