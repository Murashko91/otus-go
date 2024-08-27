package internalhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Handler struct {
	app app.Application
}

func (h Handler) eventHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetEvents(w, r, h.app)
	case http.MethodPost:
		handleInsertEvent(w, r, h.app)
	case http.MethodPatch:
		handleUpdateEvent(w, r, h.app)
	case http.MethodDelete:
		handleDeleteEvents(w, r, h.app)
	}
}

func handleGetEvents(w http.ResponseWriter, r *http.Request, a app.Application) {
	userID := r.URL.Query().Get(app.UserIDKey)
	startDate := r.URL.Query().Get("start_date")
	duration := r.URL.Query().Get("duration")

	dt, err := time.Parse(time.DateTime, startDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("missed or incorrect format start date, should be: %s \n", time.DateTime)))

		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missed or incorrect user id\n"))
		return
	}

	ctx := app.SetContextValue(r.Context(), app.UserIDKey, uid)

	events := make([]storage.Event, 0)

	switch duration {
	case "day":
		events, err = a.GetDailyEvents(ctx, dt)
	case "week":
		events, err = a.GetWeeklyEvents(ctx, dt)
	case "month":
		events, err = a.GetMonthlyEvents(ctx, dt)
	default:
		w.Write([]byte("missed or incorrect duration, should be one of (day, week, month)\n"))
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)

	jsonData, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.Write(jsonData)
}

func handleInsertEvent(w http.ResponseWriter, r *http.Request, a app.Application) {
	handleAlterEvents(w, r, a, "insert")
}

func handleUpdateEvent(w http.ResponseWriter, r *http.Request, a app.Application) {
	handleAlterEvents(w, r, a, "update")
}

func handleDeleteEvents(w http.ResponseWriter, r *http.Request, a app.Application) {
	handleAlterEvents(w, r, a, "delete")
}

func handleAlterEvents(w http.ResponseWriter, r *http.Request, a app.Application, cmd string) {
	userID := r.URL.Query().Get(app.UserIDKey)

	uid, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("missed or incorrect format start date, should be: %s \n", time.DateTime)))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var event storage.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ctx := app.SetContextValue(r.Context(), app.UserIDKey, uid)
	switch cmd {
	case "insert":
		event, err := a.CreateEvent(ctx, event)
		setInsertUpdateEventResponse(event, err, w)
	case "update":
		event, err := a.UpdateEvent(ctx, event)
		setInsertUpdateEventResponse(event, err, w)
	case "delete":
		err := a.DeleteEvent(ctx, event.ID)
		setDeleteEventResponse(err, w)
	}
}

func setInsertUpdateEventResponse(event storage.Event, err error, w http.ResponseWriter) {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(jsonResponse)
}

func setDeleteEventResponse(err error, w http.ResponseWriter) {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
