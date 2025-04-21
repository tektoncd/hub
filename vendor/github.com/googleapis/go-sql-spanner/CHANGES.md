# Changelog

## [1.7.4](https://github.com/googleapis/go-sql-spanner/compare/v1.7.3...v1.7.4) (2024-10-16)


### Bug Fixes

* Data-race with concurrent scan and close ([#304](https://github.com/googleapis/go-sql-spanner/issues/304)) ([cdbea5c](https://github.com/googleapis/go-sql-spanner/commit/cdbea5c87f90c1db40f84a837a1b8ace187dc64f))

## [1.7.3](https://github.com/googleapis/go-sql-spanner/compare/v1.7.2...v1.7.3) (2024-10-08)


### Bug Fixes

* Avoid session leak on transaction retry ([#300](https://github.com/googleapis/go-sql-spanner/issues/300)) ([#301](https://github.com/googleapis/go-sql-spanner/issues/301)) ([129248c](https://github.com/googleapis/go-sql-spanner/commit/129248c4caa86ba41e2d9b27a9a6c9a0533fd302))

## [1.7.2](https://github.com/googleapis/go-sql-spanner/compare/v1.7.1...v1.7.2) (2024-09-30)


### Bug Fixes

* Do not panic when nil is passed in as a query parameter to a pointer of type that implements driver.Valuer via a value receiver ([#296](https://github.com/googleapis/go-sql-spanner/issues/296)) ([816847b](https://github.com/googleapis/go-sql-spanner/commit/816847b09031a272af228350448d774d97ee6d7b))

## [1.7.1](https://github.com/googleapis/go-sql-spanner/compare/v1.7.0...v1.7.1) (2024-08-29)


### Bug Fixes

* Driver.Valuer method was being ignored ([#289](https://github.com/googleapis/go-sql-spanner/issues/289)) ([78fb05d](https://github.com/googleapis/go-sql-spanner/commit/78fb05dedbb13e1fce54accc795d49c1d7425f2e)), refs [#281](https://github.com/googleapis/go-sql-spanner/issues/281)
* Reconnect after all idle connections close ([#290](https://github.com/googleapis/go-sql-spanner/issues/290)) ([b0cdd9d](https://github.com/googleapis/go-sql-spanner/commit/b0cdd9d568b0fa5f55884a5ac701c14f8720dbaa))

## [1.7.0](https://github.com/googleapis/go-sql-spanner/compare/v1.6.0...v1.7.0) (2024-08-23)


### Features

* Support structs in queries ([#282](https://github.com/googleapis/go-sql-spanner/issues/282)) ([21f8b75](https://github.com/googleapis/go-sql-spanner/commit/21f8b75ae3fcfd0de3751fc8a6639d9289a35fe0)), refs [#281](https://github.com/googleapis/go-sql-spanner/issues/281)

## [1.6.0](https://github.com/googleapis/go-sql-spanner/compare/v1.5.0...v1.6.0) (2024-07-14)


### Features

* Add uint64 support ([#272](https://github.com/googleapis/go-sql-spanner/issues/272)) ([5ed831a](https://github.com/googleapis/go-sql-spanner/commit/5ed831a5739753e67679ca237372eebc54de5091))

## [1.5.0](https://github.com/googleapis/go-sql-spanner/compare/v1.4.0...v1.5.0) (2024-07-01)


### Features

* Add -race to testing ([#263](https://github.com/googleapis/go-sql-spanner/issues/263)) ([3af784b](https://github.com/googleapis/go-sql-spanner/commit/3af784b0a786c7079f14a63cbef5da4ad7ad3138))
* Add fuzzing for statement parser ([#250](https://github.com/googleapis/go-sql-spanner/issues/250)) ([2e74813](https://github.com/googleapis/go-sql-spanner/commit/2e7481307b02ec2deb98599e22362ee727d0fc1b))
* Add int support ([#260](https://github.com/googleapis/go-sql-spanner/issues/260)) ([ec15462](https://github.com/googleapis/go-sql-spanner/commit/ec1546238603f409a0aba6627459365c46a92c2b)), refs [#257](https://github.com/googleapis/go-sql-spanner/issues/257)
* **examples/data-types:** Add float32 ([#261](https://github.com/googleapis/go-sql-spanner/issues/261)) ([4d6d1ef](https://github.com/googleapis/go-sql-spanner/commit/4d6d1ef0856cd6b964675cb64a3f4c7335e010dc))


### Bug Fixes

* []*int and []*uint ([#267](https://github.com/googleapis/go-sql-spanner/issues/267)) ([d24b55a](https://github.com/googleapis/go-sql-spanner/commit/d24b55ade2c4c4f5ce01d1f3fcbec41e6ad0d50c))
* **examples:** Make compilable and add a test runner for the samples ([#256](https://github.com/googleapis/go-sql-spanner/issues/256)) ([b9c0b4d](https://github.com/googleapis/go-sql-spanner/commit/b9c0b4d5e5fbfdb3bb358d9a54f66029855a0b2f))
* Prevent potential panic in RemoveStatementHint ([#253](https://github.com/googleapis/go-sql-spanner/issues/253)) ([7d7155b](https://github.com/googleapis/go-sql-spanner/commit/7d7155b2bbc430ed88420d8bf92abbb59638adb4)), refs [#251](https://github.com/googleapis/go-sql-spanner/issues/251)

## [1.4.0](https://github.com/googleapis/go-sql-spanner/compare/v1.3.1...v1.4.0) (2024-05-30)


### Features

* Add exclude_txn_from_change_streams variable ([#236](https://github.com/googleapis/go-sql-spanner/issues/236)) ([ad95d85](https://github.com/googleapis/go-sql-spanner/commit/ad95d8522eb5ca12be089a5cd235a10aa81b0bfd))


### Bug Fixes

* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to 5315273 ([#239](https://github.com/googleapis/go-sql-spanner/issues/239)) ([a39bb8f](https://github.com/googleapis/go-sql-spanner/commit/a39bb8f6e966335531114aadd49f01ddbf072e8a))
* **deps:** Update module cloud.google.com/go to v0.113.0 ([#232](https://github.com/googleapis/go-sql-spanner/issues/232)) ([c2dd030](https://github.com/googleapis/go-sql-spanner/commit/c2dd0303a8b788947e4190659ae87dd67dbf941e))
* **deps:** Update module cloud.google.com/go/longrunning to v0.5.7 ([#228](https://github.com/googleapis/go-sql-spanner/issues/228)) ([cf3cafc](https://github.com/googleapis/go-sql-spanner/commit/cf3cafcd58160076e97bff6e6cc10f72904226d9))
* **deps:** Update module github.com/googleapis/go-sql-spanner to v1.3.1 ([#226](https://github.com/googleapis/go-sql-spanner/issues/226)) ([6ffc02d](https://github.com/googleapis/go-sql-spanner/commit/6ffc02df3b01691100f1a76532cdc109ff807f20))
* **deps:** Update module google.golang.org/api to v0.180.0 ([#227](https://github.com/googleapis/go-sql-spanner/issues/227)) ([926f621](https://github.com/googleapis/go-sql-spanner/commit/926f6214f58f1e60973c000c45be07a232b4649a))
* **deps:** Update module google.golang.org/genproto to v0.0.0-20240509183442-62759503f434 ([#229](https://github.com/googleapis/go-sql-spanner/issues/229)) ([42744bc](https://github.com/googleapis/go-sql-spanner/commit/42744bc8f73d32940629249789cff67a63f4067c))
* **deps:** Update module google.golang.org/genproto/googleapis/rpc to v0.0.0-20240509183442-62759503f434 ([#230](https://github.com/googleapis/go-sql-spanner/issues/230)) ([778aa45](https://github.com/googleapis/go-sql-spanner/commit/778aa45f4a9e1b70dd0cd040e89336c75fa0821a))
* **deps:** Update module google.golang.org/protobuf to v1.34.1 ([#233](https://github.com/googleapis/go-sql-spanner/issues/233)) ([72c5b71](https://github.com/googleapis/go-sql-spanner/commit/72c5b713a749d0f43a03dbad92ca76f583218c55))
* Support reusing named parameters ([#240](https://github.com/googleapis/go-sql-spanner/issues/240)) ([c7140a2](https://github.com/googleapis/go-sql-spanner/commit/c7140a2a4b568d81969949dd76525b67606f04ee))

## [1.3.1](https://github.com/googleapis/go-sql-spanner/compare/v1.3.0...v1.3.1) (2024-04-19)


### Bug Fixes

* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to 6e1732d ([#210](https://github.com/googleapis/go-sql-spanner/issues/210)) ([adeee59](https://github.com/googleapis/go-sql-spanner/commit/adeee596a12d16788d4787524929b8671cc59786))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to 8c6c420 ([#215](https://github.com/googleapis/go-sql-spanner/issues/215)) ([7085e4c](https://github.com/googleapis/go-sql-spanner/commit/7085e4cfea014707a35af6c9612f734fa26f3012))
* **deps:** Update module cloud.google.com/go to v0.112.2 ([#217](https://github.com/googleapis/go-sql-spanner/issues/217)) ([cc65283](https://github.com/googleapis/go-sql-spanner/commit/cc6528320afe1079cd6ebc2999ffe3d3d596fd17))
* **deps:** Update module cloud.google.com/go/longrunning to v0.5.6 ([#211](https://github.com/googleapis/go-sql-spanner/issues/211)) ([8c8a7d8](https://github.com/googleapis/go-sql-spanner/commit/8c8a7d824c5ca19b013787e0aca53d60154750b8))
* **deps:** Update module cloud.google.com/go/spanner to v1.60.0 ([#216](https://github.com/googleapis/go-sql-spanner/issues/216)) ([b22e7a6](https://github.com/googleapis/go-sql-spanner/commit/b22e7a6c4670dad5f481989ed8dbbc5dd5ef3dee))
* **deps:** Update module github.com/googleapis/go-sql-spanner to v1.3.0 ([#212](https://github.com/googleapis/go-sql-spanner/issues/212)) ([8066ff1](https://github.com/googleapis/go-sql-spanner/commit/8066ff13f9981460903fd0f0eea7be7a339a5b43))
* **deps:** Update module google.golang.org/api to v0.174.0 ([#218](https://github.com/googleapis/go-sql-spanner/issues/218)) ([2ef46fc](https://github.com/googleapis/go-sql-spanner/commit/2ef46fce92c691cfe3fc886425331aee57d67143))
* **deps:** Update module google.golang.org/grpc to v1.63.2 ([#221](https://github.com/googleapis/go-sql-spanner/issues/221)) ([69840c4](https://github.com/googleapis/go-sql-spanner/commit/69840c4b2aa226bc7f6c3deac043c5d63f450799))

## [1.3.0](https://github.com/googleapis/go-sql-spanner/compare/v1.2.1...v1.3.0) (2024-03-14)


### Features

* Support float32 data type ([#208](https://github.com/googleapis/go-sql-spanner/issues/208)) ([7db1b8f](https://github.com/googleapis/go-sql-spanner/commit/7db1b8f0be6f60a4d20f206bcebbebedd5f2d460))


### Bug Fixes

* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to c811ad7 ([#194](https://github.com/googleapis/go-sql-spanner/issues/194)) ([e58c5b6](https://github.com/googleapis/go-sql-spanner/commit/e58c5b6888078f61fba743393151f058302ff720))
* **deps:** Update module cloud.google.com/go/spanner to v1.59.0 ([#196](https://github.com/googleapis/go-sql-spanner/issues/196)) ([e03a9b1](https://github.com/googleapis/go-sql-spanner/commit/e03a9b14854560e7f85a8405876d03cbc753b32a))
* **deps:** Update module github.com/googleapis/go-sql-spanner to v1.2.1 ([#197](https://github.com/googleapis/go-sql-spanner/issues/197)) ([6c94089](https://github.com/googleapis/go-sql-spanner/commit/6c940890863a720e1849ab9946c90f4b4b3debd9))
* **deps:** Update module google.golang.org/grpc to v1.62.1 ([#195](https://github.com/googleapis/go-sql-spanner/issues/195)) ([606e7dc](https://github.com/googleapis/go-sql-spanner/commit/606e7dc07f272dd3a0999859ad2af10e6d964189))

## [1.2.1](https://github.com/googleapis/go-sql-spanner/compare/v1.2.0...v1.2.1) (2024-02-05)


### Bug Fixes

* Support uint query params ([#190](https://github.com/googleapis/go-sql-spanner/issues/190)) ([e3e2a9a](https://github.com/googleapis/go-sql-spanner/commit/e3e2a9a9fd14f2a4940dd24e73a5a79aa0431580))


### Documentation

* Add link to pgx sample for PGAdapter ([#191](https://github.com/googleapis/go-sql-spanner/issues/191)) ([e9f8e90](https://github.com/googleapis/go-sql-spanner/commit/e9f8e901b5cc345344d9369b92df8d093b4d47a8))

## [1.2.0](https://github.com/googleapis/go-sql-spanner/compare/v1.1.1...v1.2.0) (2024-02-05)


### Features

* Support of uint for gorm.Model  ([#173](https://github.com/googleapis/go-sql-spanner/issues/173)) ([b56bfa8](https://github.com/googleapis/go-sql-spanner/commit/b56bfa8df549780852674f178a1c06f56f0f4c45))


### Bug Fixes

* **deps:** Update google.golang.org/genproto digest to 1f4bbc5 ([#164](https://github.com/googleapis/go-sql-spanner/issues/164)) ([13ae178](https://github.com/googleapis/go-sql-spanner/commit/13ae17898fc30b3cf8ca24320899ae04ed86d8a7))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to 1f4bbc5 ([#165](https://github.com/googleapis/go-sql-spanner/issues/165)) ([d1b6ab4](https://github.com/googleapis/go-sql-spanner/commit/d1b6ab4737037b7997c3ccaea9ba2cf109948695))
* **deps:** Update module cloud.google.com/go/longrunning to v0.5.5 ([#169](https://github.com/googleapis/go-sql-spanner/issues/169)) ([7bea3f5](https://github.com/googleapis/go-sql-spanner/commit/7bea3f5597b8d7f3dbdec56faf63dc0538874708))
* **deps:** Update module cloud.google.com/go/spanner to v1.56.0 ([#178](https://github.com/googleapis/go-sql-spanner/issues/178)) ([c6f7169](https://github.com/googleapis/go-sql-spanner/commit/c6f7169158dcb471529fbd8d75013e1a86e38519))
* **deps:** Update module github.com/google/uuid to v1.6.0 ([#181](https://github.com/googleapis/go-sql-spanner/issues/181)) ([add2822](https://github.com/googleapis/go-sql-spanner/commit/add282209e18108f9d318a4e477f05fbe509920c))
* **deps:** Update module github.com/googleapis/go-sql-spanner to v1.1.1 ([#162](https://github.com/googleapis/go-sql-spanner/issues/162)) ([aa9c4a0](https://github.com/googleapis/go-sql-spanner/commit/aa9c4a07fc3b15f4f1e5e61cc78257e5dbbdca4c))
* **deps:** Update module google.golang.org/api to v0.161.0 ([#182](https://github.com/googleapis/go-sql-spanner/issues/182)) ([bd33213](https://github.com/googleapis/go-sql-spanner/commit/bd33213dc437c900b1b11f05ef9e4600381f7e0c))


### Documentation

* Clarify postgresql support ([#186](https://github.com/googleapis/go-sql-spanner/issues/186)) ([4aa18f6](https://github.com/googleapis/go-sql-spanner/commit/4aa18f64cabe252f71612294127b9791eb90072c))

## [1.1.1](https://github.com/googleapis/go-sql-spanner/compare/v1.1.0...v1.1.1) (2023-09-29)


### Bug Fixes

* Initalize the connectors map to avoid nil pointer exception ([#158](https://github.com/googleapis/go-sql-spanner/issues/158)) ([ec69c32](https://github.com/googleapis/go-sql-spanner/commit/ec69c321262631c3cdc4369c01ed8abc6932792f))

## [1.1.0](https://github.com/googleapis/go-sql-spanner/compare/v1.0.1...v1.1.0) (2023-07-05)


### Features

* Add missing spanner config properties ([#152](https://github.com/googleapis/go-sql-spanner/issues/152)) ([c6bda23](https://github.com/googleapis/go-sql-spanner/commit/c6bda23e86f3679bb48c33e19ebf413ca984a4ee))
* **driver:** Replace `value.Value` with `Value()` return ([#139](https://github.com/googleapis/go-sql-spanner/issues/139)) ([6f2b96e](https://github.com/googleapis/go-sql-spanner/commit/6f2b96ea14d87a3edfdbaa0738139cd09862e618))

## [1.0.1](https://github.com/googleapis/go-sql-spanner/compare/v1.0.0...v1.0.1) (2023-03-06)


### Bug Fixes

* **deps:** Update module github.com/googleapis/go-sql-spanner to v1 ([#125](https://github.com/googleapis/go-sql-spanner/issues/125)) ([e1ba360](https://github.com/googleapis/go-sql-spanner/commit/e1ba360543b59ae930b4228a03b94cc724dd14d3))

## [1.0.0](https://github.com/googleapis/go-sql-spanner/compare/v1.0.0...v1.0.0) (2022-09-15)


### Features

* Add ARRAY support ([#19](https://github.com/googleapis/go-sql-spanner/issues/19)) ([6b1556a](https://github.com/googleapis/go-sql-spanner/commit/6b1556a8db409fbab0998fdcde59521b26495472))
* Add client side statement parser ([#38](https://github.com/googleapis/go-sql-spanner/issues/38)) ([969bf52](https://github.com/googleapis/go-sql-spanner/commit/969bf52b2cda303349746de9730557242082893c))
* Add support for JSON data type ([#39](https://github.com/googleapis/go-sql-spanner/issues/39)) ([ef52036](https://github.com/googleapis/go-sql-spanner/commit/ef5203657aa8e2173e387ea6aede02d457577790))
* Add support for stale reads ([#44](https://github.com/googleapis/go-sql-spanner/issues/44)) ([2e3a264](https://github.com/googleapis/go-sql-spanner/commit/2e3a2645073d7c9174b3aec934c1e1fcbb06534c))
* Add support of positional parameter in the queries ([#110](https://github.com/googleapis/go-sql-spanner/issues/110)) ([a71a457](https://github.com/googleapis/go-sql-spanner/commit/a71a457261ea8e522d320726a8aeea6768f08acf))
* Allow host in dsn and use statement based transactions ([#10](https://github.com/googleapis/go-sql-spanner/issues/10)) ([0528e13](https://github.com/googleapis/go-sql-spanner/commit/0528e13eed6ccb0b71636554f79c9d278242987c))
* Create standalone samples that run against emulator ([#30](https://github.com/googleapis/go-sql-spanner/issues/30)) ([22b127e](https://github.com/googleapis/go-sql-spanner/commit/22b127e111dc7f8e3a8dfb83e1f3dd736640fcaf))
* Support getting the commit timestamp of a transaction or statement ([#52](https://github.com/googleapis/go-sql-spanner/issues/52)) ([802e7be](https://github.com/googleapis/go-sql-spanner/commit/802e7be6dd18dd6c75991bc129116892a45de944))
* Support mutations ([#43](https://github.com/googleapis/go-sql-spanner/issues/43)) ([2d698b7](https://github.com/googleapis/go-sql-spanner/commit/2d698b754205888fcd4487ec0793c956f42bbf56))


### Bug Fixes

* Add ddl support and change tests to run ddl from driver ([259f98b](https://github.com/googleapis/go-sql-spanner/commit/259f98b017849d158ff799de5c947bb0c39eb4f2))
* Added ddl support to driver, changed tests to call driver rather than api directly ([a9c4c8a](https://github.com/googleapis/go-sql-spanner/commit/a9c4c8a50b3823d2eff368fd672dd7e4adfff1f5))
* Allow users to specify custom credentials ([#57](https://github.com/googleapis/go-sql-spanner/issues/57)) ([1715929](https://github.com/googleapis/go-sql-spanner/commit/171592955606f2b7ff3313d2ab6ddf17ea785f3b))
* Always set a value for dest in Next ([#34](https://github.com/googleapis/go-sql-spanner/issues/34)) ([7b8190c](https://github.com/googleapis/go-sql-spanner/commit/7b8190cbb3d63eac2f4b311208d1da7ea282436b))
* Check named value parameter types ([#35](https://github.com/googleapis/go-sql-spanner/issues/35)) ([f260dd2](https://github.com/googleapis/go-sql-spanner/commit/f260dd247f6a5c69d0d8ec1ac0fb7136f978ae05))
* **deps:** Update all modules ([#108](https://github.com/googleapis/go-sql-spanner/issues/108)) ([2d13f6d](https://github.com/googleapis/go-sql-spanner/commit/2d13f6dcc272d3354a1ebe001fc711e731540aca))
* **deps:** Update google.golang.org/genproto commit hash ([#78](https://github.com/googleapis/go-sql-spanner/issues/78)) ([c9ed2ac](https://github.com/googleapis/go-sql-spanner/commit/c9ed2ac088d9b2bd4cc3bcb613aa0595b1ef9b73))
* **deps:** Update google.golang.org/genproto commit hash to 1739428 ([#81](https://github.com/googleapis/go-sql-spanner/issues/81)) ([3f6ba94](https://github.com/googleapis/go-sql-spanner/commit/3f6ba948cc361edd4190392a629dbef764979ee2))
* **deps:** Update module cloud.google.com/go to v0.100.2 ([#71](https://github.com/googleapis/go-sql-spanner/issues/71)) ([cac55f0](https://github.com/googleapis/go-sql-spanner/commit/cac55f092744104371fa6e539928b2bf73fae1ab))
* **deps:** Update module cloud.google.com/go to v0.102.1 ([#103](https://github.com/googleapis/go-sql-spanner/issues/103)) ([23d315e](https://github.com/googleapis/go-sql-spanner/commit/23d315e644d740a77ff39ac6d1553db81229f2c7))
* **deps:** Update module cloud.google.com/go/spanner to v1.29.0 ([#74](https://github.com/googleapis/go-sql-spanner/issues/74)) ([9a676ba](https://github.com/googleapis/go-sql-spanner/commit/9a676bad33664faf2e6ce937a6c7393407545723))
* **deps:** Update module github.com/google/go-cmp to v0.5.7 ([#80](https://github.com/googleapis/go-sql-spanner/issues/80)) ([cae3a7a](https://github.com/googleapis/go-sql-spanner/commit/cae3a7a7e7ed3d0b9427ebf192a4ad55a7e08728))
* **deps:** Update module github.com/google/uuid to v1.3.0 ([#75](https://github.com/googleapis/go-sql-spanner/issues/75)) ([2072930](https://github.com/googleapis/go-sql-spanner/commit/2072930d8a8171d96a08c7be86578dca0b40b60a))
* **deps:** Update module google.golang.org/api to v0.68.0 ([#76](https://github.com/googleapis/go-sql-spanner/issues/76)) ([8af9417](https://github.com/googleapis/go-sql-spanner/commit/8af94172388cf28a2b8e9dd63e2bd7b40a262a0f))
* **deps:** Update module google.golang.org/grpc to v1.44.0 ([#82](https://github.com/googleapis/go-sql-spanner/issues/82)) ([7b20269](https://github.com/googleapis/go-sql-spanner/commit/7b2026924b9efbdb3ace243596ae978542f2b18b))
* Do not parse hints as parameters ([#45](https://github.com/googleapis/go-sql-spanner/issues/45)) ([56243a5](https://github.com/googleapis/go-sql-spanner/commit/56243a5a1169e86b3b14c02ab5c47a4b950a7f14))
* Pass userAgent in client config ([#118](https://github.com/googleapis/go-sql-spanner/issues/118)) ([2c97068](https://github.com/googleapis/go-sql-spanner/commit/2c97068e002c3c83acb9f181908cf2cbe025e516))
* Race condition when opening multiple connections in parallel as first action ([#59](https://github.com/googleapis/go-sql-spanner/issues/59)) ([0971f81](https://github.com/googleapis/go-sql-spanner/commit/0971f81129f36f519a4c1385a3cd634fa4492c3e))
* Refuse DDL during transactions ([#41](https://github.com/googleapis/go-sql-spanner/issues/41)) ([4e7fa97](https://github.com/googleapis/go-sql-spanner/commit/4e7fa97469752d1605b2af80997a03fd1005b5d9)), refs [#31](https://github.com/googleapis/go-sql-spanner/issues/31)
* Standardize returned errors ([#32](https://github.com/googleapis/go-sql-spanner/issues/32)) ([e780348](https://github.com/googleapis/go-sql-spanner/commit/e7803486f424d579c1b356cf7a500fbb62ac6040)), refs [#14](https://github.com/googleapis/go-sql-spanner/issues/14)
* Use correct type for decoding bytes ([49d08fc](https://github.com/googleapis/go-sql-spanner/commit/49d08fc7ade3559774c88a87b16d13d67c1eef57))


### Documentation

* Add comments to all samples to document what they do and how to use them ([#46](https://github.com/googleapis/go-sql-spanner/issues/46)) ([17a434f](https://github.com/googleapis/go-sql-spanner/commit/17a434f71d6d682ce7974b50f57e6a4193c4f892))
* Add DDL batch sample ([#48](https://github.com/googleapis/go-sql-spanner/issues/48)) ([82a23e4](https://github.com/googleapis/go-sql-spanner/commit/82a23e44db5752d4310133597183ee7967d0efea))
* Add documentation about contributing ([fd70120](https://github.com/googleapis/go-sql-spanner/commit/fd70120b979887d389633a8ffcb8fb647b163cbb))
* Add sample for all data types ([#51](https://github.com/googleapis/go-sql-spanner/issues/51)) ([5a0129b](https://github.com/googleapis/go-sql-spanner/commit/5a0129b45ea0a3b89d900024b99115523de8b8d7))
* Add sample for DML batches ([#49](https://github.com/googleapis/go-sql-spanner/issues/49)) ([bac4a4c](https://github.com/googleapis/go-sql-spanner/commit/bac4a4cef1f628918b925a1dc944ab7e3c732480))
* Add sample for PDML ([#53](https://github.com/googleapis/go-sql-spanner/issues/53)) ([9bd832b](https://github.com/googleapis/go-sql-spanner/commit/9bd832bec3a988c0b322a0563b7300b2c4e09e89))
* Add sample for read-only transaction ([#47](https://github.com/googleapis/go-sql-spanner/issues/47)) ([306c4ea](https://github.com/googleapis/go-sql-spanner/commit/306c4eae4dfef0f38d4c24177206c05092918ab5))
* Cleanup and extend readme ([#60](https://github.com/googleapis/go-sql-spanner/issues/60)) ([2d64f82](https://github.com/googleapis/go-sql-spanner/commit/2d64f827825255743ff4ea631c0b0e8913ef4148))
* Remove disclaimer from README ([#93](https://github.com/googleapis/go-sql-spanner/issues/93)) ([12780e5](https://github.com/googleapis/go-sql-spanner/commit/12780e57be1cfa3df753e92e3dd6c51e06dfb070))


### Miscellaneous Chores

* Release 1.0.0 ([#123](https://github.com/googleapis/go-sql-spanner/issues/123)) ([e7e0d8a](https://github.com/googleapis/go-sql-spanner/commit/e7e0d8a66d7f60cb6a32e28c29ab28ca4d62d5a5))
