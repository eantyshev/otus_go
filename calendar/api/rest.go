package api

import (
	"encoding/json"
	"fmt"
	"github.com/eantyshev/otus_go/calendar/logger"
	"github.com/eantyshev/otus_go/calendar/pkg/adapters"
	"github.com/eantyshev/otus_go/calendar/pkg/appointment"
	"github.com/eantyshev/otus_go/calendar/pkg/appointment/repository"
	"github.com/eantyshev/otus_go/calendar/pkg/models"
	"go.uber.org/zap"
	"time"

	"github.com/gorilla/mux"
	"net/http"
)

type responseBody struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func handleSuccess(w http.ResponseWriter, r *http.Request, result interface{}) {
	body := &responseBody{Result: result}
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		handleError(w, r, 500, "json encoding error")
	}
}

func handleError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	w.WriteHeader(code)
	body := &responseBody{Error: msg}
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		panic(err)
	}
}

func AccessLogMiddleware(l *zap.SugaredLogger, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Infof("request to %s from %s", r.RequestURI, r.RemoteAddr)
		h(w, r)
	}
}

type MyHandler struct {
	repo appointment.Repository
	l    *zap.SugaredLogger
}

func (mh *MyHandler) CUDEvent(w http.ResponseWriter, r *http.Request) {
	var (
		ap   *models.Appointment
		err  error
	)
	if err = r.ParseForm(); err != nil {
		handleError(w, r, 400, fmt.Sprintf("failed to parse form: %s", err))
		return
	}
	if ap, err = adapters.NewFormAppointment(r); err != nil {
		handleError(w, r, 400, fmt.Sprintf("bad appointment data: %s", err))
		return
	}
	switch action := mux.Vars(r)["action"]; action {
	case "create":
		if err := mh.repo.Store(ap); err != nil {
			handleError(w, r, 500, fmt.Sprintf("cannot store appointment: %s", err))
			return
		}
	case "update":
		if err := mh.repo.Update(ap); err != nil {
			handleError(w, r, 500, fmt.Sprintf("cannot update appointment: %s", err))
			return
		}
	case "delete":
		if err := mh.repo.Delete(ap.ID); err != nil {
			handleError(w, r, 500, fmt.Sprintf("cannot delete appointment: %s", err))
			return
		}
	default:
		handleError(w, r, 500, fmt.Sprintf("Action %s not implemented", action))
		return
	}
	handleSuccess(w, r, "success")
}

func (mh *MyHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	var (
		now, since time.Time
		aps        = make([]*models.Appointment, 0)
		err        error
	)
	now = time.Now()
	switch period := mux.Vars(r)["period"]; period {
	case "day":
		since = now.AddDate(0, 0, -1)
	case "week":
		since = now.AddDate(0, 0, -7)
	case "month":
		since = now.AddDate(0, -1, 0)
	default:
		panic(fmt.Sprintf("Unknown period: %s", period))
	}
	if aps, _, err = mh.repo.Fetch(since, -1); err != nil {
		handleError(w, r, 500, fmt.Sprintf("fetch failed: %s", err))
		return
	}
	handleSuccess(w, r, aps)
}

func Server(addrPort string) {
	repo := repository.NewMapRepo()
	handler := &MyHandler{
		repo: repo,
	}
	router := mux.NewRouter()
	router.HandleFunc("/{action:(?:create|update|delete)}_event",
		AccessLogMiddleware(logger.L, handler.CUDEvent)).Methods(http.MethodPost)
	router.HandleFunc("/events_for_{period:(?:day|week|month)}",
		AccessLogMiddleware(logger.L, handler.ListEvents)).Methods(http.MethodGet)

	server := &http.Server{
		Addr:           addrPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.L.Debugf("listening at %s", addrPort)
	panic(server.ListenAndServe())
}
