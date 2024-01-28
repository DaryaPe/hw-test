package internalhttp

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/common"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/pkg/util"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				log.Fatal(rvr)
			}

			s.log.Debugf(
				"%s [%s] %s %s %dms %s",
				r.RemoteAddr,
				time.Now().Format(time.RFC3339),
				r.Method,
				r.URL.String(),
				time.Unix(time.Now().Unix()-r.Context().Value(common.CtxLatency).(int64), 0).UnixMilli(),
				r.UserAgent(),
			)
		}()

		actionID, err := util.GenerateUUID()
		if err != nil {
			panic(err)
		}
		ctx := context.WithValue(context.Background(), common.CtxActionID, actionID) // nolint:staticcheck
		ctx = context.WithValue(ctx, common.CtxLatency, time.Now().Unix())           // nolint:staticcheck
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
