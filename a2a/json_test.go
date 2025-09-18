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

package a2a

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func mustMarshal(t *testing.T, data any) string {
	t.Helper()
	bytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Marshal() failed with: %v", err)
	}
	return string(bytes)
}

func mustUnmarshal(t *testing.T, data []byte, out any) {
	t.Helper()
	if err := json.Unmarshal(data, out); err != nil {
		t.Fatalf("Unmarshal() failed with: %v", err)
	}
}

func TestFilePartJSONCodec(t *testing.T) {
	testCases := []struct {
		json string
		part FilePart
	}{
		{
			part: FilePart{File: FileURI{URI: "uri"}},
			json: `{"kind":"file","file":{"uri":"uri"}}`,
		},
		{
			part: FilePart{File: FileBytes{Bytes: "abc"}},
			json: `{"kind":"file","file":{"bytes":"abc"}}`,
		},
		{
			part: FilePart{File: FileBytes{Bytes: "abc", FileMeta: FileMeta{Name: "foo"}}},
			json: `{"kind":"file","file":{"name":"foo","bytes":"abc"}}`,
		},
		{
			part: FilePart{File: FileBytes{Bytes: "abc", FileMeta: FileMeta{Name: "foo", MimeType: "mime"}}},
			json: `{"kind":"file","file":{"mimeType":"mime","name":"foo","bytes":"abc"}}`,
		},
		{
			part: FilePart{File: FileURI{URI: "uri", FileMeta: FileMeta{Name: "foo", MimeType: "mime"}}},
			json: `{"kind":"file","file":{"mimeType":"mime","name":"foo","uri":"uri"}}`,
		},
	}
	for _, tc := range testCases {
		if got := mustMarshal(t, tc.part); got != tc.json {
			t.Fatalf("Marshal() failed:\nwant %v\ngot: %s", tc.json, got)
		}

		got := FilePart{}
		mustUnmarshal(t, []byte(tc.json), &got)
		if !reflect.DeepEqual(got, tc.part) {
			t.Fatalf("Unmarshal() failed for %s:\nwant %v\ngot: %s", tc.json, tc.part, got)
		}
	}
}

func TestFilePartJSONDecodingFailure(t *testing.T) {
	malformed := []string{
		`{"kind":"file"}`,
		`{"kind":"file","file":{}}`,
		`{"kind":"file","file":{"name":"foo","mimeType":"mime","uri":"uri","bytes":"abc"}}`,
		`{"kind":"file","file":{"name":"foo","mimeType":"mime"}}`,
		`not-a-json`,
	}
	for _, v := range malformed {
		got := FilePart{}
		if err := json.Unmarshal([]byte(v), &got); err == nil {
			t.Fatalf("Unmarshal() expected to fail for %s, got: %v", v, got)
		}
	}
}

func TestUnmarshalFailure(t *testing.T) {
	testCases := []struct {
		name      string
		target    any
		malformed []string
	}{
		{
			name:   "NamedSecuritySchemes",
			target: &NamedSecuritySchemes{},
			malformed: []string{
				`{"name1": "not-an-object"}`,
				`{"name1": {"type": "invalid-type"}}`,
				`{"name1": {"type": "apiKey", "name": 123}}`,
				`{"name1": {"type": "http", "scheme": 123}}`,
				`{"name1": {"type": "mutualTLS", "description": 123}}`,
				`{"name1": {"type": "oauth2", "flows": "not-an-object"}}`,
				`{"name1": {"type": "openIdConnect", "openIdConnectUrl": 123}}`,
			},
		},
		{
			name:   "ContentParts",
			target: &ContentParts{},
			malformed: []string{
				`["not-an-object"]`,
				`[{"kind": "invalid-kind"}]`,
				`[{"kind": "text", "text": {}}]`,
				`[{"kind": "data", "data": "not-a-map"}]`,
				`[{"kind": "file", "file": "not-an-object"}]`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, v := range tc.malformed {
				err := json.Unmarshal([]byte(v), tc.target)
				if err == nil {
					t.Fatalf("Unmarshal() expected to fail for %s", v)
				}
			}
		})
	}
}

func TestNamedSecuritySchemes_UnmarshalJSON_Error(t *testing.T) {
	var s NamedSecuritySchemes
	err := json.Unmarshal([]byte("not-a-json"), &s)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestContentParts_UnmarshalJSON_Error(t *testing.T) {
	var p ContentParts
	err := json.Unmarshal([]byte("not-a-json"), &p)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestContentPartsJSONCodec(t *testing.T) {
	parts := ContentParts{
		TextPart{Text: "hello, world"},
		FilePart{File: FileBytes{Bytes: "abc", FileMeta: FileMeta{Name: "foo", MimeType: "mime"}}},
		DataPart{Data: map[string]any{"foo": "bar"}},
		TextPart{Text: "42", Metadata: map[string]any{"foo": "bar"}},
	}

	jsons := []string{
		`{"kind":"text","text":"hello, world"}`,
		`{"kind":"file","file":{"mimeType":"mime","name":"foo","bytes":"abc"}}`,
		`{"kind":"data","data":{"foo":"bar"}}`,
		`{"kind":"text","text":"42","metadata":{"foo":"bar"}}`,
	}

	wantJSON := fmt.Sprintf("[%s]", strings.Join(jsons, ","))
	if got := mustMarshal(t, parts); got != wantJSON {
		t.Fatalf("Marshal() failed:\nwant %v\ngot: %s", wantJSON, got)
	}

	var got ContentParts
	mustUnmarshal(t, []byte(wantJSON), &got)
	if !reflect.DeepEqual(got, parts) {
		t.Fatalf("Unmarshal() failed:\nwant %v\ngot: %s", parts, got)
	}
}

func TestSecuritySchemeJSONCodec(t *testing.T) {
	schemes := NamedSecuritySchemes{
		"name1": APIKeySecurityScheme{Name: "abc", In: APIKeySecuritySchemeInCookie},
		"name2": OpenIDConnectSecurityScheme{OpenIDConnectURL: "url"},
		"name3": MutualTLSSecurityScheme{Description: "optional"},
		"name4": HTTPAuthSecurityScheme{Scheme: "Bearer", BearerFormat: "JWT"},
		"name5": OAuth2SecurityScheme{
			Flows: OAuthFlows{
				Password: &PasswordOAuthFlow{
					TokenURL: "url",
					Scopes:   map[string]string{"email": "read user emails"},
				}},
		},
	}

	entriesJSON := []string{
		`"name1":{"type":"apiKey","in":"cookie","name":"abc"}`,
		`"name2":{"type":"openIdConnect","openIdConnectUrl":"url"}`,
		`"name3":{"type":"mutualTLS","description":"optional"}`,
		`"name4":{"type":"http","bearerFormat":"JWT","scheme":"Bearer"}`,
		`"name5":{"type":"oauth2","flows":{"password":{"scopes":{"email":"read user emails"},"tokenUrl":"url"}}}`,
	}
	wantJSON := fmt.Sprintf("{%s}", strings.Join(entriesJSON, ","))

	var decodedJSON NamedSecuritySchemes
	mustUnmarshal(t, []byte(wantJSON), &decodedJSON)
	if !reflect.DeepEqual(decodedJSON, schemes) {
		t.Fatalf("Unmarshal() failed:\nwant %v\ngot: %s", schemes, decodedJSON)
	}

	encodedSchemes := mustMarshal(t, schemes)
	var decodedBack NamedSecuritySchemes
	mustUnmarshal(t, []byte(encodedSchemes), &decodedBack)
	if !reflect.DeepEqual(decodedJSON, decodedJSON) {
		t.Fatalf("Decoding back failed:\nwant %v\ngot: %s", decodedJSON, decodedBack)
	}
}
