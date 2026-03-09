# CLAUDE.md - env-output-setter

GitHub Action that sets multiple key-value pairs in both `$GITHUB_ENV` and `$GITHUB_OUTPUT`.

## Project Structure

```
cmd/main.go                  # Entrypoint
internal/
  config/                    # Input parsing, validation, option structs
  writer/                    # GITHUB_ENV / GITHUB_OUTPUT file writers
  transformer/               # Value transformation (upper, lower, URL encode, mask)
  jsonutil/                  # JSON parsing and flattening
  interpolator/              # Variable interpolation
  printer/                   # Debug/summary output
  filereader/                # File encoding handling
tests/test_local.go          # Local integration test suite
Makefile                     # Build, test, lint commands
Dockerfile                   # Multi-stage (golang:1.26-alpine → alpine:latest)
action.yml                   # GitHub Action definition (23 inputs, 4 outputs)
cliff.toml                   # git-cliff config for release notes
```

## Build & Test

```bash
make test            # Run unit tests (alias for test-unit)
make test-unit       # go test ./internal/... ./cmd/... -v -cover
make test-local      # Run local integration test suite (tests/test_local.go)
make test-all        # Run all tests (unit + local)
make cover           # Generate coverage report
make cover-html      # Open coverage in browser
make bench           # Run benchmarks
make lint            # go vet
make fmt             # gofmt
make build           # Build binary
make clean           # Remove artifacts
```

## Key Inputs

- **Required**: `env_key`, `env_value`, `output_key`, `output_value`
- **Options**: `delimiter`, `fail_on_empty`, `trim_whitespace`, `case_sensitive`, `error_on_duplicate`
- **Transform**: `to_upper`, `to_lower`, `encode_url`, `escape_newlines`, `max_length`, `mask_secrets`, `mask_pattern`
- **Advanced**: `json_support`, `group_prefix`, `export_as_env`, `enable_interpolation`, `validation_rules`

## Outputs

`set_env_count`, `set_output_count`, `action_status`, `error_message`

## Workflow Structure

| Workflow | Name | Trigger |
|----------|------|---------|
| `ci.yml` | `Continuous Integration` | push(main), PR, dispatch |
| `release.yml` | `Create release` | tag push `v*` |
| `changelog-generator.yml` | `Generate changelog` | after release, PR merge, issue close |
| `use-action.yml` | `Smoke Test (Released Action)` | after release, dispatch |
| `linter.yml` | `Lint Codebase` | push(main), PR |

### Workflow Chain
```
tag push v* → Create release
                ├→ Smoke Test (Released Action)
                └→ Generate changelog
```

### CI Structure
```
unit-tests ─┐
test-local ──→ build-and-push-docker → test-action → ... → ci-result
```

## Conventions

- **Commits**: Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `ci:`, `chore:`)
- **Branches**: `main` (production), `test` (integration tests)
- **Secrets**: `PAT_TOKEN` (cross-repo ops), `GITHUB_TOKEN` (changelog, releases)
- **Docker**: Multi-stage build, alpine base
- **Comments**: English only
- **Release**: `git switch` (not `git checkout`), git-cliff for RELEASE.md
- **cliff.toml**: Skip `^Merge`, `^Update changelog`, `^Auto commit`
- **paths-ignore**: `.github/workflows/**`, `**/*.md`, `backup/**`
- Do NOT commit directly - recommend commit messages only
