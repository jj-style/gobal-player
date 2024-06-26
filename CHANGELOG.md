# Changelog

## [v0.1.20](https://github.com/jj-style/gobal-player/compare/v0.1.19...HEAD) (2024-04-13)

### Fixes

* **workflow:** GH actions write permissions
([fbfe7c2](https://github.com/jj-style/gobal-player/commit/fbfe7c2ff436e04d2e46e7fd23113e8492507444))

### [v0.1.19](https://github.com/jj-style/gobal-player/compare/v0.1.18...v0.1.19) (2024-04-13)

#### Features

* **rss:** add enclosure length
([b158b96](https://github.com/jj-style/gobal-player/commit/b158b963e12ac998dd25a711346a761dc6415c0e))
* add cron schedule to gobal player
([0153e5d](https://github.com/jj-style/gobal-player/commit/0153e5df3a293e3a3e7c55265919b4c06f82d85c)),
closes [#22](https://github.com/jj-style/gobal-player/issues/22)
* **server:** Add server side UI for browsing APIs
([9aa9fa2](https://github.com/jj-style/gobal-player/commit/9aa9fa2748fb59aaceb871f9db53a862d57a0eff)),
closes [#21](https://github.com/jj-style/gobal-player/issues/21)

#### Fixes

* **rss:** timeout when fetching episode content length
([01f59bc](https://github.com/jj-style/gobal-player/commit/01f59bc4ee656b2e71f35cd9f2835a76ea2f555f))
* **server:** actually update rest client with new build id
([39c70a9](https://github.com/jj-style/gobal-player/commit/39c70a92cfe5b627f9bf9fea1dbacc7e56cafcac))
* **dockerfile:** copy html templates in to run layer
([06fc0cf](https://github.com/jj-style/gobal-player/commit/06fc0cf6ab3397b1fdab2b0b926ed517eefc4d8e))

### [v0.1.18](https://github.com/jj-style/gobal-player/compare/v0.1.17...v0.1.18) (2024-03-31)

#### Fixes

* **ci:** changelog by fetch depth all
([6cdc690](https://github.com/jj-style/gobal-player/commit/6cdc6902f766ce57684fd8e90579ea36733001a4))

### [v0.1.17](https://github.com/jj-style/gobal-player/compare/v0.1.16...v0.1.17) (2024-03-31)

#### Fixes

* **ci:** debug release notes
([e20a8db](https://github.com/jj-style/gobal-player/commit/e20a8dbd5ed2aef1ba0d157dd7d04fa8e9669d13))

### [v0.1.16](https://github.com/jj-style/gobal-player/compare/v0.1.15...v0.1.16) (2024-03-31)

#### Fixes

* test release
([f1e4e62](https://github.com/jj-style/gobal-player/commit/f1e4e621c3f998f4d49c1ec813752eb37f70cc6c))

### [v0.1.15](https://github.com/jj-style/gobal-player/compare/v0.1.14...v0.1.15) (2024-03-31)

#### Features

* **server:** RSS feeds for shows and episodes (#20)
([ea8b366](https://github.com/jj-style/gobal-player/commit/ea8b3662a5ec915ab3d6f77fe57e37c258cd3526)),
closes [#20](https://github.com/jj-style/gobal-player/issues/20)
[#17](https://github.com/jj-style/gobal-player/issues/17)

#### Fixes

* **ci:** convco local path to binary
([f2b35d4](https://github.com/jj-style/gobal-player/commit/f2b35d4ef8d57f7557da50f385bc2e69ea76afdd))

### [v0.1.14](https://github.com/jj-style/gobal-player/compare/v0.1.13...v0.1.14) (2024-03-31)

#### Fixes

* **ci:** convco version missing "v" prefix
([a4a3f42](https://github.com/jj-style/gobal-player/commit/a4a3f42a24e8e6586949460f86f27ff2c8357174))

### [v0.1.13](https://github.com/jj-style/gobal-player/compare/v0.1.12...v0.1.13) (2024-03-31)

#### Fixes

* **ci:** convco changelog for released version
([41acff8](https://github.com/jj-style/gobal-player/commit/41acff8c46c254811ece31ea4dfb219dffafade3)),
closes [#19](https://github.com/jj-style/gobal-player/issues/19)

### [v0.1.12](https://github.com/jj-style/gobal-player/compare/v0.1.11...v0.1.12) (2024-03-31)

#### Features

* **server:** RSS feeds for shows and episodes
([7fd44c3](https://github.com/jj-style/gobal-player/commit/7fd44c37336df966d6f271557c1201d57c718bab)),
closes [#17](https://github.com/jj-style/gobal-player/issues/17)

### [v0.1.11](https://github.com/jj-style/gobal-player/compare/v0.1.10...v0.1.11) (2024-03-28)

#### Fixes

* **docker:** fix docker image name for release
([c0f829a](https://github.com/jj-style/gobal-player/commit/c0f829a65afb5f1a69e92748e018d9185a7800b5))

### [v0.1.10](https://github.com/jj-style/gobal-player/compare/v0.1.9...v0.1.10) (2024-03-28)

### [v0.1.9](https://github.com/jj-style/gobal-player/compare/v0.1.8...v0.1.9) (2024-03-28)

#### Features

* **server:** implement server to global player (#16)
([469a7aa](https://github.com/jj-style/gobal-player/commit/469a7aa72a6dcd0b8248070da810c9a097f51b15)),
closes [#16](https://github.com/jj-style/gobal-player/issues/16)
* **tui:** add shortcut to clear cache and refetch data
([8f57aef](https://github.com/jj-style/gobal-player/commit/8f57aef220df3b6799ba371e8961dcc2ac4474de))

### [v0.1.8](https://github.com/jj-style/gobal-player/compare/v0.1.7...v0.1.8) (2024-03-20)

#### Features

* **tui:** set preferred station in config
([de5040e](https://github.com/jj-style/gobal-player/commit/de5040eb7799b98f85b0a6f8107b22dcf7511a85)),
closes [#13](https://github.com/jj-style/gobal-player/issues/13)

### [v0.1.7](https://github.com/jj-style/gobal-player/compare/v0.1.6...v0.1.7) (2024-03-13)

### [v0.1.6](https://github.com/jj-style/gobal-player/compare/v0.1.5...v0.1.6) (2024-03-12)

#### Features

* **tui:** Add keyboard shortcut help text panel
([021f8c8](https://github.com/jj-style/gobal-player/commit/021f8c8bc193c7f98f323faf197737ca02cb69c7)),
closes [#3](https://github.com/jj-style/gobal-player/issues/3)

#### Fixes

* **tests:** fix go list running in docker
([321ef11](https://github.com/jj-style/gobal-player/commit/321ef1122c1ed71f2c99987eedd5573beb370a7c))
* **tests:** remove mocks package from coverage report
([ea8bc1e](https://github.com/jj-style/gobal-player/commit/ea8bc1e715d8381eb2f5714175f2bd3fe4966383)),
closes [#10](https://github.com/jj-style/gobal-player/issues/10)

### [v0.1.5](https://github.com/jj-style/gobal-player/compare/v0.1.4...v0.1.5) (2024-03-09)

#### Fixes

* **release-ci:** add write permissions to pipeline
([fbe6383](https://github.com/jj-style/gobal-player/commit/fbe6383f7e9c69ff3fece52e645165471f72beed))

### [v0.1.4](https://github.com/jj-style/gobal-player/compare/v0.1.3...v0.1.4) (2024-03-09)

#### Fixes

* **release-ci:** fix release notes body path
([6009a31](https://github.com/jj-style/gobal-player/commit/6009a3105b70ed13875d101e15a0211d50f7ec9b))

### [v0.1.3](https://github.com/jj-style/gobal-player/compare/v0.1.2...v0.1.3) (2024-03-09)

#### Fixes

* **release-ci:** only run on v-tag
([52158be](https://github.com/jj-style/gobal-player/commit/52158be5443cde230e36b603e4a323a17577d2c9))

### v0.1.2 (2024-03-09)

#### Fixes

* **justfile:** annotated tag on git tag
([c85da80](https://github.com/jj-style/gobal-player/commit/c85da803437acd043177460d73ea8f688430be06))
* **readme:** add newline to readme
([e74b3bb](https://github.com/jj-style/gobal-player/commit/e74b3bbc75b5cc4da3e53bfb507fe1b331e691d9))
