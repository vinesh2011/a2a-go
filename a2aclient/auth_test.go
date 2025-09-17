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

func TestInMemoryCredentialsStore(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryCredentialsStore()
	sid := SessionID("test-session")
	scheme := a2a.SecuritySchemeName("test-scheme")
	cred := AuthCredential("test-credential")

	// 1. Test getting a credential that doesn't exist
	_, err := store.Get(ctx, sid, scheme)
	if err != ErrCredentialNotFound {
		t.Errorf("expected ErrCredentialNotFound, got %v", err)
	}

	// 2. Test setting and getting a credential
	store.Set(sid, scheme, cred)
	retrievedCred, err := store.Get(ctx, sid, scheme)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedCred != cred {
		t.Errorf("expected credential %q, got %q", cred, retrievedCred)
	}

	// 3. Test overwriting a credential
	newCred := AuthCredential("new-credential")
	store.Set(sid, scheme, newCred)
	retrievedCred, err = store.Get(ctx, sid, scheme)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedCred != newCred {
		t.Errorf("expected credential %q, got %q", newCred, retrievedCred)
	}

	// 4. Test getting a credential for a different scheme
	otherScheme := a2a.SecuritySchemeName("other-scheme")
	_, err = store.Get(ctx, sid, otherScheme)
	if err != ErrCredentialNotFound {
		t.Errorf("expected ErrCredentialNotFound, got %v", err)
	}

	// 5. Test getting a credential for a different session
	otherSid := SessionID("other-session")
	_, err = store.Get(ctx, otherSid, scheme)
	if err != ErrCredentialNotFound {
		t.Errorf("expected ErrCredentialNotFound, got %v", err)
	}
}
