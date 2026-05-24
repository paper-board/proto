# paper-board/proto

`paper-board/proto` is the single source of truth for every gRPC contract in the paper-board platform. All backend services and the SDK consume generated code from this repository — no service owns a proto file outside this repo.

## Package map

```
paper-board/proto
├── common/v1      — shared messages: ErrorDetail, UsageMetrics, paging, IDs, timestamps
├── identity/v1    — AuthService (VerifyAPIKey, VerifyJWT, IssueJWT, GetPublicKey)
├── agents/v1      — (Phase 4)
├── runtime/v1     — RuntimeService (Invoke — streaming)
├── compute/v1     — ComputeService (sandbox lifecycle, exec, workspace I/O)
├── platform/v1    — (Phase 4)
├── billing/v1     — (Phase 5)
└── gen/           — CI-committed codegen output (Go + TypeScript)
```

Each package maps to one backend service. Namespace convention: `<service>.v1` (e.g. `identity.v1`, `compute.v1`). Breaking change → new package version (`identity.v2`), parallel support until consumers migrate.

## Consumers

| Consumer                  | Import path                                       |
| ------------------------- | ------------------------------------------------- |
| paper-board/identity      | `github.com/paper-board/proto/gen/go/identity/v1` |
| paper-board/agents        | `github.com/paper-board/proto/gen/go/agents/v1`   |
| paper-board/runtime       | `github.com/paper-board/proto/gen/go/runtime/v1`  |
| paper-board/compute       | `github.com/paper-board/proto/gen/go/compute/v1`  |
| paper-board/platform      | `github.com/paper-board/proto/gen/go/platform/v1` |
| paper-board/billing       | `github.com/paper-board/proto/gen/go/billing/v1`  |
| paper-board/sdk           | `github.com/paper-board/proto/gen/go/common/v1`   |
| paper-board/frontend (TS) | `gen/ts/<package>/v1`                             |

Services reference the proto module via `go get github.com/paper-board/proto@<tag>`. Local development uses `replace ../proto` in `go.mod` — that directive must never be committed to CI (see [Operations](operations.md#local-replace-directive)).

## Quickstart

Install [`buf`](https://buf.build/docs/installation), then:

```bash
# Validate proto files against buf lint rules
buf lint

# Check for breaking changes against main
buf breaking --against '.git#branch=main'

# Re-generate Go + TypeScript output
buf generate

# Verify gen/ is up to date (same check CI runs)
git diff --exit-code gen/
```

CI runs all three checks — `buf lint`, `buf breaking` (PRs only), and `buf generate` + drift check — on every push and pull request. A PR that introduces a lint violation, a breaking change, or uncommitted generated output will not merge.

## Adding a new service package

```bash
# 1. Create the package directory
mkdir -p <svc>/v1

# 2. Write <svc>/v1/<svc>.proto
#    - package <svc>.v1;
#    - option go_package = "github.com/paper-board/proto/gen/go/<svc>/v1;<svc>v1";

# 3. Validate and generate
buf lint
buf breaking --against '.git#branch=main'
buf generate

# 4. Commit proto sources and generated output together
git add <svc>/v1/ gen/
git commit -m "feat(<svc>): add v1 contract"

# 5. Tag a new semver release (required — see Operations)
git tag v<NEXT>
git push origin v<NEXT>
```

See [Operations](operations.md) for the tagging requirement before any consumer service can pin the new package.

## Standards

Contracts follow [http-api-conventions.md](https://github.com/paper-board/.github/blob/main/docs/standards/http-api-conventions.md) (HTTP gateway annotations) and [ADR-0005](https://github.com/paper-board/.github/blob/main/docs/adr/0005-rest-public-grpc-internal.md) (REST public / gRPC internal split).
