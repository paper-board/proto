# paper-board/proto

Cross-service contract source of truth — `.proto` files compiled to Go + TypeScript clients.

## Versioning

Major version in package name: `paperboard.identity.v1`, `paperboard.identity.v2`. Breaking change → new package, parallel support.

`buf breaking` enforces forward compatibility on PRs.

## Layout

```
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

## Codegen

```bash
buf generate
```

Outputs:
- Go: `gen/go/<package>/v1/...` → `github.com/paper-board/proto/gen/go/...`
- TS: `gen/ts/<package>/v1/...`

CI verifies `gen/` stays in sync with `.proto` sources.

## Phase 1.0b

This repo is **skeleton-only** for Phase 1.0b — only `common/v1/*.proto` is populated; service packages are filled in Phase 4. Versioning policy + buf config + CI scaffolding is set up now so that adding a service in Phase 4 is cheap.

## License

[MIT](./LICENSE).
