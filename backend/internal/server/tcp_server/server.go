package tcp_server

import (
	"context"
	"log/slog"
	"msg_app/internal/message"
	"net"
	"os"
)

type Server struct {
	hostIP   string
	listener net.Listener
	logger   *slog.Logger

	groups      []uint
	users       []uint
	group_count uint
	user_count  uint
}

func Init(logger *slog.Logger) *Server {
	return &Server{
		hostIP:      os.Getenv("TCP_SERVER_IP"),
		listener:    nil,
		groups:      make([]uint, 0),
		users:       make([]uint, 0),
		group_count: 0,
		user_count:  0,
		logger:      logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.hostIP)
	if err != nil {
		s.logger.Error("failed to listen tcp", "error", err)
		return err
	}

	s.listener = listener
	s.logger.Debug("server running", "ip", s.hostIP)

	defer s.listener.Close()

	go s.handle_acception()

	<-ctx.Done()

	s.listener.Close()

	return ctx.Err()
}

func (s *Server) GetUserCount() uint {
	return s.user_count
}

func (s *Server) GetGroupCount() uint {
	return s.group_count
}

func (s *Server) handle_acception() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error("error accepting connection", "error", err)
			return
		}

		go s.handle_connection(conn)
	}
}

func (s *Server) handle_connection(conn net.Conn) {
	defer conn.Close()

	for {
		msg := message.NewMsg(s.logger)
		msg_data, err := msg.HandleMsg(&conn)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}

		s.logger.Info("message received", "ip", conn.RemoteAddr(), "message", msg_data)
	}
}
