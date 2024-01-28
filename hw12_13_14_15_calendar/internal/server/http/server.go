package internalhttp

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/common"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/domain"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/config"
)

type Server struct { // TODO
	cfg     *config.Config
	log     common.Logger
	event   domain.EventInter
	ln      net.Listener
	handler *http.ServeMux
}

func NewServer(log common.Logger, event domain.EventInter, cfg *config.Config) *Server {
	return &Server{
		cfg:   cfg,
		log:   log,
		event: event,
	}
}

func (s *Server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s.CreateMux()

	var err error
	s.ln, err = net.Listen("tcp4", s.cfg.Server.Host+":"+s.cfg.Server.Port)
	if err != nil {
		return err
	}

	server := http.Server{
		Handler: s.handler,
	}

	osSignals := make(chan os.Signal, 1)
	listenErr := make(chan error, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.log.Infof("web server started on %s", s.ln.Addr())
		listenErr <- server.Serve(s.ln)
	}()

	for {
		select {
		case err = <-listenErr:
			return err
		case <-osSignals:
			s.log.Info("web server stopped")
			if err := server.Shutdown(ctx); err != nil {
				return err
			}
			s.log.Info("web server exited properly")
		}
	}
}
