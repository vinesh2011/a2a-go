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

package agentcard

import (
	"context"
	"testing"
)

func TestResolver_Resolve(t *testing.T) {
	resolver := &Resolver{BaseURL: "http://localhost"}
	ctx := context.Background()

	// Test with no options
	_, err := resolver.Resolve(ctx)
	if err == nil || err.Error() != "not implemented" {
		t.Errorf("expected 'not implemented' error, got %v", err)
	}

	// Test with WithPath option
	_, err = resolver.Resolve(ctx, WithPath("/new-path"))
	if err == nil || err.Error() != "not implemented" {
		t.Errorf("expected 'not implemented' error, got %v", err)
	}

	// Test with WithRequestHeaders option
	headers := map[string]string{"X-Test": "true"}
	_, err = resolver.Resolve(ctx, WithRequestHeaders(headers))
	if err == nil || err.Error() != "not implemented" {
		t.Errorf("expected 'not implemented' error, got %v", err)
	}
}

// Test options to ensure they don't panic and can be created
func TestResolveOptions(t *testing.T) {
	pathOpt := WithPath("/some/path")
	if pathOpt == nil {
		t.Error("WithPath returned nil")
	}

	headers := map[string]string{"Content-Type": "application/json"}
	headersOpt := WithRequestHeaders(headers)
	if headersOpt == nil {
		t.Error("WithRequestHeaders returned nil")
	}
}
