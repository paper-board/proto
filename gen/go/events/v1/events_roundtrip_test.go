package eventsv1_test

import (
	"testing"
	"time"

	eventsv1 "github.com/paper-board/proto/gen/go/events/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEnvelopeRoundTrip(t *testing.T) {
	t.Parallel()
	inner := &eventsv1.UserCreated{
		UserId:       "user-001",
		Email:        "test@example.com",
		DisplayName:  "Test User",
		AuthProvider: "password",
		EmittedAt:    timestamppb.New(time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)),
	}
	anyPayload, err := anypb.New(inner)
	if err != nil {
		t.Fatalf("anypb.New: %v", err)
	}

	orig := &eventsv1.Envelope{
		EventId:       "evt-001",
		EventType:     "identity.user.created",
		SourceService: "identity",
		OccurredAt:    timestamppb.New(time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)),
		TraceId:       "trace-abc",
		OrgId:         "",
		Payload:       anyPayload,
		SchemaVersion: 1,
	}

	data, err := protojson.Marshal(orig)
	if err != nil {
		t.Fatalf("protojson.Marshal(Envelope): %v", err)
	}

	got := &eventsv1.Envelope{}
	if err := protojson.Unmarshal(data, got); err != nil {
		t.Fatalf("protojson.Unmarshal(Envelope): %v", err)
	}

	if got.EventId != orig.EventId {
		t.Errorf("EventId: got %q, want %q", got.EventId, orig.EventId)
	}
	if got.EventType != orig.EventType {
		t.Errorf("EventType: got %q, want %q", got.EventType, orig.EventType)
	}
	if got.SourceService != orig.SourceService {
		t.Errorf("SourceService: got %q, want %q", got.SourceService, orig.SourceService)
	}
	if got.SchemaVersion != orig.SchemaVersion {
		t.Errorf("SchemaVersion: got %d, want %d", got.SchemaVersion, orig.SchemaVersion)
	}

	innerGot := &eventsv1.UserCreated{}
	if err := got.Payload.UnmarshalTo(innerGot); err != nil {
		t.Fatalf("UnmarshalTo(UserCreated): %v", err)
	}
	if innerGot.UserId != inner.UserId {
		t.Errorf("UserCreated.UserId: got %q, want %q", innerGot.UserId, inner.UserId)
	}
	if innerGot.Email != inner.Email {
		t.Errorf("UserCreated.Email: got %q, want %q", innerGot.Email, inner.Email)
	}
}

func TestUserCreatedRoundTrip(t *testing.T) {
	t.Parallel()
	orig := &eventsv1.UserCreated{
		UserId:       "user-002",
		Email:        "user2@example.com",
		DisplayName:  "User Two",
		AuthProvider: "oauth_google",
		EmittedAt:    timestamppb.Now(),
	}
	data, err := protojson.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := &eventsv1.UserCreated{}
	if err := protojson.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.UserId != orig.UserId || got.Email != orig.Email || got.AuthProvider != orig.AuthProvider {
		t.Errorf("field mismatch: got %+v, want %+v", got, orig)
	}
}

func TestOnboardingCompletedRoundTrip(t *testing.T) {
	t.Parallel()
	orig := &eventsv1.OnboardingCompleted{
		UserId:    "user-003",
		OrgId:     "org-001",
		VaultId:   "vault-001",
		EnvId:     "env-001",
		AgentId:   "agent-001",
		StartedAt: timestamppb.Now(),
		EmittedAt: timestamppb.Now(),
	}
	data, err := protojson.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := &eventsv1.OnboardingCompleted{}
	if err := protojson.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.UserId != orig.UserId || got.VaultId != orig.VaultId || got.AgentId != orig.AgentId {
		t.Errorf("field mismatch: got %+v, want %+v", got, orig)
	}
}

func TestOnboardingFailedRoundTrip(t *testing.T) {
	t.Parallel()
	orig := &eventsv1.OnboardingFailed{
		UserId:       "user-004",
		OrgId:        "org-002",
		LastStep:     "STEP_VAULT_PROVISION",
		ErrorCode:    "UNAVAILABLE",
		ErrorMessage: "vault service unreachable",
		Attempts:     3,
		EmittedAt:    timestamppb.Now(),
	}
	data, err := protojson.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := &eventsv1.OnboardingFailed{}
	if err := protojson.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.LastStep != orig.LastStep || got.ErrorCode != orig.ErrorCode || got.Attempts != orig.Attempts {
		t.Errorf("field mismatch: got %+v, want %+v", got, orig)
	}
}

func TestOnboardingDLQMovedRoundTrip(t *testing.T) {
	t.Parallel()
	orig := &eventsv1.OnboardingDLQMoved{
		UserId:            "user-005",
		OrgId:             "org-003",
		LastStep:          "STEP_ENV_PROVISION",
		FinalErrorMessage: "exceeded max attempts",
		Attempts:          5,
		EmittedAt:         timestamppb.Now(),
	}
	data, err := protojson.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := &eventsv1.OnboardingDLQMoved{}
	if err := protojson.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.LastStep != orig.LastStep || got.Attempts != orig.Attempts {
		t.Errorf("field mismatch: got %+v, want %+v", got, orig)
	}
}

func TestSessionStartedRoundTrip(t *testing.T) {
	t.Parallel()
	orig := &eventsv1.SessionStarted{
		SessionId: "sess-001",
		OrgId:     "org-001",
		UserId:    "user-001",
		AgentId:   "agent-001",
		EnvId:     "env-001",
		StartedAt: timestamppb.Now(),
		EmittedAt: timestamppb.Now(),
	}
	data, err := protojson.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := &eventsv1.SessionStarted{}
	if err := protojson.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.SessionId != orig.SessionId || got.EnvId != orig.EnvId {
		t.Errorf("field mismatch: got %+v, want %+v", got, orig)
	}
}

func TestSessionEndedRoundTrip(t *testing.T) {
	t.Parallel()
	orig := &eventsv1.SessionEnded{
		SessionId:       "sess-002",
		OrgId:           "org-001",
		UserId:          "user-001",
		AgentId:         "agent-001",
		StartedAt:       timestamppb.Now(),
		EndedAt:         timestamppb.Now(),
		DurationSeconds: 120,
		Reason:          eventsv1.TerminationReason_TERMINATION_REASON_USER_QUIT,
		EmittedAt:       timestamppb.Now(),
	}
	data, err := protojson.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := &eventsv1.SessionEnded{}
	if err := protojson.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.DurationSeconds != orig.DurationSeconds || got.Reason != orig.Reason {
		t.Errorf("field mismatch: got %+v, want %+v", got, orig)
	}
}

func TestLLMTokenEmittedRoundTrip(t *testing.T) {
	t.Parallel()
	orig := &eventsv1.LLMTokenEmitted{
		SessionId:        "sess-003",
		OrgId:            "org-001",
		UserId:           "user-001",
		AgentId:          "agent-001",
		Model:            "claude-opus-4-7",
		InputTokens:      1024,
		OutputTokens:     512,
		CacheReadTokens:  256,
		CacheWriteTokens: 128,
		EmittedAt:        timestamppb.Now(),
	}
	data, err := protojson.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := &eventsv1.LLMTokenEmitted{}
	if err := protojson.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Model != orig.Model || got.InputTokens != orig.InputTokens || got.CacheReadTokens != orig.CacheReadTokens {
		t.Errorf("field mismatch: got %+v, want %+v", got, orig)
	}
}

func TestSandboxLifecycleRoundTrip(t *testing.T) {
	t.Parallel()
	orig := &eventsv1.SandboxLifecycle{
		SandboxId: "sb-001",
		SessionId: "sess-001",
		OrgId:     "org-001",
		Phase:     eventsv1.SandboxLifecycle_PHASE_READY,
		Reason:    "",
		EmittedAt: timestamppb.Now(),
	}
	data, err := protojson.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := &eventsv1.SandboxLifecycle{}
	if err := protojson.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.SandboxId != orig.SandboxId || got.Phase != orig.Phase {
		t.Errorf("field mismatch: got %+v, want %+v", got, orig)
	}
}
