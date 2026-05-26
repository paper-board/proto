# paper-board/proto

gRPC + Protobuf contracts for paperboard. Phase 4 (v0.4.0). OpenAPI auto-generated from proto annotations.

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![buf lint](https://img.shields.io/badge/buf-lint-passing-brightgreen)](https://buf.build/paper-board/proto)
[![buf breaking](https://img.shields.io/badge/buf-breaking-enforced-orange)](https://buf.build/paper-board/proto)

Cross-service contract source of truth — `.proto` files compiled to Go + TypeScript clients.

## Layout

```text
.
├── common/v1/          — shared types: errors, paging, ids, timestamps, usage
├── identity/v1/        — AuthService (Phase 2 ✅)
├── agents/v1/          — AgentService, SessionService (Phase 1.1 ✅)
├── runtime/v1/         — RuntimeService (Phase 3 ✅)
├── compute/v1/         — ComputeService (Phase 3 ✅)
├── audit/v1/           — AuditService (Phase 4)
├── metering/v1/        — MeteringService (Phase 4)
├── notifications/v1/   — NotificationsService (Phase 4)
├── onboarding/v1/      — OnboardingService (Phase 4)
├── environments/v1/    — EnvironmentsService (Phase 4)
├── vaults/v1/          — VaultsService (Phase 4)
├── events/v1/          — Outbox Envelope + 8 Phase 4 event payloads (Phase 4)
├── billing/v1/         — BillingService (Phase 6)
└── gen/                — CI-committed codegen output (Go + TypeScript)
```

## Consuming this module

```bash
go get github.com/paper-board/proto@v0.4.0
```

### GOPRIVATE + authentication (required for sibling services)

All paper-board backend services (`agents`, `identity`, `runtime`, etc.) are private repos.
Their `go.mod` files list `github.com/paper-board/proto` as a dependency. Because sibling repos
are private, Go's module proxy cannot serve them — **you must set `GOPRIVATE`**:

```bash
# local development
export GOPRIVATE=github.com/paper-board/*
export GONOSUMCHECK=github.com/paper-board/*

# also add credentials to ~/.netrc so git fetches over HTTPS succeed:
# machine github.com login <your-github-user> password <personal-access-token>
```

In CI, set `GOPRIVATE=github.com/paper-board/*` as an env var and supply `GITHUB_TOKEN`
(see `paper-board/.github/go-ci.yml` for the canonical job template).

> **DO NOT commit `replace github.com/paper-board/proto => ../proto`.**
> This directive works locally when repos are checked out side-by-side, but it breaks CI
> because the private sibling paths are not present in the runner. Always use a proper
> tagged version (`@v0.4.0`, etc.) in committed `go.mod` files.

## Tagging discipline

This repo uses **manual git tags** — it is NOT managed by release-please (unlike `paper-board/sdk`).

Every time you add a new proto package or make a breaking change:

```bash
git tag v0.X.Y
git push origin v0.X.Y
```

Tag bump rules:

| Change                                                    | Bump                                                                              |
| --------------------------------------------------------- | --------------------------------------------------------------------------------- |
| New package or additive RPC                               | MINOR                                                                             |
| Bugfix to generated output, no proto change               | PATCH                                                                             |
| Breaking change (new package version, e.g. `identity.v2`) | coordinate with all consumers before tagging; `buf breaking` enforces this on PRs |

Tag history:

| Tag    | Contents added                                                                     |
| ------ | ---------------------------------------------------------------------------------- |
| v0.4.0 | Outbox envelope + 6 Phase 4 service contracts + 8 event payloads                   |
| v0.3.0 | `common.v1` UsageMetrics, `compute.v1` ComputeService, `runtime.v1` RuntimeService |
| v0.2.1 | `common.v1` package rename (`paperboard.common.v1` → `common.v1`)                  |
| v0.2.0 | `identity.v1` AuthService                                                          |
| v0.1.0 | Initial scaffold, `common.v1` (errors, paging, ids, timestamps)                    |

## Versioning

Major version in package name: `paperboard.identity.v1`, `paperboard.identity.v2`.
Breaking change → new package, parallel support.

`buf breaking` enforces backward compatibility on PRs.

## buf

Install buf (v2): <https://buf.build/docs/installation>

```bash
# lint
buf lint

# check for breaking changes against main
buf breaking --against '.git#branch=main'

# codegen
buf generate
```

Config files:

- [`buf.yaml`](buf.yaml) — module definition + lint/breaking rules
- [`buf.gen.yaml`](buf.gen.yaml) — plugin config (Go + TypeScript codegen)

Outputs:

- Go: `gen/go/<package>/v1/...` → import path `github.com/paper-board/proto/gen/go/...`
- TypeScript: `gen/ts/<package>/v1/...`

CI verifies `gen/` stays in sync with `.proto` sources.

## Adding a new service package

```bash
# 1. create the package dir
mkdir -p <svc>/v1

# 2. add <svc>/v1/<svc>.proto (see common/v1 for shape)

# 3. validate + generate
buf lint
buf breaking --against '.git#branch=main'
buf generate

# 4. commit both .proto and gen/
git add <svc>/v1/ gen/
git commit -m "feat(<svc>): add v1 contract"

# 5. tag + push (manual — not release-please)
git tag v0.X.Y
git push origin v0.X.Y
```

## Standards

Contract design follows
[http-api-conventions.md](https://github.com/paper-board/.github/blob/main/docs/standards/http-api-conventions.md)
(HTTP annotations) and
[ADR-0005](https://github.com/paper-board/.github/blob/main/docs/adr/0005-rest-public-grpc-internal.md)
(REST public, gRPC internal split).

## Further Reading

- [ADR-0005](https://github.com/paper-board/.github/blob/main/docs/adr/0005-rest-public-grpc-internal.md) — REST public, gRPC internal
- [ADR-0008](https://github.com/paper-board/.github/blob/main/docs/adr/0008-license-coc-commit-conventions.md) — license + versioning + commit conventions
- [docs/operations.md](docs/operations.md) — tagging procedure, tag history, codegen workflow

## License

MIT — see [LICENSE](LICENSE) and
[ADR-0008](https://github.com/paper-board/.github/blob/main/docs/adr/0008-license-coc-commit-conventions.md)
for licensing policy.

[Code of Conduct](https://github.com/paper-board/.github/blob/main/CODE_OF_CONDUCT.md) — Contributor Covenant v2.1.
