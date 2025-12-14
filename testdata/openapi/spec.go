package openapi

import (
	"context"

	"github.com/goatx/goat"
	"github.com/goatx/goat/openapi"
)

type ClientStateMachine struct {
	goat.StateMachine
	Service *UserService
}

type UserService struct {
	goat.StateMachine
	Client *ClientStateMachine
}

type ClientState struct {
	goat.State
}

type ServiceState struct {
	goat.State
}

type CreateUserRequest struct {
	openapi.Schema[*ClientStateMachine, *UserService]
	Username string `json:"username" openapi:"required"`
	Email    string `json:"email" openapi:"required"`
}

type CreateUserResponse struct {
	openapi.Schema[*UserService, *ClientStateMachine]
	UserID int    `json:"user_id"`
	Status string `json:"status"`
}

func createUserServiceModel() {
	clientSpec := goat.NewStateMachineSpec(&ClientStateMachine{})
	serviceSpec := openapi.NewServiceSpec(&UserService{})

	clientState := &ClientState{}
	serviceState := &ServiceState{}

	clientSpec.DefineStates(clientState).SetInitialState(clientState)
	serviceSpec.DefineStates(serviceState).SetInitialState(serviceState)

	goat.OnEntry(clientSpec, clientState,
		func(ctx context.Context, client *ClientStateMachine) {
			request := &CreateUserRequest{
				Username: "test",
				Email:    "test@example.com",
			}
			openapi.SendTo(ctx, client.Service, request)
		})

	openapi.OnRequest(serviceSpec, serviceState, openapi.HTTPMethodPost, "/users",
		func(ctx context.Context, req *CreateUserRequest, service *UserService) openapi.Response[*CreateUserResponse] {
			response := &CreateUserResponse{
				UserID: 1,
				Status: "created",
			}
			return openapi.SendTo(ctx, service.Client, response)
		},
		openapi.WithOperationID("createUser"),
	)
}
