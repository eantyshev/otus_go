package api

import (
	"encoding/json"
	"fmt"
	"github.com/eantyshev/otus_go/calendar/logger"
	"github.com/eantyshev/otus_go/calendar/pkg/appointment"
	"github.com/eantyshev/otus_go/calendar/pkg/appointment/repository"
	"github.com/eantyshev/otus_go/calendar/pkg/models"
	"go.uber.org/zap"
	"strconv"
	"time"

	"net/http"
)

func NewFormAppointment(r *http.Request) (ap *models.Appointment, err error) {
	var (
		id int
		startsAt time.Time
		durationMinutes int
		isRegular bool
	)
	if err = r.ParseForm(); err != nil {
		return nil, err
	}
	if id, err = strconv.Atoi(r.FormValue("id")); err != nil {
		return nil, err
	}
	if startsAt, err = time.Parse(time.RFC3339, r.FormValue("starts_at")); err != nil {
		return nil, err
	}
	if durationMinutes, err = strconv.Atoi(r.FormValue("duration_minutes")); err != nil {
		return nil, err
	}
	if isRegular, err = strconv.ParseBool(r.FormValue("is_regular")); err != nil {
		return nil, err
	}
	ap = &models.Appointment{
		ID:              int64(id),
		Summary:         r.FormValue("summary"),
		Description:     r.FormValue("description"),
		StartsAt:        startsAt,
		DurationMinutes: uint16(durationMinutes),
		IsRegular:       isRegular,
		DaysOfWeek:      nil,
	}
	return ap, nil
}


type successBody struct {
	Result string `json:"result"`
}

func handleSuccess(w http.ResponseWriter, r *http.Request, result string) {
	body := &successBody{Result: result}
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		handleError(w, r, 500, "json encoding error")
	}
}

type errorBody struct {
	Error string `json:"error"`
}

func handleError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	w.WriteHeader(code)
	body := &errorBody{Error: msg}
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		panic(err)
	}
}

func AccessLogMiddleware(l *zap.SugaredLogger, h http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		l.Infof("request to %s from %s", r.RequestURI, r.RemoteAddr)
		h(w, r)
	}
}

type MyHandler struct {
	repo appointment.Repository
	l    *zap.SugaredLogger
}


func (mh *MyHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var (
		ap *models.Appointment
		err error
	)
	if err = r.ParseForm(); err != nil {
		handleError(w, r, 400, fmt.Sprintf("failed to parse form: %s", err))
		return
	}
	if ap, err = NewFormAppointment(r); err != nil {
		handleError(w, r, 400, fmt.Sprintf("failed to construct the appointment: %s", err))
		return
	}
	if err := mh.repo.Store(ap); err != nil {
		handleError(w, r, 500, fmt.Sprintf("cannot store appointment: %s", err))
		return
	}
	handleSuccess(w, r, "stored successfully")
}

func Server(addrPort string) {
	repo := repository.NewMapRepo()
	handler := &MyHandler{
		repo: repo,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", AccessLogMiddleware(logger.L, handler.CreateEvent))
	server := &http.Server{
		Addr:           addrPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.L.Debugf("listening at %s", addrPort)
	panic(server.ListenAndServe())
}
