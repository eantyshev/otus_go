package main

import (
	"github.com/eantyshev/otus_go/calendar/api/server"
	pb "github.com/eantyshev/otus_go/calendar/internal/adapters"
	"github.com/eantyshev/otus_go/calendar/internal/adapters/db"
	"github.com/eantyshev/otus_go/calendar/internal/config"
	"github.com/eantyshev/otus_go/calendar/internal/interfaces"
	"github.com/eantyshev/otus_go/calendar/internal/logger"
	"github.com/eantyshev/otus_go/calendar/internal/usecases"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
)

func newCalendarServer(repo interfaces.Repository) *server.CalendarService {
	return &server.CalendarService{
		Usecases: &usecases.Usecases{Repo: repo, L: logger.L},
		L:        logger.L,
	}
}

func Server(addrPort string, pgDsn string) {
	repo, err := db.NewPgRepo(pgDsn)
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", addrPort)
	if err != nil {
		panic(err)
	}
	logger.L.Debugf("listening at %s", addrPort)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_zap.UnaryServerInterceptor(logger.L.Desugar()),
	))
	pb.RegisterCalendarServer(grpcServer, newCalendarServer(repo))
	panic(grpcServer.Serve(lis))
}

func main() {
	config.Configure()
	grpcHostPort := viper.GetString("grpc_listen")
	pgDsn := viper.GetString("pg_dsn")
	Server(grpcHostPort, pgDsn)
}
