
### Bug Fixes
- use git switch in release workflow to avoid branch/file ambiguity by @somaz94

### Build
- bump golang in the docker-minor group by @dependabot[bot]
- bump docker/build-push-action from 6 to 7 by @dependabot[bot]
- bump docker/setup-buildx-action from 3 to 4 by @dependabot[bot]

### Features
- add variable interpolation, file input support, and output validation by @somaz94

### Refactoring
- modernize CI/CD workflows and simplify smoke tests by @somaz94
- extract shared jsonutil package, simplify test helpers, and improve test coverage to 98% by @somaz94
- fix error handling, optimize whitespace, unify comments to English by @somaz94
- remove emojis, refactor main for testability, add Makefile by @somaz94

**Full Changelog**: https://github.com/somaz94/env-output-setter/compare/v1.5.1...v1.6.0
