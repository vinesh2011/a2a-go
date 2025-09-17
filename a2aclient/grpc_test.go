// Copyright 2025 The A2A Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package a2aclient

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/a2aproject/a2a-go/a2a"
)

func newTestGRPCServer(t *testing.T) (*grpc.Server, *bufconn.Listener) {
	t.Helper()
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	return s, lis
}

func TestWithGRPCTransport(t *testing.T) {
	s, lis := newTestGRPCServer(t)
	defer s.Stop()

	opt := WithGRPCTransport(
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
	)

	factory := NewFactory(WithDefaultsDisabled(), opt)
	transportFactory := factory.transports[a2a.TransportProtocolGRPC]
	if transportFactory == nil {
		t.Fatal("gRPC transport factory not registered")
	}

	transport, err := transportFactory.Create(context.Background(), "bufnet", nil)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}
	if transport == nil {
		t.Fatal("created transport is nil")
	}
	defer transport.Destroy()

	// Test with a failing dial option
	opt = WithGRPCTransport(grpc.WithAuthority("invalid-authority"))
	factory = NewFactory(WithDefaultsDisabled(), opt)
	transportFactory = factory.transports[a2a.TransportProtocolGRPC]
	if transportFactory == nil {
		t.Fatal("gRPC transport factory not registered")
	}
	_, err = transportFactory.Create(context.Background(), "bufnet", nil)
	if err == nil {
		t.Error("expected an error from Create, got nil")
	}
}

func TestGRPCTransport_Destroy(t *testing.T) {
	closeCalled := false
	transport := &grpcTransport{
		closeConnFn: func() error {
			closeCalled = true
			return nil
		},
	}
	err := transport.Destroy()
	if err != nil {
		t.Errorf("expected no error from Destroy, got %v", err)
	}
	if !closeCalled {
		t.Error("expected close function to be called")
	}
}

func TestGRPCTransport_NotImplemented(t *testing.T) {
	transport := &grpcTransport{}
	ctx := context.Background()

	_, err := transport.GetTask(ctx, a2a.TaskQueryParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	_, err = transport.CancelTask(ctx, a2a.TaskIDParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	_, err = transport.SendMessage(ctx, a2a.MessageSendParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	resubscribeSeq := transport.ResubscribeToTask(ctx, a2a.TaskIDParams{})
	resubscribeSeq(func(e a2a.Event, err error) bool {
		if err != ErrNotImplemented {
			t.Errorf("expected ErrNotImplemented, got %v", err)
		}
		return false
	})

	sendStreamingSeq := transport.SendStreamingMessage(ctx, a2a.MessageSendParams{})
	sendStreamingSeq(func(e a2a.Event, err error) bool {
		if err != ErrNotImplemented {
			t.Errorf("expected ErrNotImplemented, got %v", err)
		}
		return false
	})

	_, err = transport.GetTaskPushConfig(ctx, a2a.GetTaskPushConfigParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	_, err = transport.ListTaskPushConfig(ctx, a2a.ListTaskPushConfigParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	_, err = transport.SetTaskPushConfig(ctx, a2a.TaskPushConfig{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	err = transport.DeleteTaskPushConfig(ctx, a2a.DeleteTaskPushConfigParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	_, err = transport.GetAgentCard(ctx)
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}
}
