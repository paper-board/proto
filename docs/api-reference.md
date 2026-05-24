# API Reference

Generated Go code lives under `gen/go/<package>/v1/`. TypeScript under `gen/ts/<package>/v1/`. Import via `github.com/paper-board/proto/gen/go/<package>/v1`.

## Packages

| Package       | Status  | Go import suffix         | Proto file(s)                          |
| ------------- | ------- | ------------------------ | -------------------------------------- |
| `common.v1`   | stable  | `common/v1;commonv1`     | errors, ids, paging, timestamps, usage |
| `identity.v1` | stable  | `identity/v1;identityv1` | auth.proto                             |
| `runtime.v1`  | stable  | `runtime/v1;runtimev1`   | runtime.proto                          |
| `compute.v1`  | stable  | `compute/v1;computev1`   | compute.proto                          |
| `agents.v1`   | Phase 4 | —                        | —                                      |
| `platform.v1` | Phase 4 | —                        | —                                      |
| `billing.v1`  | Phase 5 | —                        | —                                      |

______________________________________________________________________

## common.v1

Shared types imported by every service package.

### ErrorDetail

Structured error returned in gRPC status details. Mapped from `paper-board/sdk/errors` sentinels.

| Field                 | Type                 | Description                                                                     |
| --------------------- | -------------------- | ------------------------------------------------------------------------------- |
| `code`                | `string`             | Sentinel code (`not_found`, `conflict`, `unauthorized`, `permission_denied`, …) |
| `message`             | `string`             | Human-readable; do not parse programmatically                                   |
| `fields`              | `map<string,string>` | Optional structured fields for client-side handling                             |
| `retry_after_seconds` | `optional int32`     | Hint for rate-limited / unavailable responses                                   |

### UsageMetrics

Resource consumption snapshot. Emitted by `compute.v1` on destroy / exec completion.

| Field                   | Type     | Unit  |
| ----------------------- | -------- | ----- |
| `cpu_ms`                | `uint64` | ms    |
| `memory_peak_bytes`     | `uint64` | bytes |
| `memory_avg_bytes`      | `uint64` | bytes |
| `disk_read_bytes`       | `uint64` | bytes |
| `disk_write_bytes`      | `uint64` | bytes |
| `network_egress_bytes`  | `uint64` | bytes |
| `network_ingress_bytes` | `uint64` | bytes |
| `gpu_ms`                | `uint64` | ms    |
| `gpu_memory_bytes`      | `uint64` | bytes |

Other shared messages: `PageRequest` / `PageResponse` (cursor-based paging), `ResourceId` (UUID wrapper), `Timestamps` (created_at / updated_at).

______________________________________________________________________

## identity.v1

Internal gRPC service. Consumed by `paper-board/agents`, `paper-board/runtime`, and `paper-board/gateway` (Phase 7) to validate caller credentials.

### AuthService

| RPC            | Request               | Response               | Description                               |
| -------------- | --------------------- | ---------------------- | ----------------------------------------- |
| `VerifyAPIKey` | `VerifyAPIKeyRequest` | `VerifyAPIKeyResponse` | Validates an API key; returns AuthContext |
| `VerifyJWT`    | `VerifyJWTRequest`    | `VerifyJWTResponse`    | Validates a JWT; returns AuthContext      |
| `IssueJWT`     | `IssueJWTRequest`     | `IssueJWTResponse`     | Issues a signed JWT for a user/org        |
| `GetPublicKey` | `GetPublicKeyRequest` | `GetPublicKeyResponse` | Returns the signing public key by kid     |

### AuthContext

Propagated downstream via gRPC metadata (`x-user-id`, `x-org-id`, `x-roles`, `x-trace-id` per ADR-0005).

| Field         | Type                             |
| ------------- | -------------------------------- |
| `user_id`     | `string`                         |
| `org_id`      | `string`                         |
| `mode`        | `AuthMode` enum                  |
| `auth_key_id` | `string`                         |
| `api_key_id`  | `string`                         |
| `method`      | `string`                         |
| `env`         | `Env` enum (`LIVE` / `TEST`)     |
| `role`        | `Role` enum (`OWNER` / `MEMBER`) |
| `expires_at`  | `google.protobuf.Timestamp`      |

______________________________________________________________________

## runtime.v1

Internal service. `paper-board/runtime` implements the server; `paper-board/agents` is the caller.

### RuntimeService

| RPC      | Request         | Response             | Description                                      |
| -------- | --------------- | -------------------- | ------------------------------------------------ |
| `Invoke` | `InvokeRequest` | `stream InvokeEvent` | Forwards agent prompt invocation; streams events |

Tenant context is propagated via gRPC metadata (`x-tenant-id`, `x-org-id`) — not in the request body.

### InvokeRequest

| Field        | Type                     |
| ------------ | ------------------------ |
| `session_id` | `string`                 |
| `request_id` | `string`                 |
| `trace_id`   | `string`                 |
| `messages`   | `repeated PromptMessage` |

### InvokeEvent (oneof)

| Variant      | Fields                                |
| ------------ | ------------------------------------- |
| `text_chunk` | `content string`                      |
| `turn_done`  | `tokens_in int32`, `tokens_out int32` |
| `error`      | `code string`, `message string`       |

All variants carry `emitted_at google.protobuf.Timestamp` at field 10.

______________________________________________________________________

## compute.v1

Internal service. `paper-board/compute` implements the server; `paper-board/runtime` is the primary caller.

### ComputeService

| RPC               | Request                  | Response                  | Description                            |
| ----------------- | ------------------------ | ------------------------- | -------------------------------------- |
| `CreateSandbox`   | `CreateSandboxRequest`   | `CreateSandboxResponse`   | Provisions a gVisor sandbox pod        |
| `DestroySandbox`  | `DestroySandboxRequest`  | `DestroySandboxResponse`  | Tears down sandbox; returns usage      |
| `DescribeSandbox` | `DescribeSandboxRequest` | `DescribeSandboxResponse` | Current status and live usage snapshot |
| `ExecCommand`     | `ExecCommandRequest`     | `stream ExecEvent`        | Runs a command; streams stdout/stderr  |
| `ReadFile`        | `ReadFileRequest`        | `ReadFileResponse`        | Reads up to 16 MiB from workspace      |
| `WriteFile`       | `WriteFileRequest`       | `WriteFileResponse`       | Writes up to 16 MiB to workspace       |
| `ListFiles`       | `ListFilesRequest`       | `ListFilesResponse`       | Lists workspace entries (default 1000) |

A bidirectional `ExecSession` RPC is deferred to Phase 5 (interactive REPL; per D24).

### SandboxStatus enum

`UNSPECIFIED` · `STARTING` · `READY` · `FAILED` · `DESTROYED`

### ExecEvent (oneof)

| Variant     | Fields                                                                          |
| ----------- | ------------------------------------------------------------------------------- |
| `stdout`    | `data bytes`                                                                    |
| `stderr`    | `data bytes`                                                                    |
| `completed` | `exit_code int32`, `wall_time Duration`, `cancelled bool`, `usage UsageMetrics` |

All variants carry `emitted_at google.protobuf.Timestamp` at field 10.

### Workspace path policy

All paths must be under `WorkspaceSpec.mount_path`. Symlink escapes are rejected via in-pod `realpath` validation. Read/Write capped at 16 MiB per call. `ListFiles` defaults to 1000 entries.

______________________________________________________________________

## Field validation patterns

These patterns are applied consistently across all packages:

| Pattern               | Example fields                                        | Rule                                                                       |
| --------------------- | ----------------------------------------------------- | -------------------------------------------------------------------------- |
| UUID string IDs       | `sandbox_id`, `session_id`, `tenant_id`               | Non-empty; validated at service boundary, not here                         |
| Enum zero value       | `SANDBOX_STATUS_UNSPECIFIED`, `AUTH_MODE_UNSPECIFIED` | Always defined; callers treat zero as unset                                |
| `optional` scalar     | `retry_after_seconds`                                 | Use proto3 `optional`; absence is meaningful                               |
| `oneof` event streams | `InvokeEvent`, `ExecEvent`                            | Exactly one variant set per message                                        |
| Reserved field 10     | `emitted_at`, `current_usage`, `total_usage`          | High-number field reserved for cross-cutting metadata to avoid renumbering |

## PACKAGE_DIRECTORY_MATCH discipline

`buf lint` (STANDARD ruleset) requires the proto `package` statement to match the directory path. Every `.proto` file under `<svc>/v1/` must declare `package <svc>.v1;`. Mismatches fail CI. The `PACKAGE_VERSION_SUFFIX` rule is excepted in `buf.yaml` to allow `v1` without a `v1` suffix in the package name literal.
