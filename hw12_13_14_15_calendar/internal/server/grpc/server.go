package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/pb"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedEventsApiServer
	app  app.Application
	logg logger.Logg
	conf config.Config
}

func NewServer(logger logger.Logg, app app.Application, config config.Config) *Server {
	return &Server{
		logg: logger,
		app:  app,
		conf: config,
	}
}

func (s *Server) Start(ctx context.Context, wg *sync.WaitGroup) error {
	lsn, err := net.Listen("tcp", s.conf.GRPCServer.Address)
	if err != nil {
		s.logg.Error(err)
		wg.Done()
		return nil
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryServerMiddleWareInterceptor(s.loggingReq)),
	)
	pb.RegisterEventsApiServer(server, s)
	s.logg.Info("starting grpc server on ", lsn.Addr().String())
	go func() {
		<-ctx.Done()
		s.logg.Info("server grpc is stopping...")
		server.Stop()
		lsn.Close()
		wg.Done()
	}()
	if err := server.Serve(lsn); err != nil {
		s.logg.Error(err)
		lsn.Close()
	}
	return nil
}

func (s *Server) Create(ctx context.Context, req *pb.Event) (*pb.Responce, error) {
	ctx, cancel := context.WithTimeout(ctx, 3000*time.Millisecond)
	defer cancel()
	userID, _ := strconv.Atoi(req.UserID)
	appEv := app.Event{
		Title:       req.Title,
		DateTime:    req.GetDateTime().AsTime(),
		Duration:    req.GetDuration().AsDuration(),
		TimeBefore:  req.GetTimeBefore().AsDuration(),
		Description: req.Description,
		UserID:      userID,
	}
	u := s.app.CreateEvent(ctx, appEv)
	if u != uuid.Nil {
		s.logg.Info("create new event with uuid: ", u.String())
		return &pb.Responce{Result: u.String()}, nil
	}
	return &pb.Responce{Result: "error with adding an event"}, nil
}

func (s *Server) Update(ctx context.Context, req *pb.Event) (*pb.Responce, error) {
	ctx, cancel := context.WithTimeout(ctx, 3000*time.Millisecond)
	defer cancel()
	appEv := app.Event{
		Duration:    req.GetDuration().AsDuration(),
		TimeBefore:  req.GetTimeBefore().AsDuration(),
		Description: req.Description,
	}
	uuid := uuid.FromStringOrNil(req.ID)
	err := s.app.UpgateEvent(ctx, uuid, appEv)
	if err == nil {
		s.logg.Info("update en event with uuid: ", req.ID)
		return &pb.Responce{Result: req.ID}, nil
	}
	return &pb.Responce{Result: "error with updating an event"}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.RequestUuid) (*pb.Responce, error) {
	ctx, cancel := context.WithTimeout(ctx, 3000*time.Millisecond)
	defer cancel()
	u := req.GetUuid()
	err := s.app.DeleteEvent(ctx, uuid.FromStringOrNil(u))
	if err == nil {
		s.logg.Info("delete an event with uuid: ", req.GetUuid())
		return &pb.Responce{Result: u}, nil
	}
	return &pb.Responce{Result: "error when delete an event"}, nil
}

func eventsFormAppToView(eventsApp []app.Event) pb.Events {
	events := make([]*pb.Event, len(eventsApp))
	for i, ev := range eventsApp {
		timestamp := timestamp.Timestamp{
			Seconds: ev.DateTime.Unix(),
			Nanos:   0,
		}
		duration := duration.Duration{
			Seconds: int64(ev.Duration.Seconds()),
			Nanos:   0,
		}
		events[i] = &pb.Event{
			ID:          ev.ID.String(),
			Title:       ev.Title,
			DateTime:    &timestamp,
			Duration:    &duration,
			Description: ev.Description,
			UserID:      strconv.Itoa(ev.UserID),
		}
	}
	return pb.Events{Event: events}
}

func (s *Server) ListOnDay(ctx context.Context, req *pb.RequestDate) (*pb.Events, error) {
	ctx, cancel := context.WithTimeout(ctx, 3000*time.Millisecond)
	defer cancel()
	d, err := time.Parse("2006-01-02", req.GetDate())
	if err != nil {
		s.logg.Error(err)
		return nil, status.Error(codes.InvalidArgument, "date is not specified")
	}
	evApp := s.app.GetListOnDay(ctx, d)
	if len(evApp) == 0 {
		return &pb.Events{}, status.Error(codes.NotFound, "nothing was selected")
	}
	events := eventsFormAppToView(evApp)
	s.logg.Info(fmt.Sprintf("it was selected on the day %s - %d rows: ", req.GetDate(), len(events.Event)))
	return &events, nil
}

func (s *Server) ListOnWeek(ctx context.Context, req *pb.RequestDate) (*pb.Events, error) {
	ctx, cancel := context.WithTimeout(ctx, 3000*time.Millisecond)
	defer cancel()
	d, err := time.Parse("2006-01-02", req.GetDate())
	if err != nil {
		s.logg.Error(err)
		return nil, status.Error(codes.InvalidArgument, "date is not specified")
	}
	evApp := s.app.GetListOnWeek(ctx, d)
	if len(evApp) == 0 {
		return &pb.Events{}, status.Error(codes.NotFound, "nothing was selected")
	}
	events := eventsFormAppToView(evApp)
	_, w := d.ISOWeek()
	s.logg.Info(fmt.Sprintf("it was selected on the week %d - %d rows: ", w, len(events.Event)))
	return &events, nil
}

func (s *Server) ListOnMonth(ctx context.Context, req *pb.RequestDate) (*pb.Events, error) {
	ctx, cancel := context.WithTimeout(ctx, 3000*time.Millisecond)
	defer cancel()
	d, err := time.Parse("2006-01-02", req.GetDate())
	if err != nil {
		s.logg.Error(err)
		return nil, status.Error(codes.InvalidArgument, "date is not specified")
	}
	evApp := s.app.GetListOnMonth(ctx, d)
	if len(evApp) == 0 {
		return &pb.Events{}, status.Error(codes.NotFound, "nothing was selected")
	}
	events := eventsFormAppToView(evApp)
	s.logg.Info(fmt.Sprintf("it was selected on the month %s - %d rows: ", d.Format("2006-01"), len(events.Event)))
	return &events, nil
}
