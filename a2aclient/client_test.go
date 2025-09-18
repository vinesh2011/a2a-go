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
	"iter"
	"testing"

	"github.com/a2aproject/a2a-go/a2a"
)

// mockTransport is a mock implementation of the Transport interface for testing.
type mockTransport struct {
	destroyCalled bool
}

func (m *mockTransport) GetTask(ctx context.Context, query a2a.TaskQueryParams) (*a2a.Task, error) {
	return nil, nil
}
func (m *mockTransport) CancelTask(ctx context.Context, id a2a.TaskIDParams) (*a2a.Task, error) {
	return nil, nil
}
func (m *mockTransport) SendMessage(ctx context.Context, message a2a.MessageSendParams) (a2a.SendMessageResult, error) {
	return nil, nil
}
func (m *mockTransport) ResubscribeToTask(ctx context.Context, id a2a.TaskIDParams) iter.Seq2[a2a.Event, error] {
	return nil
}
func (m *mockTransport) SendStreamingMessage(ctx context.Context, message a2a.MessageSendParams) iter.Seq2[a2a.Event, error] {
	return nil
}
func (m *mockTransport) GetTaskPushConfig(ctx context.Context, params a2a.GetTaskPushConfigParams) (a2a.TaskPushConfig, error) {
	return a2a.TaskPushConfig{}, nil
}
func (m *mockTransport) ListTaskPushConfig(ctx context.Context, params a2a.ListTaskPushConfigParams) ([]a2a.TaskPushConfig, error) {
	return nil, nil
}
func (m *mockTransport) SetTaskPushConfig(ctx context.Context, params a2a.TaskPushConfig) (a2a.TaskPushConfig, error) {
	return a2a.TaskPushConfig{}, nil
}
func (m *mockTransport) DeleteTaskPushConfig(ctx context.Context, params a2a.DeleteTaskPushConfigParams) error {
	return nil
}
func (m *mockTransport) GetAgentCard(ctx context.Context) (*a2a.AgentCard, error) {
	return nil, nil
}
func (m *mockTransport) Destroy() error {
	m.destroyCalled = true
	return nil
}

func TestClient_AddCallInterceptor(t *testing.T) {
	client := &Client{}
	interceptor := &AuthInterceptor{}
	client.AddCallInterceptor(interceptor)

	if len(client.interceptors) != 1 {
		t.Errorf("expected 1 interceptor, got %d", len(client.interceptors))
	}
}

func TestClient_Destroy(t *testing.T) {
	transport := &mockTransport{}
	client := &Client{transport: transport}
	err := client.Destroy()

	if err != nil {
		t.Errorf("expected no error from Destroy, got %v", err)
	}
	if !transport.destroyCalled {
		t.Error("expected transport.Destroy to be called")
	}
}

func TestClient_NotImplemented(t *testing.T) {
	client := &Client{}
	ctx := context.Background()

	_, err := client.GetTask(ctx, a2a.TaskQueryParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	_, err = client.CancelTask(ctx, a2a.TaskIDParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	_, err = client.SendMessage(ctx, a2a.MessageSendParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	resubscribeSeq := client.ResubscribeToTask(ctx, a2a.TaskIDParams{})
	resubscribeSeq(func(e a2a.Event, err error) bool {
		if err != ErrNotImplemented {
			t.Errorf("expected ErrNotImplemented, got %v", err)
		}
		return false
	})

	sendStreamingSeq := client.SendStreamingMessage(ctx, a2a.MessageSendParams{})
	sendStreamingSeq(func(e a2a.Event, err error) bool {
		if err != ErrNotImplemented {
			t.Errorf("expected ErrNotImplemented, got %v", err)
		}
		return false
	})

	_, err = client.GetTaskPushConfig(ctx, a2a.GetTaskPushConfigParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	_, err = client.ListTaskPushConfig(ctx, a2a.ListTaskPushConfigParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	_, err = client.SetTaskPushConfig(ctx, a2a.TaskPushConfig{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	err = client.DeleteTaskPushConfig(ctx, a2a.DeleteTaskPushConfigParams{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	_, err = client.GetAgentCard(ctx)
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}
}
