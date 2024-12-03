# Changelog

## [0.1.6](https://github.com/anexia/external-dns-webhook/compare/v0.1.5...v0.1.6) (2024-09-25)


### Miscellaneous Chores

* **deps:** bump anchore/sbom-action from 0.16.1 to 0.17.0 ([#31](https://github.com/anexia/external-dns-webhook/issues/31)) ([b0fea23](https://github.com/anexia/external-dns-webhook/commit/b0fea23766e0b28ad765c8bbad9e90b81e552733))
* **deps:** bump anchore/sbom-action from 0.17.0 to 0.17.1 ([#37](https://github.com/anexia/external-dns-webhook/issues/37)) ([9a87b20](https://github.com/anexia/external-dns-webhook/commit/9a87b20eca191761d8028dd043d5f488fc7caf40))
* **deps:** bump anchore/sbom-action from 0.17.1 to 0.17.2 ([#38](https://github.com/anexia/external-dns-webhook/issues/38)) ([37dc2ff](https://github.com/anexia/external-dns-webhook/commit/37dc2ff2323c7026126e3d40f80dc8c1b1c29b6b))
* **deps:** bump github.com/caarlos0/env/v11 from 11.1.0 to 11.2.2 ([#35](https://github.com/anexia/external-dns-webhook/issues/35)) ([fd33c0a](https://github.com/anexia/external-dns-webhook/commit/fd33c0af66694c0d438277d108bd959042e89758))
* **deps:** bump go.anx.io/go-anxcloud from 0.7.2 to 0.7.3 ([#34](https://github.com/anexia/external-dns-webhook/issues/34)) ([0dc6320](https://github.com/anexia/external-dns-webhook/commit/0dc63205f19ed23c6f7535b96127349a3bb77a39))
* **deps:** bump sigs.k8s.io/external-dns from 0.14.2 to 0.15.0 ([#39](https://github.com/anexia/external-dns-webhook/issues/39)) ([44d35fa](https://github.com/anexia/external-dns-webhook/commit/44d35fa518f0d2310e3d3a1b090bae98d023e8ef))
* **deps:** bump sigstore/cosign-installer from 3.5.0 to 3.6.0 ([#36](https://github.com/anexia/external-dns-webhook/issues/36)) ([2e672d1](https://github.com/anexia/external-dns-webhook/commit/2e672d1efa69dd03c91f26c3e0e64d2c67c97b59))

## [0.1.5](https://github.com/anexia/external-dns-webhook/compare/v0.1.4...v0.1.5) (2024-07-10)


### Bug Fixes

* do not require image signing during pr ([#27](https://github.com/anexia/external-dns-webhook/issues/27)) ([63138df](https://github.com/anexia/external-dns-webhook/commit/63138df80e3add4e6afcbf1e2d52aed41d3e03ef))
* fix typo ([#21](https://github.com/anexia/external-dns-webhook/issues/21)) ([3d30e2e](https://github.com/anexia/external-dns-webhook/commit/3d30e2ec5f6aa6f6baa05aec6b8719a3c7c7d87a))
* set the sign skip correctly in pul_request action ([#29](https://github.com/anexia/external-dns-webhook/issues/29)) ([2811151](https://github.com/anexia/external-dns-webhook/commit/28111518d25d42873d353a615d8d5d94da4938db))


### Miscellaneous Chores

* **deps:** bump anchore/sbom-action from 0.16.0 to 0.16.1 ([#30](https://github.com/anexia/external-dns-webhook/issues/30)) ([c7d40d9](https://github.com/anexia/external-dns-webhook/commit/c7d40d9d5c9b0e5bf1429f83391b81b260b78c74))

## [0.1.4](https://github.com/anexia/external-dns-webhook/compare/v0.1.3...v0.1.4) (2024-07-06)


### Bug Fixes

* update goreleaser config to work with v2 ([#26](https://github.com/anexia/external-dns-webhook/issues/26)) ([2ee51d8](https://github.com/anexia/external-dns-webhook/commit/2ee51d8b16f54cf89af4d2d2b1b8a03fc0143a3b))


### Miscellaneous Chores

* **deps:** bump github.com/caarlos0/env/v11 from 11.0.1 to 11.1.0 ([#20](https://github.com/anexia/external-dns-webhook/issues/20)) ([9357d69](https://github.com/anexia/external-dns-webhook/commit/9357d693f271a04903590a66f9209a3e9b2403a8))
* **deps:** bump github.com/go-chi/chi/v5 from 5.0.12 to 5.1.0 ([#24](https://github.com/anexia/external-dns-webhook/issues/24)) ([0fe0722](https://github.com/anexia/external-dns-webhook/commit/0fe07221b59895b956e1d82a9ff183118ba48c37))
* **deps:** bump go.anx.io/go-anxcloud from 0.6.4 to 0.7.0 ([#15](https://github.com/anexia/external-dns-webhook/issues/15)) ([d5eacbf](https://github.com/anexia/external-dns-webhook/commit/d5eacbf84c3aec027d0a15090f527b6dbe2244d6))
* **deps:** bump go.anx.io/go-anxcloud from 0.7.0 to 0.7.1 ([#18](https://github.com/anexia/external-dns-webhook/issues/18)) ([3c70b4a](https://github.com/anexia/external-dns-webhook/commit/3c70b4ac9055d76c77b78bca4d09f62dbaa3acda))
* **deps:** bump go.anx.io/go-anxcloud from 0.7.1 to 0.7.2 ([#22](https://github.com/anexia/external-dns-webhook/issues/22)) ([938b348](https://github.com/anexia/external-dns-webhook/commit/938b348eaacdae868809e55471599c958348f59c))
* **deps:** bump goreleaser/goreleaser-action from 5 to 6 ([#17](https://github.com/anexia/external-dns-webhook/issues/17)) ([e6e0270](https://github.com/anexia/external-dns-webhook/commit/e6e02708099b8b6de19c6fd5649e7c2e47786b1a))

## [0.1.3](https://github.com/anexia/external-dns-webhook/compare/v0.1.2...v0.1.3) (2024-05-07)


### Bug Fixes

* reorder record parsing in GetRecords function ([#9](https://github.com/anexia/external-dns-webhook/issues/9)) ([4b5a0cb](https://github.com/anexia/external-dns-webhook/commit/4b5a0cbebd245bf2c7e60ec1ad1b6777b165a182))

## [0.1.2](https://github.com/anexia/external-dns-webhook/compare/v0.1.1...v0.1.2) (2024-05-07)


### Bug Fixes

* update Dockerfile entrypoint to use Anexia webhook ([#8](https://github.com/anexia/external-dns-webhook/issues/8)) ([cb23803](https://github.com/anexia/external-dns-webhook/commit/cb23803208d016143aceb9fad4b5c5557d413286))


### Miscellaneous Chores

* update License file ([726710f](https://github.com/anexia/external-dns-webhook/commit/726710f9845919b8e0b775e7d7c045309826893d))

## [0.1.1](https://github.com/anexia/external-dns-webhook/compare/v0.1.0...v0.1.1) (2024-05-06)


### Miscellaneous Chores

* add release-please configuration file ([f317571](https://github.com/anexia/external-dns-webhook/commit/f3175717f5420f5bdf83c2572b102113b0b3f96e))
* Update .conform.yaml and release.yml ([d4e427b](https://github.com/anexia/external-dns-webhook/commit/d4e427b19a318e852ab49cb236fb41bb72820220))
