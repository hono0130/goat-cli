package protobuf

import (
	"context"

	"github.com/goatx/goat"
	"github.com/goatx/goat/protobuf"
)

type ClientStateMachine struct {
	goat.StateMachine
	Service *ServiceStateMachine
}

type ServiceStateMachine struct {
	goat.StateMachine
	Client *ClientStateMachine
}

type ClientState struct {
	goat.State
}

type ServiceState struct {
	goat.State
}

type Request struct {
	protobuf.ProtobufMessage[*ClientStateMachine, *ServiceStateMachine]
}

type Response struct {
	protobuf.ProtobufMessage[*ServiceStateMachine, *ClientStateMachine]
}

func createProtobufModel() {
	clientSpec := goat.NewStateMachineSpec(&ClientStateMachine{})
	serviceSpec := protobuf.NewProtobufServiceSpec(&ServiceStateMachine{})

	clientState := &ClientState{}
	serviceState := &ServiceState{}

	clientSpec.DefineStates(clientState).SetInitialState(clientState)
	serviceSpec.DefineStates(serviceState).SetInitialState(serviceState)

	goat.OnEntry(clientSpec, clientState,
		func(ctx context.Context, client *ClientStateMachine) {
			request := &Request{}
			protobuf.ProtobufSendTo(ctx, client.Service, request)
		})

	protobuf.OnProtobufMessage(serviceSpec, serviceState, "HandleRequest",
		func(ctx context.Context, request *Request, service *ServiceStateMachine) protobuf.ProtobufResponse[*Response] {
			response := &Response{}
			return protobuf.ProtobufSendTo(ctx, service.Client, response)
		})
}
