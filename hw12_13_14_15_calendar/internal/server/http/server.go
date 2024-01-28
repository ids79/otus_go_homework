package internalhttp

import (
	"context"
	"encoding/json"
	"errors"

	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	internaljson "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/json"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	uuid "github.com/satori/go.uuid"
)

type Server struct {
	logg logger.Logg
	app  app.Application
	conf *config.Config
	srv  *http.Server
}

func NewServer(logger logger.Logg, app app.Application, config config.Config) *Server {
	return &Server{
		logg: logger,
		app:  app,
		conf: &config,
	}
}

func (s *Server) Start(ctx context.Context) error {
	handler := s.loggingMiddleware()
	server := &http.Server{
		Addr:         s.conf.HTTPServer.Address,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.srv = server
	s.logg.Info("starting http server on ", server.Addr)
	server.ListenAndServe()
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		s.logg.Error("HTTP server error: %v", err)
	}
	return nil
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	s.logg.Info("server http is stopping...")
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logg.Error("failed to stop http server: " + err.Error())
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("hello-otus")
	case "/create/":
		s.create(w, r)
	case "/update/":
		s.update(w, r)
	case "/delete/":
		s.delete(w, r)
	case "/list-on-day/":
		s.list(w, r, "day")
	case "/list-on-week/":
		s.list(w, r, "week")
	case "/list-on-month/":
		s.list(w, r, "month")
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var ev internaljson.Event
		err := decoder.Decode(&ev)
		defer r.Body.Close()
		if err != nil {
			s.logg.Info(err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
		defer cancel()
		userID, _ := strconv.Atoi(ev.UserID)
		appEv := app.Event{
			Title:       ev.Title,
			DateTime:    time.Time(ev.DateTime),
			Duration:    time.Duration(ev.Duration),
			TimeBefore:  time.Duration(ev.TimeBefore),
			Description: ev.Description,
			UserID:      userID,
		}
		u := s.app.CreateEvent(ctx, appEv)
		if u != uuid.Nil {
			s.logg.Info("create new event with uuid: ", u.String())
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(u.String()))
		} else {
			http.Error(w, "error when adding an event", http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var ev internaljson.Event
		err := decoder.Decode(&ev)
		defer r.Body.Close()
		if err != nil {
			s.logg.Info(err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
		defer cancel()
		appEv := app.Event{
			Duration:    time.Duration(ev.Duration),
			TimeBefore:  time.Duration(ev.TimeBefore),
			Description: ev.Description,
		}
		err = s.app.UpgateEvent(ctx, ev.ID, appEv)
		if err == nil {
			s.logg.Info("update en event with uuid: ", ev.ID.String())
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(ev.ID.String()))
		} else {
			http.Error(w, "error when update an event", http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			s.logg.Error(err.Error())
			return
		}
		st := string(buf)
		ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
		defer cancel()
		err = s.app.DeleteEvent(ctx, uuid.FromStringOrNil(st))
		if err == nil {
			s.logg.Info("delete an event with uuid: ", st)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("delete an event with uuid: %s", st)))
		} else {
			http.Error(w, "error when delete an event", http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}
}

func (s *Server) list(w http.ResponseWriter, r *http.Request, period string) {
	if r.Method != http.MethodPost {
		return
	}
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		s.logg.Error(err.Error())
		return
	}
	st := string(buf)
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()
	day, err := time.Parse("2006-01-02", st)
	if err != nil {
		s.logg.Error(err.Error())
		return
	}
	var evApp []app.Event
	switch period {
	case "day":
		evApp = s.app.GetListOnDay(ctx, day)
	case "week":
		evApp = s.app.GetListOnWeek(ctx, day)
	case "month":
		evApp = s.app.GetListOnMonth(ctx, day)
	}
	if evApp == nil {
		http.Error(w, "error selecting events for the day", http.StatusInternalServerError)
		return
	}
	if len(evApp) == 0 {
		var resp string
		switch period {
		case "day":
			resp = fmt.Sprintf("nothing was selected on the day %s", day.Format("2006-01-02"))
		case "week":
			_, w := day.ISOWeek()
			resp = fmt.Sprintf("nothing was selected on the week %d", w)
		case "month":
			resp = fmt.Sprintf("nothing was selected on the month %s", day.Format("2006-01"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	} else {
		events := internaljson.EventsFormAppToView(evApp)
		body, err := json.Marshal(events)
		if err != nil {
			s.logg.Error(err.Error())
			return
		}
		s.logg.Info(fmt.Sprintf("it was selected on the day %s - %d rows: ", st, len(events)))
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}
