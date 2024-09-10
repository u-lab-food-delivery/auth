package server

import (
	"auth_service/config"
	"auth_service/genproto/auth"
	"auth_service/service"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func Run(s interface{}, cnf config.Config) error {
	server := grpc.NewServer()

	switch ser := s.(type) {
	case *service.AuthService:
		auth.RegisterAuthServiceServer(server, ser)
	default:
		return fmt.Errorf("unsupported service type: %T", ser)
	}

	lst, err := net.Listen("tcp", cnf.AuthServer.Host+":"+cnf.AuthServer.Port)
	if err != nil {
		return err
	}

	return server.Serve(lst)
}
