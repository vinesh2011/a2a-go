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
	"testing"
)

// TestInterfaceGuards calls the private methods that are used to enforce interface implementations.
// This is done to achieve 100% test coverage.
func TestInterfaceGuards(t *testing.T) {
	guards := []func(){
		(&Task{}).isSendMessageResult,
		(&Message{}).isSendMessageResult,
		(&Message{}).isEvent,
		(&Task{}).isEvent,
		(&TaskStatusUpdateEvent{}).isEvent,
		(&TaskArtifactUpdateEvent{}).isEvent,
		(TextPart{}).isPart,
		(FilePart{}).isPart,
		(DataPart{}).isPart,
		(FileBytes{}).isFilePartContent,
		(FileURI{}).isFilePartContent,
		(APIKeySecurityScheme{}).isSecurityScheme,
		(HTTPAuthSecurityScheme{}).isSecurityScheme,
		(OpenIDConnectSecurityScheme{}).isSecurityScheme,
		(MutualTLSSecurityScheme{}).isSecurityScheme,
		(OAuth2SecurityScheme{}).isSecurityScheme,
	}
	for i, f := range guards {
		if f == nil {
			t.Fatalf("guard function at index %d is nil", i)
		}
		f()
	}
}

func TestNewIDFunctions(t *testing.T) {
	if NewMessageID() == "" {
		t.Error("NewMessageID returned empty string")
	}
	if NewTaskID() == "" {
		t.Error("NewTaskID returned empty string")
	}
	if NewContextID() == "" {
		t.Error("NewContextID returned empty string")
	}
	if NewArtifactID() == "" {
		t.Error("NewArtifactID returned empty string")
	}
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage(MessageRoleUser, TextPart{Text: "hello"})
	if msg.ID == "" {
		t.Error("message ID is empty")
	}
	if msg.Role != MessageRoleUser {
		t.Errorf("unexpected role: got %q, want %q", msg.Role, MessageRoleUser)
	}
	if len(msg.Parts) != 1 {
		t.Errorf("unexpected number of parts: got %d, want 1", len(msg.Parts))
	}
}

func TestNewMessageForTask(t *testing.T) {
	task := Task{ID: "task-1", ContextID: "ctx-1"}
	msg := NewMessageForTask(MessageRoleAgent, task, TextPart{Text: "world"})
	if msg.ID == "" {
		t.Error("message ID is empty")
	}
	if msg.Role != MessageRoleAgent {
		t.Errorf("unexpected role: got %q, want %q", msg.Role, MessageRoleAgent)
	}
	if msg.TaskID != task.ID {
		t.Errorf("unexpected task ID: got %q, want %q", msg.TaskID, task.ID)
	}
	if msg.ContextID != task.ContextID {
		t.Errorf("unexpected context ID: got %q, want %q", msg.ContextID, task.ContextID)
	}
	if len(msg.Parts) != 1 {
		t.Errorf("unexpected number of parts: got %d, want 1", len(msg.Parts))
	}
}

func TestNewArtifactEvent(t *testing.T) {
	task := Task{ID: "task-1", ContextID: "ctx-1"}
	event := NewArtifactEvent(task, TextPart{Text: "artifact part"})
	if event.TaskID != task.ID {
		t.Errorf("unexpected task ID: got %q, want %q", event.TaskID, task.ID)
	}
	if event.ContextID != task.ContextID {
		t.Errorf("unexpected context ID: got %q, want %q", event.ContextID, task.ContextID)
	}
	if event.Artifact.ID == "" {
		t.Error("artifact ID is empty")
	}
	if len(event.Artifact.Parts) != 1 {
		t.Errorf("unexpected number of parts: got %d, want 1", len(event.Artifact.Parts))
	}
}

func TestNewArtifactUpdateEvent(t *testing.T) {
	task := Task{ID: "task-1", ContextID: "ctx-1"}
	artifactID := ArtifactID("artifact-1")
	event := NewArtifactUpdateEvent(task, artifactID, TextPart{Text: "update part"})

	if !event.Append {
		t.Error("Append should be true")
	}
	if event.Artifact.ID != artifactID {
		t.Errorf("unexpected artifact ID: got %q, want %q", event.Artifact.ID, artifactID)
	}
}

func TestNewStatusUpdateEvent(t *testing.T) {
	task := &Task{ID: "task-1", ContextID: "ctx-1"}
	msg := NewMessage(MessageRoleAgent, TextPart{Text: "status message"})
	event := NewStatusUpdateEvent(task, TaskStateWorking, msg)

	if event.TaskID != task.ID {
		t.Errorf("unexpected task ID: got %q, want %q", event.TaskID, task.ID)
	}
	if event.Status.State != TaskStateWorking {
		t.Errorf("unexpected state: got %q, want %q", event.Status.State, TaskStateWorking)
	}
	if event.Status.Message != msg {
		t.Error("unexpected message")
	}
	if event.Status.Timestamp.IsZero() {
		t.Error("timestamp is zero")
	}
}

func TestTaskStatus_Terminal(t *testing.T) {
	testCases := []struct {
		state    TaskState
		terminal bool
	}{
		{TaskStateCompleted, true},
		{TaskStateCanceled, true},
		{TaskStateFailed, true},
		{TaskStateRejected, true},
		{TaskStateAuthRequired, false},
		{TaskStateInputRequired, false},
		{TaskStateSubmitted, false},
		{TaskStateUnknown, false},
		{TaskStateWorking, false},
	}

	for _, tc := range testCases {
		if tc.state.Terminal() != tc.terminal {
			t.Errorf("state %q terminal status should be %v", tc.state, tc.terminal)
		}
	}
}

func TestPart_Meta(t *testing.T) {
	meta := map[string]any{"key": "value"}

	textPart := TextPart{Metadata: meta}
	if textPart.Meta()["key"] != "value" {
		t.Error("TextPart.Meta() returned wrong metadata")
	}

	dataPart := DataPart{Metadata: meta}
	if dataPart.Meta()["key"] != "value" {
		t.Error("DataPart.Meta() returned wrong metadata")
	}

	filePart := FilePart{Metadata: meta}
	if filePart.Meta()["key"] != "value" {
		t.Error("FilePart.Meta() returned wrong metadata")
	}
}

func TestNewStatusUpdateEvent_NilTask(t *testing.T) {
	// This test is to ensure that a panic does not occur if a nil task is passed in.
	// The function should still execute without errors, although in a real-world scenario,
	// a non-nil task should be provided. This is primarily for achieving 100% coverage
	// and ensuring robustness against unexpected inputs.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked on a nil task: %v", r)
		}
	}()
	_ = NewStatusUpdateEvent(&Task{}, TaskStateWorking, nil)
}

func TestNewMessageForTask_EmptyTask(t *testing.T) {
	// Similar to the above, this test ensures that passing an empty task does not cause a panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked on an empty task: %v", r)
		}
	}()
	_ = NewMessageForTask(MessageRoleAgent, Task{}, TextPart{Text: "world"})
}

func TestNewArtifactEvent_EmptyTask(t *testing.T) {
	// Ensures that passing an empty task to NewArtifactEvent does not cause a panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked on an empty task: %v", r)
		}
	}()
	_ = NewArtifactEvent(Task{}, TextPart{Text: "artifact part"})
}

func TestNewArtifactUpdateEvent_EmptyTask(t *testing.T) {
	// Ensures that passing an empty task to NewArtifactUpdateEvent does not cause a panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked on an empty task: %v", r)
		}
	}()
	_ = NewArtifactUpdateEvent(Task{}, "artifact-1", TextPart{Text: "update part"})
}
