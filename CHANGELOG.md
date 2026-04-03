# Changelog

All notable changes to this project will be documented in this file.

## [v1.7.0](https://github.com/somaz94/env-output-setter/compare/v1.6.1...v1.7.0) (2026-04-03)

### Bug Fixes

- apache license -> mit license ([6829e78](https://github.com/somaz94/env-output-setter/commit/6829e78c185170f4a4b3157021d58bfbbc03469e))
- skip major version tag deletion on first release ([063dfc0](https://github.com/somaz94/env-output-setter/commit/063dfc0f619383631b51d5caa4cb5ac17d8e3f3f))

### Code Refactoring

- extract shared transform logic, remove unused groupPrefix param, and use random EOF delimiter ([455d45b](https://github.com/somaz94/env-output-setter/commit/455d45bdc3f1805427ae2e799941982f7ef66f48))

### Documentation

- remove duplicate rules covered by global CLAUDE.md ([e77f407](https://github.com/somaz94/env-output-setter/commit/e77f40771259b63d6e7cbc5a9e1f7ad614dc65ee))
- add no-push rule to CLAUDE.md ([dbfa35c](https://github.com/somaz94/env-output-setter/commit/dbfa35cac1202da04757d46a8b28f827d4bd3407))
- update CLAUDE.md with commit guidelines and language ([a1c3ae2](https://github.com/somaz94/env-output-setter/commit/a1c3ae271a7ca43f1d7dc4004752fa5136315d2d))
- CLAUDE.md ([ee2cc1a](https://github.com/somaz94/env-output-setter/commit/ee2cc1a59bc5e371c1480926dae1f142693af3ef))

### Continuous Integration

- skip auto-generated changelog and contributors commits in release notes ([49991de](https://github.com/somaz94/env-output-setter/commit/49991de064a46d623797a5a01b53b24bf2319ff2))
- revert to body_path RELEASE.md in release workflow ([68c8e06](https://github.com/somaz94/env-output-setter/commit/68c8e06ddb4ff2202b2239cc7f1cee38efb78a4e))
- use generate_release_notes instead of RELEASE.md ([379d213](https://github.com/somaz94/env-output-setter/commit/379d21391cadb5d4d9e1d8f7dbd2019f2fc4fac4))
- migrate gitlab-mirror workflow to multi-git-mirror action ([9306813](https://github.com/somaz94/env-output-setter/commit/9306813f68c1fb22924830336d671cd4e2b8e201))
- use somaz94/contributors-action@v1 for contributors generation ([bdead1c](https://github.com/somaz94/env-output-setter/commit/bdead1c2414e480f1bb6e07f4a26ac9aec96a160))
- use major-tag-action for version tag updates ([512ee6c](https://github.com/somaz94/env-output-setter/commit/512ee6c09aabc76c6054d48e1f8689758384afee))
- migrate changelog generator to go-changelog-action ([a92b3ee](https://github.com/somaz94/env-output-setter/commit/a92b3eea2d7f9423237e36c8c5cfa9a0ccbe0411))
- add contributors and dependabot auto-merge workflows ([4787dcc](https://github.com/somaz94/env-output-setter/commit/4787dcc260c6c9850d0bc3d3a359282809e543c8))
- update Go version from 1.23 to 1.26 ([699ec24](https://github.com/somaz94/env-output-setter/commit/699ec2426f45f43eb0aa970ed384a02da2446849))
- unify changelog-generator with flexible tag pattern ([d072250](https://github.com/somaz94/env-output-setter/commit/d072250847f24f6b4089e4b5185095b84736f77b))
- use conventional commit message in changelog-generator workflow ([29626a7](https://github.com/somaz94/env-output-setter/commit/29626a7b1fccba2d7f0a7efac7dd32a9a09cc9ac))

### Chores

- remove duplicate rules from CLAUDE.md (moved to global) ([920fd9f](https://github.com/somaz94/env-output-setter/commit/920fd9f9c3a053c130fafd003a4659db662af7fe))
- add git config protection to CLAUDE.md ([d1695fc](https://github.com/somaz94/env-output-setter/commit/d1695fc55bd6fcbe33a500cab5fbc6a6cf0ebc87))
- upgrade Go version to 1.26 ([4e56fda](https://github.com/somaz94/env-output-setter/commit/4e56fda6c947963b514df99f066e4fc9c968474f))
- change license from MIT to Apache 2.0 ([bed0dbf](https://github.com/somaz94/env-output-setter/commit/bed0dbf48bdbafc2348f26cd0e41e3e0c85e3210))
- remove linter workflow and config files ([b1e5c47](https://github.com/somaz94/env-output-setter/commit/b1e5c475ffb52480bb2dbf3b4cfa9d9568135278))
- migrate devcontainer feature from devcontainers-contrib to devcontainers-extra ([c630669](https://github.com/somaz94/env-output-setter/commit/c630669d59d56bd2db47a6eaeac6235833d1085d))
- switch changelog generator from git-cliff to github_changelog_generator ([2818e83](https://github.com/somaz94/env-output-setter/commit/2818e836056836051802a58b769b5683a24fb0e9))
- regenerate CHANGELOG.md with version-based format ([0b83e9a](https://github.com/somaz94/env-output-setter/commit/0b83e9a01140d5d04c179abeccfe59703b1a7073))
- regenerate changelog and revert cliff.toml to simple format ([af850c8](https://github.com/somaz94/env-output-setter/commit/af850c83edd7a6bfc192f68fb0bf130486beb695))
- update cliff.toml for version-based changelog format ([c757c16](https://github.com/somaz94/env-output-setter/commit/c757c1679a3d3fc91cf4909f8220774168d43621))

### Contributors

- somaz

<br/>

## [v1.6.1](https://github.com/somaz94/env-output-setter/compare/v1.6.0...v1.6.1) (2026-03-09)

### Bug Fixes

- preserve JSON values with commas during delimiter split ([acc826b](https://github.com/somaz94/env-output-setter/commit/acc826b459f4888bef164064c361688454279d26))

### Contributors

- somaz

<br/>

## [v1.6.0](https://github.com/somaz94/env-output-setter/compare/v1.5.1...v1.6.0) (2026-03-09)

### Features

- add variable interpolation, file input support, and output validation ([213e4ca](https://github.com/somaz94/env-output-setter/commit/213e4cae0aa7b000334e555079d935389df63a32))

### Bug Fixes

- use git switch in release workflow to avoid branch/file ambiguity ([71ff95d](https://github.com/somaz94/env-output-setter/commit/71ff95dc43dbd2bfc92a44fa9d4b17b512dcf893))

### Code Refactoring

- remove emojis, refactor main for testability, add Makefile ([5f7ba05](https://github.com/somaz94/env-output-setter/commit/5f7ba05a2e54365433df51a35aee7778b72eed27))
- fix error handling, optimize whitespace, unify comments to English ([4ce0714](https://github.com/somaz94/env-output-setter/commit/4ce0714b03ca07ac8514cf1c92ee9990993f84b8))
- extract shared jsonutil package, simplify test helpers, and improve test coverage to 98% ([571733f](https://github.com/somaz94/env-output-setter/commit/571733fa0b8b646c8d5accf1e9c3b950b7198b6b))
- modernize CI/CD workflows and simplify smoke tests ([cdd73a4](https://github.com/somaz94/env-output-setter/commit/cdd73a452ba7f5a9b16932971d39e71695ae3c04))

### Builds

- **deps:** bump docker/setup-buildx-action from 3 to 4 ([4667b59](https://github.com/somaz94/env-output-setter/commit/4667b595fb968033db98467d8bf966c9d6eab95f))
- **deps:** bump docker/build-push-action from 6 to 7 ([bdea599](https://github.com/somaz94/env-output-setter/commit/bdea599d66eb9486001304bb0380a917f59f9bfe))
- **deps:** bump golang in the docker-minor group ([24b8c40](https://github.com/somaz94/env-output-setter/commit/24b8c40c29aea2f263e6688ad2c82f21063a77fe))

### Contributors

- somaz

<br/>

## [v1.5.1](https://github.com/somaz94/env-output-setter/compare/v1.5.0...v1.5.1) (2025-11-27)

### Builds

- **deps:** bump actions/checkout from 5 to 6 ([ec1b77f](https://github.com/somaz94/env-output-setter/commit/ec1b77fc9e0aace4a4a6d745796675772a237849))
- **deps:** bump actions/setup-go from 5 to 6 ([408222e](https://github.com/somaz94/env-output-setter/commit/408222e06f3241eac84eefc1f4b3e9c0c6a7b5fc))

### Chores

- test code ([af95448](https://github.com/somaz94/env-output-setter/commit/af954486c03dc3c98a2a977df27bb3f434fa62a0))
- stale-issues, issue-greeting ([0698bb4](https://github.com/somaz94/env-output-setter/commit/0698bb4eb3374d1f7150f6fc7017ecd4ccb21dc5))
- dockerignore ([6eec257](https://github.com/somaz94/env-output-setter/commit/6eec257d19bc46ac77b4ba56cf03fa6f755c676b))
- release.yml ([d4d9ec7](https://github.com/somaz94/env-output-setter/commit/d4d9ec73ad90497b90b520f5a13436f4beb69c44))
- workflows ([01aefb0](https://github.com/somaz94/env-output-setter/commit/01aefb086196d8aaa9fb493e053e852ab73c7ad2))

### Contributors

- somaz

<br/>

## [v1.5.0](https://github.com/somaz94/env-output-setter/compare/v1.4.1...v1.5.0) (2025-10-30)

### Code Refactoring

- all ([7b2dd69](https://github.com/somaz94/env-output-setter/commit/7b2dd69430b2b66766af6b658dfe53892caca4f8))
- ci.yml,action.yml,writer.go ([9bbda7b](https://github.com/somaz94/env-output-setter/commit/9bbda7b9f5b364b0296c30013e6f08795009b6a9))
- all ([648e212](https://github.com/somaz94/env-output-setter/commit/648e21209da1955c0847c0754f2fbee278cbdad2))
- cmd/main.go ([1003061](https://github.com/somaz94/env-output-setter/commit/10030613d82c3185cfde0177387646e0e2b6d9bc))

### Builds

- **deps:** bump actions/checkout from 4 to 5 ([268c491](https://github.com/somaz94/env-output-setter/commit/268c491a907d6fa43737fb23412935ad37762ece))
- **deps:** bump golang in the docker-minor group ([874fd1d](https://github.com/somaz94/env-output-setter/commit/874fd1d2513a7cc406d85959b1bd5e5b5e58d518))
- **deps:** bump super-linter/super-linter from 7 to 8 ([a411d4c](https://github.com/somaz94/env-output-setter/commit/a411d4caab5956e535d25e9bed8eaf0339464518))

### Chores

- writer.go ([94b7176](https://github.com/somaz94/env-output-setter/commit/94b7176ae7b7f7dcfcc5643f33d5080a3295e917))

### Contributors

- somaz

<br/>

## [v1.4.1](https://github.com/somaz94/env-output-setter/compare/v1.4.0...v1.4.1) (2025-04-15)

### Bug Fixes

- config.go, printer.go ([142a880](https://github.com/somaz94/env-output-setter/commit/142a880c2189237a32afe040bd084dc3a53792fd))
- transformer.go ([04e9273](https://github.com/somaz94/env-output-setter/commit/04e9273bbc175fd09f50c091f90b3d9d83a369e4))
- writer.go ([03ff523](https://github.com/somaz94/env-output-setter/commit/03ff523bf802e217c2219d180dfa82b52a6ed22f))

### Documentation

- README.md ([96d0967](https://github.com/somaz94/env-output-setter/commit/96d09676c444efee6e011069e0883d48f6862306))
- README.md ([a8e6340](https://github.com/somaz94/env-output-setter/commit/a8e6340b152f7efd08135e9797932450915867ed))

### Contributors

- somaz

<br/>

## [v1.4.0](https://github.com/somaz94/env-output-setter/compare/v1.3.1...v1.4.0) (2025-04-10)

### Bug Fixes

- writer.go ([ca33ecb](https://github.com/somaz94/env-output-setter/commit/ca33ecb811488ceb6277bf302425bdb4b7000b0b))
- transformer.go, writer.go ([e86e2d5](https://github.com/somaz94/env-output-setter/commit/e86e2d5a247c252bc8322dac247aa23bb4dc0ce9))
- ci.yml, writer.go ([2ef5ee6](https://github.com/somaz94/env-output-setter/commit/2ef5ee6518013fd89e368c96f6cc4acf57c029de))
- ci.yml, README.md ([8953f47](https://github.com/somaz94/env-output-setter/commit/8953f47637d4719b141fbce4a3a897396c5f3942))
- ci.yml ([649d0e4](https://github.com/somaz94/env-output-setter/commit/649d0e41efb4e1f1e09a5d81ef327c6399b907f7))
- ci.yml ([1cf4fef](https://github.com/somaz94/env-output-setter/commit/1cf4fef0fa86190aa34b3fec453f3bf4162f942e))
- ci.yml, action.yml, config,go, transformer.go, writer.go ([5f67244](https://github.com/somaz94/env-output-setter/commit/5f67244e576fa860ec5620b5c091a129a73b620e))

### Contributors

- somaz

<br/>

## [v1.3.1](https://github.com/somaz94/env-output-setter/compare/v1.3.0...v1.3.1) (2025-03-04)

### Bug Fixes

- config, printer, transformer.go ([9797994](https://github.com/somaz94/env-output-setter/commit/97979943926182dd83c414bc6086183f2b3528c8))
- config, printer, transformer.go ([86b1a15](https://github.com/somaz94/env-output-setter/commit/86b1a156c42d0074b4171f667e8c589e5ab70673))
- backup/config, printer, transformer ([5de21a4](https://github.com/somaz94/env-output-setter/commit/5de21a41a1317e55e84d53e4cbf760ebdff4be55))
- writer.go ([cab3a5d](https://github.com/somaz94/env-output-setter/commit/cab3a5de67c8decffd7c659718419b2e9c36b14b))
- writer.go ([685152a](https://github.com/somaz94/env-output-setter/commit/685152a50399161f126bd665e822aa7fefe8eec4))
- backup/writer.go ([216c824](https://github.com/somaz94/env-output-setter/commit/216c824d2e1dee28387f330781b560ff3114fb02))
- changelog-generator.yml ([c02e6b0](https://github.com/somaz94/env-output-setter/commit/c02e6b0cc0a41642a6785f77c5a1a1981eca046f))

### Documentation

- README.md ([d21487f](https://github.com/somaz94/env-output-setter/commit/d21487f7432224fed9c58f783470102ab2b2aad2))
- README.md ([5202efa](https://github.com/somaz94/env-output-setter/commit/5202efa7f4d3580d11b307a0d6cab837086f9e26))

### Add

- gitlab-mirror.yml ([1b34602](https://github.com/somaz94/env-output-setter/commit/1b34602ac34606e2001d322b8a782a1d08410adc))

### Contributors

- somaz

<br/>

## [v1.3.0](https://github.com/somaz94/env-output-setter/compare/v1.2.1...v1.3.0) (2025-02-20)

### Bug Fixes

- printer.go, writer,go, use-action.yml ([b0a9d8d](https://github.com/somaz94/env-output-setter/commit/b0a9d8d874a72c8ce30191c292b0317ffafb3248))
- writer.go ([ed16eea](https://github.com/somaz94/env-output-setter/commit/ed16eeaf1986374b9c415fec73a16dcb4f40e8d4))
- ci.yml, action.yml, config.go, writer.go ([e7d15b7](https://github.com/somaz94/env-output-setter/commit/e7d15b76ca284b4609664786e7f9db33b72adb2c))
- writer.go ([6b0d6a1](https://github.com/somaz94/env-output-setter/commit/6b0d6a1cc38148dc1dd11cc03380309108393504))
- writer.go ([07769b1](https://github.com/somaz94/env-output-setter/commit/07769b19a312ace278e4bad6a1348454ef3c90f3))
- ci.yml, writer.go ([805bacd](https://github.com/somaz94/env-output-setter/commit/805bacd62d578e7ae43fbf97f74116ff69ac2c4f))
- ci.yml ([b0e5bdf](https://github.com/somaz94/env-output-setter/commit/b0e5bdf3a6f44b2bdb769f0476642ec2dc3ba424))
- wrtier.go, README.md ([81d036a](https://github.com/somaz94/env-output-setter/commit/81d036a22c2b85f8c27967d0a1ba279edb0ce67c))
- ci.yml ([405c241](https://github.com/somaz94/env-output-setter/commit/405c2410508ce4e26da67f3ccfb085247057d7d0))
- writer.go ([8085510](https://github.com/somaz94/env-output-setter/commit/80855103e86d69f7a8f24e50ee3661bbd5066a01))

### Contributors

- somaz

<br/>

## [v1.2.1](https://github.com/somaz94/env-output-setter/compare/v1.2.0...v1.2.1) (2025-02-19)

### Bug Fixes

- main.go, transformer.go, writer.go ([ccfdfbf](https://github.com/somaz94/env-output-setter/commit/ccfdfbfe53b75d9601d935df71bd50e4abc0a283))
- backup/*, ci.yml ([71df3a7](https://github.com/somaz94/env-output-setter/commit/71df3a7fd757553d9d731bc21a365665baa0dfb4))
- dependabot.yml ([cd0aaa1](https://github.com/somaz94/env-output-setter/commit/cd0aaa1ee6b16847321aeba63795029281784f9c))
- changelog-generator.yml ([5b48296](https://github.com/somaz94/env-output-setter/commit/5b482968c96c9e55b5dc450189fa98dc9b2d90b5))

### Documentation

- README.md ([de34aea](https://github.com/somaz94/env-output-setter/commit/de34aea65f00046eedc69cba8a10bb3f982d9762))

### Builds

- **deps:** bump golang in the docker-minor group ([25c152f](https://github.com/somaz94/env-output-setter/commit/25c152f9af40cf4db767b0811e30ccd1bdf7b335))

### Contributors

- somaz

<br/>

## [v1.2.0](https://github.com/somaz94/env-output-setter/compare/v1.1.0...v1.2.0) (2025-02-14)

### Bug Fixes

- action.yml, ci.yml, config.go, transformer.go, writer.go ([fece63c](https://github.com/somaz94/env-output-setter/commit/fece63cd57219fc0fd59de771d02b75027df020c))

### Contributors

- somaz

<br/>

## [v1.1.0](https://github.com/somaz94/env-output-setter/compare/v1.0.4...v1.1.0) (2025-02-13)

### Bug Fixes

- use-action.yml ([5a74fb3](https://github.com/somaz94/env-output-setter/commit/5a74fb34c5cb80078fa0aa71848fd8a12cbfd0ec))

### Documentation

- README.md ([9548514](https://github.com/somaz94/env-output-setter/commit/95485147082f915675cd1c719468e5f5052b2066))
- README.md ([5f68cc4](https://github.com/somaz94/env-output-setter/commit/5f68cc4e903f1b6a35674666585b793b87c36998))
- README.md ([59e189e](https://github.com/somaz94/env-output-setter/commit/59e189e138b4d04b607acb5ef1fc0930b825cd1b))

### Add

- transformer.go ([91e629f](https://github.com/somaz94/env-output-setter/commit/91e629f0ececa67d4bc0ba828ad7a7e716070fdb))

### Contributors

- somaz

<br/>

## [v1.0.4](https://github.com/somaz94/env-output-setter/compare/v1.0.3...v1.0.4) (2025-02-13)

### Bug Fixes

- README.md , ci.yml ([ab1814c](https://github.com/somaz94/env-output-setter/commit/ab1814c0de3ed3b465e5419d9118ba58e85efe96))
- ci.yml ([7ba5834](https://github.com/somaz94/env-output-setter/commit/7ba5834ff781c143dd2cd0e6c4e5ac09470560d0))
- writer.go ([1edb86a](https://github.com/somaz94/env-output-setter/commit/1edb86a0f1ea40e34a36a04dbdd3d232383c23d2))
- main.go ([ae6449d](https://github.com/somaz94/env-output-setter/commit/ae6449d4845fc44d76e71da3d4dd089cbb5af534))
- file structure ([98220a4](https://github.com/somaz94/env-output-setter/commit/98220a4ba63a2cb8be3292b472afbc0f308c45c6))

### Contributors

- somaz

<br/>

## [v1.0.3](https://github.com/somaz94/env-output-setter/compare/v1.0.2...v1.0.3) (2025-02-13)

### Bug Fixes

- file-structure ([aba4dbc](https://github.com/somaz94/env-output-setter/commit/aba4dbc0c669719ab18a5053c40f6ba2533df273))

### Contributors

- somaz

<br/>

## [v1.0.2](https://github.com/somaz94/env-output-setter/compare/v1.0.1...v1.0.2) (2025-02-13)

### Bug Fixes

- Dockerfile ([b582e94](https://github.com/somaz94/env-output-setter/commit/b582e94c07177e3bf24c8adb4e025b4e8d65cbca))

### Documentation

- README.md ([1e82167](https://github.com/somaz94/env-output-setter/commit/1e82167964a0022dfa4cd4690593e9a44a95a195))

### Contributors

- somaz

<br/>

## [v1.0.1](https://github.com/somaz94/env-output-setter/compare/v1.0.0...v1.0.1) (2025-02-07)

### Bug Fixes

- main.go ([a9b080f](https://github.com/somaz94/env-output-setter/commit/a9b080f39c95fa202126d09dfe6530dc5724b826))
- use-action.yml ([a117a56](https://github.com/somaz94/env-output-setter/commit/a117a56e49acacd6341cf5b15ae6cb1b7deaf64c))
- main.go ([e7de71c](https://github.com/somaz94/env-output-setter/commit/e7de71c0059094cf085706ab4eefe225283ae7cf))

### Contributors

- somaz

<br/>

## [v1.0.0](https://github.com/somaz94/env-output-setter/releases/tag/v1.0.0) (2025-02-05)

### Bug Fixes

- ci.yml ([83219d7](https://github.com/somaz94/env-output-setter/commit/83219d70a76da8c6f933881225a5490853a59e72))
- ci.yml ([d90681c](https://github.com/somaz94/env-output-setter/commit/d90681c38608c63cc8f10d4e67bf7e16d1d7fe98))
- linter.yml & README.md ([4df4c45](https://github.com/somaz94/env-output-setter/commit/4df4c45de3f3e66f280f13ed27a9cf5d2411fdff))
- .github/* & README.md & action.yml ([d9a7815](https://github.com/somaz94/env-output-setter/commit/d9a7815db6246376ba1a9b44fc2da38ac182bd4d))
- linter.yml & use-action.yml ([bbfecc8](https://github.com/somaz94/env-output-setter/commit/bbfecc80ab61f7b3269062c6f5d833ecef4d457e))
- Dockerfile ([4db894a](https://github.com/somaz94/env-output-setter/commit/4db894a0a7b852ad37cab1f102f7b8454c862d1d))
- Dockerfile ([db1a10f](https://github.com/somaz94/env-output-setter/commit/db1a10fb1dd60087135482b9bf28f94593fe0469))
- use-action.yml ([37c3906](https://github.com/somaz94/env-output-setter/commit/37c390652e783815420a692bd80f3368d48ee3f1))
- action.yml ([7d9ce59](https://github.com/somaz94/env-output-setter/commit/7d9ce59caa41733811d4e011766611063a593a83))
- action.yml ([8e7b17b](https://github.com/somaz94/env-output-setter/commit/8e7b17b0721e96a8beb001f52afc112d308afa29))
- ci.yml ([43e9893](https://github.com/somaz94/env-output-setter/commit/43e98933b64b9558f08699b3026bdf429bee1e93))
- ci.yml & main.go ([b4c3c39](https://github.com/somaz94/env-output-setter/commit/b4c3c3999c6c343130cc0442085453c03951f899))
- Dockerfile ([cd6a55e](https://github.com/somaz94/env-output-setter/commit/cd6a55efb4bb7bd175e914a59c39c49de51d6ba9))
- ci.yml & Dockerfile ([e73f391](https://github.com/somaz94/env-output-setter/commit/e73f39165a983354b2346abf6c07b30645ca93be))
- ci.yml & Dockerfile ([9ac55a7](https://github.com/somaz94/env-output-setter/commit/9ac55a7ebe0b840c7e9b29c70ef03970098c9d9e))
- Dockerfile & main.go ([8c9f6d3](https://github.com/somaz94/env-output-setter/commit/8c9f6d337ae7b1050df1625cfbaf2cc293a9d305))
- main.go ([cdfddcd](https://github.com/somaz94/env-output-setter/commit/cdfddcdb94d5edf351e66cd4a32ac2af177d6e7a))
- main.go ([442d367](https://github.com/somaz94/env-output-setter/commit/442d36799247623849f17a5be543f89d3531c663))

### Documentation

- CODEOWNERS ([1a55650](https://github.com/somaz94/env-output-setter/commit/1a55650e9495538d8c80a24cba9055a40510139e))
- README.md ([888684e](https://github.com/somaz94/env-output-setter/commit/888684eafd4047df7f9ca369c0959c737a6de6b0))
- README.md ([e997093](https://github.com/somaz94/env-output-setter/commit/e9970932f9a529c6f169202d904c60e653f7fa54))

### Builds

- **deps:** bump janheinrichmerker/action-github-changelog-generator ([2727f8d](https://github.com/somaz94/env-output-setter/commit/2727f8dd6a554ff2d77dc4ca72b495ac9ebcdbf1))
- **deps:** bump golang from 1.20 to 1.23 in the docker-minor group ([8170345](https://github.com/somaz94/env-output-setter/commit/81703451e4fb00867c701edfa04f51205f2d9085))

### Chores

- fix changelog-generator.yml ([ddb246c](https://github.com/somaz94/env-output-setter/commit/ddb246c5b88be831b448eaae4f88aae269e60841))
- fix changelog workflow ([9fe0238](https://github.com/somaz94/env-output-setter/commit/9fe0238a3f87385296ee32e172dd63483f99a8bb))
- add changelog workflow ([627b842](https://github.com/somaz94/env-output-setter/commit/627b8426df0f15245267a7e920069e7bf231a278))

### Contributors

- somaz

<br/>

