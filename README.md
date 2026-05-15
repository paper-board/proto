# paper-board/proto

**Phase:** 2 — `common/v1` + `identity/v1` (AuthService) shipped at v0.2.1.

Cross-service contract source of truth — `.proto` files compiled to Go + TypeScript clients.

## Versioning

Major version in package name: `paperboard.identity.v1`, `paperboard.identity.v2`. Breaking change → new package, parallel support.

`buf breaking` enforces forward compatibility on PRs.

## Layout

```text
.
├── identity/v1/          (Phase 4: auth, tenant, rbac)
├── agents/v1/            (Phase 4: agents, sessions, prompt)
├── billing/v1/           (Phase 4: subscription, budget, meter)
├── platform/v1/          (Phase 4: audit, notify, webhook)
├── runtime/v1/           (Phase 5: runtime)
├── compute/v1/           (Phase 6: exec)
├── common/v1/            (errors, paging, ids, timestamps)
└── gen/                  (CI-committed codegen output)
```

## Prerequisites

[Install `buf`](https://buf.build/docs/installation) (Mac/Linux/Windows). All commands below assume `buf` is on `PATH`.

## Codegen

```bash
buf generate
```

Outputs:
- Go: `gen/go/<package>/v1/...` → `github.com/paper-board/proto/gen/go/...`
- TS: `gen/ts/<package>/v1/...`

CI verifies `gen/` stays in sync with `.proto` sources.

## Adding a new service package

```bash
# 1. Create the package dir
mkdir -p <svc>/v1

# 2. Add <svc>/v1/<svc>.proto (see common/v1 for shape)

# 3. Validate + generate
buf lint
buf breaking --against '.git#branch=main'
buf generate

# 4. Commit both .proto and gen/
git add <svc>/v1/ gen/
git commit -m "feat(<svc>): add v1 contract"
```

`buf breaking` is enforced in CI — incompatible changes fail PR checks.

## Status

`common/v1` (errors, paging, ids, timestamps) and `identity/v1` (AuthService) are stable at v0.2.1. Other service packages fill in by Phase 4 (see [ADR-0009](https://github.com/paper-board/.github/blob/main/docs/adr/0009-product-first-sequencing.md)). Versioning policy + buf config + CI scaffolding set up early so adding a service in Phase 4 is cheap.

## Standards

Contract design follows [http-api-conventions.md](https://github.com/paper-board/.github/blob/main/docs/standards/http-api-conventions.md) (HTTP annotations) and [ADR-0005](https://github.com/paper-board/.github/blob/main/docs/adr/0005-rest-public-grpc-internal.md) (REST + gRPC split).

## Further Reading

- [ADR-0005](https://github.com/paper-board/.github/blob/main/docs/adr/0005-rest-public-grpc-internal.md) — REST public, gRPC internal
- [ADR-0008](https://github.com/paper-board/.github/blob/main/docs/adr/0008-license-coc-commit-conventions.md) — license + versioning

## License

MIT — see [ADR-0008](https://github.com/paper-board/.github/blob/main/docs/adr/0008-license-coc-commit-conventions.md) for licensing policy. `LICENSE` file lands Phase 5.
