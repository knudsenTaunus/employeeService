package proto

import (
	"context"
	userpb "github.com/knudsenTaunus/employeeService/generated/go/proto/user/v1"
	"github.com/knudsenTaunus/employeeService/internal/model"
	"github.com/rs/zerolog"
	"time"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	logger      zerolog.Logger
	updateChan  chan model.User
	subscribers map[string]userpb.UserService_RegisterServer
}

func NewServer(userChan chan model.User, logger zerolog.Logger) *UserServer {
	return &UserServer{
		UnimplementedUserServiceServer: userpb.UnimplementedUserServiceServer{},
		logger:                         logger,
		updateChan:                     userChan,
		subscribers:                    make(map[string]userpb.UserService_RegisterServer),
	}

}

func (u *UserServer) Invoke() {
	for {
		select {
		case update := <-u.updateChan:
			for k, v := range u.subscribers {
				err := v.Send(&userpb.UserServiceRegisterResponse{
					FirstName: update.FirstName,
					LastName:  update.LastName,
					Nickname:  "knudsenTaunus",
					Email:     "foo@bar.de",
					Country:   "germany",
					UpdatedAt: time.Now().String(),
				})
				if err != nil {
					u.logger.Error().Err(err).Msg("failed to send update")
					continue
				}

				u.logger.Info().Msgf("sent to %s", k)
			}
		}
	}
}

func (u *UserServer) Register(req *userpb.UserServiceRegisterRequest, stream userpb.UserService_RegisterServer) error {
	u.logger.Info().Msgf("got request from %s", req.ClientName)
	u.subscribers[req.ClientName] = stream
	for {
		select {
		case <-stream.Context().Done():
			delete(u.subscribers, req.ClientName)
			u.logger.Info().Msg("stream closed, client removed")
			return nil
		}
	}
}
func (u *UserServer) Deregister(ctx context.Context, req *userpb.UserServiceDeregisterRequest) (*userpb.UserServiceDeregisterResponse, error) {
	u.logger.Info().Msgf("remove client for updates %s", req.ClientName)
	delete(u.subscribers, req.ClientName)
	return &userpb.UserServiceDeregisterResponse{ClientName: req.GetClientName()}, nil
}
