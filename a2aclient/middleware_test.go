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
)

func TestCallMetaFrom(t *testing.T) {
	ctx := context.Background()
	meta := CallMeta{"key": "value"}
	ctxWithMeta := context.WithValue(ctx, callMetaKey{}, meta)

	// Test with meta
	retrievedMeta, ok := CallMetaFrom(ctxWithMeta)
	if !ok {
		t.Fatal("expected to find meta")
	}
	if retrievedMeta["key"] != "value" {
		t.Errorf("unexpected meta value: got %q", retrievedMeta["key"])
	}

	// Test without meta
	_, ok = CallMetaFrom(ctx)
	if ok {
		t.Fatal("expected not to find meta")
	}
}

func TestCallContextFrom(t *testing.T) {
	ctx := context.Background()
	callCtx := CallContext{SessionID: "test-sid"}
	ctxWithCallCtx := context.WithValue(ctx, callContextKey{}, callCtx)

	// Test with call context
	retrievedCallCtx, ok := CallContextFrom(ctxWithCallCtx)
	if !ok {
		t.Fatal("expected to find call context")
	}
	if retrievedCallCtx.SessionID != "test-sid" {
		t.Errorf("unexpected session id: got %q", retrievedCallCtx.SessionID)
	}

	// Test without call context
	_, ok = CallContextFrom(ctx)
	if ok {
		t.Fatal("expected not to find call context")
	}
}

func TestWithSessionID(t *testing.T) {
	ctx := context.Background()
	sid := SessionID("test-sid")

	// Test adding a new session id
	ctxWithSid := WithSessionID(ctx, sid)
	callCtx, ok := CallContextFrom(ctxWithSid)
	if !ok {
		t.Fatal("expected to find call context")
	}
	if callCtx.SessionID != sid {
		t.Errorf("unexpected session id: got %q, want %q", callCtx.SessionID, sid)
	}

	// Test overwriting an existing session id
	newSid := SessionID("new-test-sid")
	ctxWithNewSid := WithSessionID(ctxWithSid, newSid)
	callCtx, ok = CallContextFrom(ctxWithNewSid)
	if !ok {
		t.Fatal("expected to find call context")
	}
	if callCtx.SessionID != newSid {
		t.Errorf("unexpected session id: got %q, want %q", callCtx.SessionID, newSid)
	}
}

func TestPassthroughInterceptor(t *testing.T) {
	interceptor := PassthroughInterceptor{}
	ctx := context.Background()
	req := &Request{}
	resp := &Response{}

	// Test Before
	newCtx, err := interceptor.Before(ctx, req)
	if err != nil {
		t.Errorf("unexpected error from Before: %v", err)
	}
	if newCtx != ctx {
		t.Error("context was modified by Before")
	}

	// Test After
	err = interceptor.After(ctx, resp)
	if err != nil {
		t.Errorf("unexpected error from After: %v", err)
	}
}
