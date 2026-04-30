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

This repo is **skeleton-only** for Phase 1.0b — sadece `common/v1/*.proto` dolu, servis paketleri Phase 4'te doldurulacak. Versioning policy + buf config + CI yapısı şimdiden kuruluyor ki Phase 4'te servis ekleme ucuz olsun.

## License

[MIT](./LICENSE).
