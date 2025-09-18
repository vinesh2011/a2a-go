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
	"testing"

	"github.com/a2aproject/a2a-go/a2a"
)

type mockTransportFactory struct{}

func (m *mockTransportFactory) Create(ctx context.Context, url string, card *a2a.AgentCard) (Transport, error) {
	return &mockTransport{}, nil
}

func TestFactory_NewFactory(t *testing.T) {
	// Test with default options
	factory := NewFactory()
	if len(factory.transports) != 1 {
		t.Errorf("expected 1 default transport, got %d", len(factory.transports))
	}

	// Test with WithDefaultsDisabled
	factory = NewFactory(WithDefaultsDisabled())
	if len(factory.transports) != 0 {
		t.Errorf("expected 0 transports with defaults disabled, got %d", len(factory.transports))
	}

	// Test with custom options
	config := Config{AcceptedOutputModes: []string{"test"}}
	interceptor := &AuthInterceptor{}
	transportProtocol := a2a.TransportProtocol("test")
	transportFactory := &mockTransportFactory{}

	factory = NewFactory(
		WithConfig(config),
		WithInterceptors(interceptor),
		WithTransport(transportProtocol, transportFactory),
	)

	if factory.config.AcceptedOutputModes[0] != "test" {
		t.Error("config not applied")
	}
	if len(factory.interceptors) != 1 {
		t.Error("interceptors not applied")
	}
	if _, ok := factory.transports[transportProtocol]; !ok {
		t.Error("transport not applied")
	}
}

func TestFactory_WithAdditionalOptions(t *testing.T) {
	config := Config{AcceptedOutputModes: []string{"test"}}
	interceptor := &AuthInterceptor{}
	transportProtocol := a2a.TransportProtocol("test")
	transportFactory := &mockTransportFactory{}

	baseFactory := NewFactory(
		WithDefaultsDisabled(),
		WithConfig(config),
		WithInterceptors(interceptor),
		WithTransport(transportProtocol, transportFactory),
	)

	newConfig := Config{AcceptedOutputModes: []string{"new-test"}}
	newInterceptor := &AuthInterceptor{}
	newTransportProtocol := a2a.TransportProtocol("new-test")
	newTransportFactory := &mockTransportFactory{}

	extendedFactory := WithAdditionalOptions(*baseFactory,
		WithConfig(newConfig),
		WithInterceptors(newInterceptor),
		WithTransport(newTransportProtocol, newTransportFactory),
	)

	if extendedFactory.config.AcceptedOutputModes[0] != "new-test" {
		t.Error("new config not applied")
	}
	if len(extendedFactory.interceptors) != 2 {
		t.Errorf("expected 2 interceptors, got %d", len(extendedFactory.interceptors))
	}
	if _, ok := extendedFactory.transports[transportProtocol]; !ok {
		t.Error("base transport not preserved")
	}
	if _, ok := extendedFactory.transports[newTransportProtocol]; !ok {
		t.Error("new transport not applied")
	}
}

func TestFactory_CreateNotImplemented(t *testing.T) {
	// Test defaultsDisabledOpt.apply
	opt := WithDefaultsDisabled()
	opt.apply(&Factory{})

	factory := NewFactory()
	ctx := context.Background()

	_, err := factory.CreateFromCard(ctx, &a2a.AgentCard{})
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	// With options
	_, err = factory.CreateFromCard(ctx, &a2a.AgentCard{}, WithConfig(Config{}))
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	_, err = factory.CreateFromURL(ctx, "", nil)
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}

	// With options
	_, err = factory.CreateFromURL(ctx, "", nil, WithConfig(Config{}))
	if err != ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}
}
