package web

import (
	"context"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	"github.com/julienschmidt/httprouter"
	"github.com/syllabix/kafkaless/producer"
)

type Server interface {
	Shutdown(context.Context) error
}

type server struct {
	weaver.Implements[Server]
	api      weaver.Listener
	producer weaver.Ref[producer.Service]
}

func (s *server) Init(ctx context.Context) error {
	router := httprouter.New()
	router.PUT("/emit", s.emit)
	router.GET("/healthz", s.healthCheck)

	s.Logger().Info("api server ready to accept new connections", "addr", s.api)
	go http.Serve(s.api, router)
	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	if err := s.api.Close(); err != nil {
		s.Logger().ErrorCtx(ctx, "failed to gracefully shutdown web server", "error", err)
		return err
	}
	return nil
}

func (s *server) emit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	evt := r.URL.Query().Get("event")
	if evt == "" {
		evt = "<event missing>"
	}

	err := s.producer.Get().EmitEvent(r.Context(), evt)
	if err != nil {
		http.Error(w, "oops... sorry about that", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) healthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}
