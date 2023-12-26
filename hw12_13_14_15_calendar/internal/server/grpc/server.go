package internalgrpc

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/pb"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
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

func (s *Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", s.conf.GRPCServer.Address)
	if err != nil {
		s.logg.Error(err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryServerMiddleWareInterceptor(s.loggingReq)),
	)
	pb.RegisterEventsApiServer(server, s)

	s.logg.Info("starting grpc server on ", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		s.logg.Error(err)
	}
	return nil
}

func (s *Server) Create(ctx context.Context, req *pb.Event) (*pb.Responce, error) {
	ctx, cancel := context.WithTimeout(ctx, 3000*time.Millisecond)
	defer cancel()
	userID, _ := strconv.Atoi(req.UserID)
	appEv := app.Event{
		Title:       req.Title,
		DateTime:    req.DateTime.AsTime(),
		Duration:    req.Duration.AsDuration(),
		TimeBefore:  req.TimeBefore.AsDuration(),
		Description: req.Description,
		UserID:      userID,
	}
	u := s.app.CreateEvent(ctx, appEv)
	if u != uuid.Nil {
		s.logg.Info("create new event with uuid: ", u.String())
		return &pb.Responce{Resalt: u.String()}, nil
	}
	return &pb.Responce{Resalt: "error with adding an event"}, nil
}

func (s *Server) Update(ctx context.Context, req *pb.Event) (*pb.Responce, error) {
	ctx, cancel := context.WithTimeout(ctx, 3000*time.Millisecond)
	defer cancel()
	appEv := app.Event{
		Duration:    req.Duration.AsDuration(),
		TimeBefore:  req.TimeBefore.AsDuration(),
		Description: req.Description,
	}
	uuid := uuid.FromStringOrNil(req.ID)
	err := s.app.UpgateEvent(ctx, uuid, appEv)
	if err == nil {
		s.logg.Info("update en event with uuid: ", req.ID)
		return &pb.Responce{Resalt: req.ID}, nil
	}
	return &pb.Responce{Resalt: "error with updating an event"}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.Request) (*pb.Responce, error) {
	ctx, cancel := context.WithTimeout(ctx, 3000*time.Millisecond)
	defer cancel()
	u := req.GetUuid()
	err := s.app.DeleteEvent(ctx, uuid.FromStringOrNil(u))
	if err == nil {
		s.logg.Info("delete an event with uuid: ", req.GetUuid())
		return &pb.Responce{Resalt: u}, nil
	}
	return &pb.Responce{Resalt: "error when delete an event"}, nil
}
