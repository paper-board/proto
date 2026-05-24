# Operations

## Tag discipline

Every new proto package addition requires a new semver tag before any consumer service can pin it.

**Rule:** `go get github.com/paper-board/proto@<tag>` in a consumer's CI resolves a module version, not a branch. A branch reference works locally but fails in CI because the Go module proxy requires a tagged release.

### Tag sequence

| Tag    | Contents added                                                                     |
| ------ | ---------------------------------------------------------------------------------- |
| v0.1.0 | Initial repo scaffold, `common.v1` (errors, paging, ids, timestamps)               |
| v0.2.0 | `identity.v1` AuthService                                                          |
| v0.2.1 | `common.v1` package rename (`paperboard.common.v1` → `common.v1`)                  |
| v0.3.0 | `common.v1` UsageMetrics, `compute.v1` ComputeService, `runtime.v1` RuntimeService |

### Tagging procedure

After merging a PR that adds a new package or makes a semver-worthy change:

```bash
git checkout main
git pull
git tag v<NEXT>
git push origin v<NEXT>
```

Increment rules (SemVer 2.0):

- New package or additive RPC → minor bump (`v0.3.0` → `v0.4.0`)
- Bugfix to generated output, no proto change → patch bump
- Breaking change (new package version, e.g. `identity.v2`) → coordinate with all consumers before tagging

`buf breaking` enforced in CI prevents accidental breaking changes in an existing package version. Breaking changes must go in a new package (e.g. `identity/v2/`).

______________________________________________________________________

## GOPRIVATE + GITHUB_TOKEN in CI

`paper-board/proto` is a private repository. Consumer service CI must authenticate to fetch it.

### Required CI configuration

In each consumer service's GitHub Actions workflow:

```yaml
env:
  GOPRIVATE: github.com/paper-board/*
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

steps:
  - name: Configure git for private modules
    run: git config --global url."https://x-access-token:${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
```

This is already centralised in `paper-board/.github/go-ci.yml` (reusable workflow). Services that extend it inherit the configuration automatically.

Without `GOPRIVATE`, the Go toolchain attempts to resolve `github.com/paper-board/proto` via the public module proxy (proxy.golang.org), which returns 404 for private repos. The `insteadOf` git rewrite ensures the module fetch uses the Actions token.

______________________________________________________________________

## Local `replace` directive

During local development, before a new package is tagged, you can point a consumer's `go.mod` at the local checkout:

```go
// go.mod (local dev only)
replace github.com/paper-board/proto => ../proto
```

**Never commit this directive.** CI has no `../proto` sibling directory; the build will fail. The replace must be removed before opening a PR. Verify with:

```bash
grep 'replace.*proto' go.mod  # must return nothing on a PR branch
```

______________________________________________________________________

## CI workflow summary

Three jobs run on every push / PR (see `.github/workflows/ci.yaml`):

| Job        | Trigger   | What it checks                                                            |
| ---------- | --------- | ------------------------------------------------------------------------- |
| `lint`     | push + PR | `buf lint` — STANDARD ruleset, PACKAGE_DIRECTORY_MATCH                    |
| `breaking` | PR only   | `buf breaking --against main` — no field removals, no type changes        |
| `generate` | push + PR | `buf generate` then `git diff --exit-code gen/` — gen/ must match sources |

The `generate` job is the enforcement mechanism for keeping committed generated code in sync. A PR that modifies a `.proto` file without re-running `buf generate` and committing the output will fail.

______________________________________________________________________

## buf configuration notes

`buf.yaml` (v2) at repo root:

- `lint.use: [STANDARD]` — full standard ruleset
- `lint.except: [PACKAGE_VERSION_SUFFIX]` — allows `identity.v1` without requiring a `v1` literal suffix in the package name
- `lint.except: [RPC_RESPONSE_STANDARD_NAME]` — allows response message names that don't end in `Response` where semantically appropriate
- `breaking.use: [FILE]` — file-level breaking change detection
- `breaking.ignore_only.FILE_SAME_PACKAGE: [common/v1]` — common package rename from `paperboard.common.v1` is grandfathered; new violations are not permitted

`buf.gen.yaml` (v2) generates:

- Go stubs via `buf.build/protocolbuffers/go` → `gen/go/`
- gRPC server/client via `buf.build/grpc/go` → `gen/go/`
- HTTP gateway via `buf.build/grpc-ecosystem/gateway` → `gen/go/`
- TypeScript via `buf.build/bufbuild/es` → `gen/ts/`
