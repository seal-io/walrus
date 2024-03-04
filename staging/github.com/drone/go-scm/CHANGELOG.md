# Changelog

## [Unreleased](https://github.com/drone/go-scm/tree/HEAD)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.34.3...HEAD)

**Closed issues:**

- Cron Jobs don't run with Gitea-scm [\#292](https://github.com/drone/go-scm/issues/292)

**Merged pull requests:**

- feat: change as per new contract of webhook in harness code [\#294](https://github.com/drone/go-scm/pull/294) ([abhinav-harness](https://github.com/abhinav-harness))
- Stash pr commits pagination [\#293](https://github.com/drone/go-scm/pull/293) ([raghavharness](https://github.com/raghavharness))
- Added support for branch names containing '&' and '\#' for GetFile Operations. [\#291](https://github.com/drone/go-scm/pull/291) ([senjucanon2](https://github.com/senjucanon2))

## [v1.34.3](https://github.com/drone/go-scm/tree/v1.34.3) (2023-12-20)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.34.2...v1.34.3)

**Merged pull requests:**

- feat: add pr link as coming from new webhook [\#290](https://github.com/drone/go-scm/pull/290) ([abhinav-harness](https://github.com/abhinav-harness))

## [v1.34.2](https://github.com/drone/go-scm/tree/v1.34.2) (2023-12-20)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.34.1...v1.34.2)

**Merged pull requests:**

- feat: support more events in webhook parse in go-scm for gitness [\#289](https://github.com/drone/go-scm/pull/289) ([abhinav-harness](https://github.com/abhinav-harness))
- fix: ref should be branch name for harness code [\#288](https://github.com/drone/go-scm/pull/288) ([abhinav-harness](https://github.com/abhinav-harness))

## [v1.34.1](https://github.com/drone/go-scm/tree/v1.34.1) (2023-12-08)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.34.0...v1.34.1)

**Fixed bugs:**

- fix: use opts for harness list commits [\#286](https://github.com/drone/go-scm/pull/286) ([abhinav-harness](https://github.com/abhinav-harness))

**Merged pull requests:**

- \(maint\)  v1.34.1 release prep [\#287](https://github.com/drone/go-scm/pull/287) ([tphoney](https://github.com/tphoney))

## [v1.34.0](https://github.com/drone/go-scm/tree/v1.34.0) (2023-12-07)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.33.0...v1.34.0)

**Implemented enhancements:**

- feat: add support for branch update for gitness [\#283](https://github.com/drone/go-scm/pull/283) ([abhinav-harness](https://github.com/abhinav-harness))

**Merged pull requests:**

- \(maint\) v1.34.0 prep [\#284](https://github.com/drone/go-scm/pull/284) ([tphoney](https://github.com/tphoney))

## [v1.33.0](https://github.com/drone/go-scm/tree/v1.33.0) (2023-10-27)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.32.3...v1.33.0)

**Implemented enhancements:**

- feat: Add pr\_comment webhook for harness [\#280](https://github.com/drone/go-scm/pull/280) ([abhinav-harness](https://github.com/abhinav-harness))

**Merged pull requests:**

- \(maint\) prep 1.33.0 [\#281](https://github.com/drone/go-scm/pull/281) ([tphoney](https://github.com/tphoney))

## [v1.32.3](https://github.com/drone/go-scm/tree/v1.32.3) (2023-10-11)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.32.2...v1.32.3)

**Fixed bugs:**

- fix: ref should have pullreq instead of pull for gitness [\#279](https://github.com/drone/go-scm/pull/279) ([abhinav-harness](https://github.com/abhinav-harness))
- fix: ref should have pullreq instead of pull for gitness [\#278](https://github.com/drone/go-scm/pull/278) ([abhinav-harness](https://github.com/abhinav-harness))

## [v1.32.2](https://github.com/drone/go-scm/tree/v1.32.2) (2023-10-03)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.32.1...v1.32.2)

**Implemented enhancements:**

- feat: Harness list commits api update as per new spec [\#277](https://github.com/drone/go-scm/pull/277) ([abhinav-harness](https://github.com/abhinav-harness))

## [v1.32.1](https://github.com/drone/go-scm/tree/v1.32.1) (2023-09-27)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.32.0...v1.32.1)

**Fixed bugs:**

- fix: Gitness get content missing query param [\#275](https://github.com/drone/go-scm/pull/275) ([abhinav-harness](https://github.com/abhinav-harness))

**Merged pull requests:**

- \(maint\) prep for 1.32.1 [\#276](https://github.com/drone/go-scm/pull/276) ([tphoney](https://github.com/tphoney))
- \(maint\) clean integration testing for stash [\#273](https://github.com/drone/go-scm/pull/273) ([tphoney](https://github.com/tphoney))

## [v1.32.0](https://github.com/drone/go-scm/tree/v1.32.0) (2023-09-12)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.31.2...v1.32.0)

**Implemented enhancements:**

- \[feat\]: \[CDS-75848\]: Add new action type for github provider [\#270](https://github.com/drone/go-scm/pull/270) ([rathodmeetsatish](https://github.com/rathodmeetsatish))

**Merged pull requests:**

- \(maint\) release prep for 1.32.0 [\#272](https://github.com/drone/go-scm/pull/272) ([tphoney](https://github.com/tphoney))

## [v1.31.2](https://github.com/drone/go-scm/tree/v1.31.2) (2023-08-31)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.31.1...v1.31.2)

**Fixed bugs:**

- fix: \[CODE-727\]: change branch in source and target for harness provider [\#264](https://github.com/drone/go-scm/pull/264) ([abhinav-harness](https://github.com/abhinav-harness))

## [v1.31.1](https://github.com/drone/go-scm/tree/v1.31.1) (2023-08-29)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.31.0...v1.31.1)

**Fixed bugs:**

- Fix diff api response conversion for harness compareChange [\#269](https://github.com/drone/go-scm/pull/269) ([shubham149](https://github.com/shubham149))
- Fix api name for fetching diff in harness driver [\#268](https://github.com/drone/go-scm/pull/268) ([shubham149](https://github.com/shubham149))
- Fix compare change api result for harness [\#267](https://github.com/drone/go-scm/pull/267) ([shubham149](https://github.com/shubham149))

## [v1.31.0](https://github.com/drone/go-scm/tree/v1.31.0) (2023-08-15)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.30.0...v1.31.0)

**Implemented enhancements:**

- \[IAC-941\]: PR comment creation for BitBucket [\#265](https://github.com/drone/go-scm/pull/265) ([scottyw-harness](https://github.com/scottyw-harness))
- Implemented FindMembership method in organization service for gitea driver [\#263](https://github.com/drone/go-scm/pull/263) ([cod3rboy](https://github.com/cod3rboy))

**Closed issues:**

- \(missing feature\) add support to check organization membership in gitea driver [\#262](https://github.com/drone/go-scm/issues/262)

**Merged pull requests:**

- \(maint\) v1.31.0 release prep [\#266](https://github.com/drone/go-scm/pull/266) ([tphoney](https://github.com/tphoney))

## [v1.30.0](https://github.com/drone/go-scm/tree/v1.30.0) (2023-07-19)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.29.1...v1.30.0)

**Implemented enhancements:**

- \[feat\]: \[CDS-73572\]: Support List Repo Live Search for all git providers [\#261](https://github.com/drone/go-scm/pull/261) ([adivishy1](https://github.com/adivishy1))
- \[feat\]: \[CDS-73030\]: Support for text based branch filtration [\#260](https://github.com/drone/go-scm/pull/260) ([adivishy1](https://github.com/adivishy1))
- feat: \[CDS-69341\]: add find user email api for github in go-scm [\#256](https://github.com/drone/go-scm/pull/256) ([shalini-agr](https://github.com/shalini-agr))

**Fixed bugs:**

- fix: \[CDS-67745\]: fix find user email api for bitbucket in go-scm [\#255](https://github.com/drone/go-scm/pull/255) ([shalini-agr](https://github.com/shalini-agr))
- fix: \[CI-6978\] fixed gitlab webhook parse [\#253](https://github.com/drone/go-scm/pull/253) ([devkimittal](https://github.com/devkimittal))
- Add required header for bitbucket server in commit API use-case to handle csrf failures [\#252](https://github.com/drone/go-scm/pull/252) ([mohitg0795](https://github.com/mohitg0795))

**Merged pull requests:**

- \(maint\) stash/bitbucket on prem v5 add push webhook test [\#257](https://github.com/drone/go-scm/pull/257) ([tphoney](https://github.com/tphoney))

## [v1.29.1](https://github.com/drone/go-scm/tree/v1.29.1) (2023-02-16)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.29.0...v1.29.1)

**Fixed bugs:**

- \(fix\) - azure content list queryparam incorrect [\#249](https://github.com/drone/go-scm/pull/249) ([eoinmcafee00](https://github.com/eoinmcafee00))

## [v1.29.0](https://github.com/drone/go-scm/tree/v1.29.0) (2023-02-15)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.28.1...v1.29.0)

**Implemented enhancements:**

- \(feat\) harness, add finduser [\#250](https://github.com/drone/go-scm/pull/250) ([tphoney](https://github.com/tphoney))
- \(feat\) harness, fix create branch, PR calls [\#247](https://github.com/drone/go-scm/pull/247) ([tphoney](https://github.com/tphoney))
- \(feat\) harness, add user and compare branches [\#246](https://github.com/drone/go-scm/pull/246) ([tphoney](https://github.com/tphoney))
- \(feat\) harness, add list commits / branches [\#245](https://github.com/drone/go-scm/pull/245) ([tphoney](https://github.com/tphoney))
- \(feat\) harness, add webhook parsing [\#244](https://github.com/drone/go-scm/pull/244) ([tphoney](https://github.com/tphoney))
- fetch branch for bitbucket onprem [\#242](https://github.com/drone/go-scm/pull/242) ([devkimittal](https://github.com/devkimittal))
- \(feat\) harness, add repo list [\#241](https://github.com/drone/go-scm/pull/241) ([tphoney](https://github.com/tphoney))
- Harness move [\#237](https://github.com/drone/go-scm/pull/237) ([tphoney](https://github.com/tphoney))

**Fixed bugs:**

- \(fix\) harness, webhook fixes [\#248](https://github.com/drone/go-scm/pull/248) ([tphoney](https://github.com/tphoney))
- fix: \[PIE-7927\]: Fix header value typo issue for BB OnPrem CSRF header [\#236](https://github.com/drone/go-scm/pull/236) ([mohitg0795](https://github.com/mohitg0795))

**Merged pull requests:**

- \(maint\) prep for 1.29.0 [\#251](https://github.com/drone/go-scm/pull/251) ([tphoney](https://github.com/tphoney))

## [v1.28.1](https://github.com/drone/go-scm/tree/v1.28.1) (2023-01-27)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.28.0...v1.28.1)

**Fixed bugs:**

- feat: \[PIE-7927\]: added header to avoid/bypass csrf check [\#234](https://github.com/drone/go-scm/pull/234) ([mohitg0795](https://github.com/mohitg0795))

**Closed issues:**

- Gogs commit fails to deserialize commitDetails in some cases [\#231](https://github.com/drone/go-scm/issues/231)

**Merged pull requests:**

- \(maint\) prep 1.28.1 release [\#235](https://github.com/drone/go-scm/pull/235) ([tphoney](https://github.com/tphoney))

## [v1.28.0](https://github.com/drone/go-scm/tree/v1.28.0) (2022-11-22)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.27.0...v1.28.0)

**Implemented enhancements:**

- Add Actor UUID to push and branch create events for Bitbucket [\#230](https://github.com/drone/go-scm/pull/230) ([rutvijmehta-harness](https://github.com/rutvijmehta-harness))
- Add support for Github release webhook [\#229](https://github.com/drone/go-scm/pull/229) ([vcalasansh](https://github.com/vcalasansh))
- Add Actor UUID to Sender for all webhooks responses for Bitbucket [\#227](https://github.com/drone/go-scm/pull/227) ([rutvijmehta-harness](https://github.com/rutvijmehta-harness))
- added date info for commits in push hook [\#223](https://github.com/drone/go-scm/pull/223) ([raghavharness](https://github.com/raghavharness))
- Added support for branch in list commits bb onprem API [\#215](https://github.com/drone/go-scm/pull/215) ([mohitg0795](https://github.com/mohitg0795))
- \[PL-26239\]: added api to list installation for github app [\#213](https://github.com/drone/go-scm/pull/213) ([bhavya181](https://github.com/bhavya181))

**Fixed bugs:**

- fixbug: gitee convert repository [\#226](https://github.com/drone/go-scm/pull/226) ([kit101](https://github.com/kit101))
- Bitbucket sha fix for merged pr [\#225](https://github.com/drone/go-scm/pull/225) ([raghavharness](https://github.com/raghavharness))
- decoding projectName for azure repo [\#224](https://github.com/drone/go-scm/pull/224) ([raghavharness](https://github.com/raghavharness))
- added omitempty annotation for secret [\#221](https://github.com/drone/go-scm/pull/221) ([raghavharness](https://github.com/raghavharness))
- \[PL-26239\]: fix for list response [\#218](https://github.com/drone/go-scm/pull/218) ([bhavya181](https://github.com/bhavya181))

**Closed issues:**

- gitlab: force\_remove\_source\_branch type is inconsistent [\#228](https://github.com/drone/go-scm/issues/228)
- gitee: When the name and path are inconsistent, got 404 error [\#217](https://github.com/drone/go-scm/issues/217)
- file naming conventions [\#208](https://github.com/drone/go-scm/issues/208)
- Support for Azure Devops git repos? [\#53](https://github.com/drone/go-scm/issues/53)

**Merged pull requests:**

- \(maint\) release prep for 1.28.0 [\#232](https://github.com/drone/go-scm/pull/232) ([tphoney](https://github.com/tphoney))
- \(maint\) fixing naming and add more go best practice [\#211](https://github.com/drone/go-scm/pull/211) ([tphoney](https://github.com/tphoney))

## [v1.27.0](https://github.com/drone/go-scm/tree/v1.27.0) (2022-07-19)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.26.0...v1.27.0)

**Merged pull requests:**

- Update scm version 1.27.0 [\#206](https://github.com/drone/go-scm/pull/206) ([raghavharness](https://github.com/raghavharness))
- Using resource version 2.0 for Azure [\#205](https://github.com/drone/go-scm/pull/205) ([raghavharness](https://github.com/raghavharness))

## [v1.26.0](https://github.com/drone/go-scm/tree/v1.26.0) (2022-07-01)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.25.0...v1.26.0)

**Implemented enhancements:**

- Support parsing PR comment events for Bitbucket Cloud [\#202](https://github.com/drone/go-scm/pull/202) ([rutvijmehta-harness](https://github.com/rutvijmehta-harness))
- added issue comment hook support for Azure [\#200](https://github.com/drone/go-scm/pull/200) ([raghavharness](https://github.com/raghavharness))

**Fixed bugs:**

- \[CI-4623\] - Azure webhook parseAPI changes [\#198](https://github.com/drone/go-scm/pull/198) ([raghavharness](https://github.com/raghavharness))

**Merged pull requests:**

- Fixed formatting in README.md [\#199](https://github.com/drone/go-scm/pull/199) ([hemanthmantri](https://github.com/hemanthmantri))

## [v1.25.0](https://github.com/drone/go-scm/tree/v1.25.0) (2022-06-16)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.24.0...v1.25.0)

**Implemented enhancements:**

- Support parsing Gitlab Note Hook event [\#194](https://github.com/drone/go-scm/pull/194) ([rutvijmehta-harness](https://github.com/rutvijmehta-harness))

**Fixed bugs:**

- \[PL-25889\]: fix list branches Azure API [\#195](https://github.com/drone/go-scm/pull/195) ([bhavya181](https://github.com/bhavya181))
- Return project specific hooks only in ListHooks API for Azure. [\#192](https://github.com/drone/go-scm/pull/192) ([raghavharness](https://github.com/raghavharness))

**Merged pull requests:**

- Update scm version 1.25.0 [\#197](https://github.com/drone/go-scm/pull/197) ([rutvijmehta-harness](https://github.com/rutvijmehta-harness))

## [v1.24.0](https://github.com/drone/go-scm/tree/v1.24.0) (2022-06-07)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.23.0...v1.24.0)

**Implemented enhancements:**

- Added PR find and listCommit API support for Azure [\#188](https://github.com/drone/go-scm/pull/188) ([raghavharness](https://github.com/raghavharness))

**Fixed bugs:**

- remove redundant slash from list commits api [\#190](https://github.com/drone/go-scm/pull/190) ([aman-harness](https://github.com/aman-harness))
- Using target commit instead of source in base info for azure [\#189](https://github.com/drone/go-scm/pull/189) ([raghavharness](https://github.com/raghavharness))

**Closed issues:**

- gitee client pagination bug [\#187](https://github.com/drone/go-scm/issues/187)

**Merged pull requests:**

- release\_prep\_v1.24.0 [\#191](https://github.com/drone/go-scm/pull/191) ([tphoney](https://github.com/tphoney))

## [v1.23.0](https://github.com/drone/go-scm/tree/v1.23.0) (2022-05-23)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.22.0...v1.23.0)

**Implemented enhancements:**

- Add the support to fetch commit of a particular file [\#182](https://github.com/drone/go-scm/pull/182) ([DeepakPatankar](https://github.com/DeepakPatankar))

**Fixed bugs:**

- Remove the null value de-reference issue when the bitbucket server url is nil [\#183](https://github.com/drone/go-scm/pull/183) ([DeepakPatankar](https://github.com/DeepakPatankar))
- \[PL-24913\]: Handle the error raised while creating a multipart input [\#181](https://github.com/drone/go-scm/pull/181) ([DeepakPatankar](https://github.com/DeepakPatankar))

**Merged pull requests:**

- Upgrade the scm version [\#185](https://github.com/drone/go-scm/pull/185) ([DeepakPatankar](https://github.com/DeepakPatankar))

## [v1.22.0](https://github.com/drone/go-scm/tree/v1.22.0) (2022-05-10)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.21.1...v1.22.0)

**Implemented enhancements:**

- \[feat\]: \[PL-24913\]: Add the support for create and update file in bitbucket server [\#177](https://github.com/drone/go-scm/pull/177) ([DeepakPatankar](https://github.com/DeepakPatankar))
- \[feat\]: \[PL-24915\]: Add the support for create branches in bitbucket server [\#174](https://github.com/drone/go-scm/pull/174) ([DeepakPatankar](https://github.com/DeepakPatankar))
- \[PL-24911\]: Make project name as optional param in Azure Repo APIs [\#173](https://github.com/drone/go-scm/pull/173) ([DeepakPatankar](https://github.com/DeepakPatankar))

**Fixed bugs:**

- \[feat\]: \[PL-25025\]: Updated Project validation for Azure API [\#179](https://github.com/drone/go-scm/pull/179) ([mankrit-singh](https://github.com/mankrit-singh))
- \[fix\]: \[PL-24880\]: Trim the ref when fetching default branch in get Repo API [\#172](https://github.com/drone/go-scm/pull/172) ([DeepakPatankar](https://github.com/DeepakPatankar))
- fixbug: gitee populatePageValues [\#167](https://github.com/drone/go-scm/pull/167) ([kit101](https://github.com/kit101))

**Closed issues:**

- gitea find commit   [\#125](https://github.com/drone/go-scm/issues/125)

**Merged pull requests:**

- \[feat\]: \[PL-25025\]: Changelog Updated/New Version [\#180](https://github.com/drone/go-scm/pull/180) ([mankrit-singh](https://github.com/mankrit-singh))

## [v1.21.1](https://github.com/drone/go-scm/tree/v1.21.1) (2022-04-22)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.21.0...v1.21.1)

**Fixed bugs:**

- remove double invocation of convertRepository [\#170](https://github.com/drone/go-scm/pull/170) ([d1wilko](https://github.com/d1wilko))

**Merged pull requests:**

- \(maint\) release prep for 1.21.1 [\#171](https://github.com/drone/go-scm/pull/171) ([d1wilko](https://github.com/d1wilko))

## [v1.21.0](https://github.com/drone/go-scm/tree/v1.21.0) (2022-04-22)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.20.0...v1.21.0)

**Implemented enhancements:**

- Add support for repository find in azure [\#164](https://github.com/drone/go-scm/pull/164) ([goelsatyam2](https://github.com/goelsatyam2))
- \(feat\) add azure webhook parsing, creation deletion & list [\#163](https://github.com/drone/go-scm/pull/163) ([tphoney](https://github.com/tphoney))
- \(DRON-242\) azure add compare commits,get commit,list repos [\#162](https://github.com/drone/go-scm/pull/162) ([tphoney](https://github.com/tphoney))

**Fixed bugs:**

- \(fix\) handle nil repos in github responses [\#168](https://github.com/drone/go-scm/pull/168) ([tphoney](https://github.com/tphoney))

**Closed issues:**

- When attempting to clone my git repo from GitLab drone hangs on git fetch. [\#161](https://github.com/drone/go-scm/issues/161)
- Fix dump response [\#119](https://github.com/drone/go-scm/issues/119)

**Merged pull requests:**

- \(maint\) release prep for 1.21.0 [\#169](https://github.com/drone/go-scm/pull/169) ([d1wilko](https://github.com/d1wilko))

## [v1.20.0](https://github.com/drone/go-scm/tree/v1.20.0) (2022-03-08)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.19.1...v1.20.0)

**Implemented enhancements:**

- \(DRON-242\) initial implementation of azure devops support [\#158](https://github.com/drone/go-scm/pull/158) ([tphoney](https://github.com/tphoney))

**Fixed bugs:**

- fixed raw response dumping in client [\#159](https://github.com/drone/go-scm/pull/159) ([marko-gacesa](https://github.com/marko-gacesa))

**Merged pull requests:**

- \(maint\) release prep for 1.20.0 [\#160](https://github.com/drone/go-scm/pull/160) ([d1wilko](https://github.com/d1wilko))

## [v1.19.1](https://github.com/drone/go-scm/tree/v1.19.1) (2022-02-23)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.19.0...v1.19.1)

**Fixed bugs:**

- Bitbucket list files fix [\#154](https://github.com/drone/go-scm/pull/154) ([mohitg0795](https://github.com/mohitg0795))
- GitHub list commits fix [\#152](https://github.com/drone/go-scm/pull/152) ([mohitg0795](https://github.com/mohitg0795))
- Bitbucket compare changes fix for rename and removed file ops [\#151](https://github.com/drone/go-scm/pull/151) ([mohitg0795](https://github.com/mohitg0795))

**Merged pull requests:**

- prep for v1.19.1 [\#155](https://github.com/drone/go-scm/pull/155) ([tphoney](https://github.com/tphoney))

## [v1.19.0](https://github.com/drone/go-scm/tree/v1.19.0) (2022-02-09)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.18.0...v1.19.0)

**Implemented enhancements:**

- \(feat\) add path support for list commits on github and gitlab [\#149](https://github.com/drone/go-scm/pull/149) ([tphoney](https://github.com/tphoney))
- Extending bitbucket listCommits API to fetch commits for a given file [\#148](https://github.com/drone/go-scm/pull/148) ([mohitg0795](https://github.com/mohitg0795))
- Update GitHub signature header to use sha256 [\#123](https://github.com/drone/go-scm/pull/123) ([nlecoy](https://github.com/nlecoy))

**Merged pull requests:**

- v1.19.0 release prep [\#150](https://github.com/drone/go-scm/pull/150) ([tphoney](https://github.com/tphoney))

## [v1.18.0](https://github.com/drone/go-scm/tree/v1.18.0) (2022-01-18)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.17.0...v1.18.0)

**Implemented enhancements:**

- Added support for parsing prevFilePath field from github compare commits API response [\#143](https://github.com/drone/go-scm/pull/143) ([mohitg0795](https://github.com/mohitg0795))

**Fixed bugs:**

- Implement parsing/handling for missing pull request webhook events for BitBucket Server \(Stash\) driver [\#130](https://github.com/drone/go-scm/pull/130) ([raphendyr](https://github.com/raphendyr))

**Closed issues:**

- Bitbucket Stash driver doesn't handle event `pr:from_ref_updated` \(new commits / force push\) [\#116](https://github.com/drone/go-scm/issues/116)

**Merged pull requests:**

- release prep v1.18.0 [\#147](https://github.com/drone/go-scm/pull/147) ([eoinmcafee00](https://github.com/eoinmcafee00))

## [v1.17.0](https://github.com/drone/go-scm/tree/v1.17.0) (2022-01-07)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.16.3...v1.17.0)

**Implemented enhancements:**

- \(feat\) map archive flag to repo response [\#141](https://github.com/drone/go-scm/pull/141) ([eoinmcafee00](https://github.com/eoinmcafee00))
- Add the support for delete of the bitbucket file [\#139](https://github.com/drone/go-scm/pull/139) ([DeepakPatankar](https://github.com/DeepakPatankar))

**Fixed bugs:**

- Fix the syntax error of the example code [\#135](https://github.com/drone/go-scm/pull/135) ([LinuxSuRen](https://github.com/LinuxSuRen))

**Closed issues:**

- The deprecation of Bitbucket API endpoint /2.0/teams breaks user registration [\#136](https://github.com/drone/go-scm/issues/136)

**Merged pull requests:**

- release prep for v1.17.0 [\#142](https://github.com/drone/go-scm/pull/142) ([eoinmcafee00](https://github.com/eoinmcafee00))

## [v1.16.3](https://github.com/drone/go-scm/tree/v1.16.3) (2021-12-30)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.16.2...v1.16.3)

**Fixed bugs:**

- fix the deprecation of Bitbucket API endpoint /2.0/teams breaks user registration \(136\) [\#137](https://github.com/drone/go-scm/pull/137) ([eoinmcafee00](https://github.com/eoinmcafee00))

**Closed issues:**

- Any plans to support manage wehook [\#134](https://github.com/drone/go-scm/issues/134)

**Merged pull requests:**

- V1.16.3 [\#138](https://github.com/drone/go-scm/pull/138) ([eoinmcafee00](https://github.com/eoinmcafee00))

## [v1.16.2](https://github.com/drone/go-scm/tree/v1.16.2) (2021-11-30)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.16.1...v1.16.2)

**Merged pull requests:**

- release prep v1.16.2 [\#132](https://github.com/drone/go-scm/pull/132) ([marko-gacesa](https://github.com/marko-gacesa))
- fixbug: gitee webhook parse [\#131](https://github.com/drone/go-scm/pull/131) ([kit101](https://github.com/kit101))

## [v1.16.1](https://github.com/drone/go-scm/tree/v1.16.1) (2021-11-19)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.16.0...v1.16.1)

**Fixed bugs:**

- swap repo and target in bitbucket CompareChanges [\#127](https://github.com/drone/go-scm/pull/127) ([jimsheldon](https://github.com/jimsheldon))

**Merged pull requests:**

- release prep v1.16.1 [\#129](https://github.com/drone/go-scm/pull/129) ([eoinmcafee00](https://github.com/eoinmcafee00))

## [v1.16.0](https://github.com/drone/go-scm/tree/v1.16.0) (2021-11-19)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.15.2...v1.16.0)

**Implemented enhancements:**

- release prep for 1.16.0 [\#128](https://github.com/drone/go-scm/pull/128) ([eoinmcafee00](https://github.com/eoinmcafee00))
- Feat: implemented gitee provider [\#124](https://github.com/drone/go-scm/pull/124) ([kit101](https://github.com/kit101))
- add release & milestone functionality [\#121](https://github.com/drone/go-scm/pull/121) ([eoinmcafee00](https://github.com/eoinmcafee00))

**Fixed bugs:**

- Fix Gitea example code on README.md [\#126](https://github.com/drone/go-scm/pull/126) ([lunny](https://github.com/lunny))

## [v1.15.2](https://github.com/drone/go-scm/tree/v1.15.2) (2021-07-20)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.15.1...v1.15.2)

**Fixed bugs:**

- Fixing Gitea commit API in case `ref/heads/` prefix is added to ref [\#108](https://github.com/drone/go-scm/pull/108) ([Vici37](https://github.com/Vici37))
- use access json header / extend error message parsing for stash [\#89](https://github.com/drone/go-scm/pull/89) ([bakito](https://github.com/bakito))

**Closed issues:**

- Drone and Bitbucket broken for write permission detection for drone build restart permission. [\#87](https://github.com/drone/go-scm/issues/87)

**Merged pull requests:**

- \(maint\) prep for v.1.15.2 release [\#118](https://github.com/drone/go-scm/pull/118) ([tphoney](https://github.com/tphoney))
- Add a vet step to drone config [\#83](https://github.com/drone/go-scm/pull/83) ([tboerger](https://github.com/tboerger))

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.15.1
### Added
- (DRON-88) github fix pr ListChanges deleted/renamed status, from [@tphoney](https://github.com/tphoney). See [#113](https://github.com/drone/go-scm/pull/113).
- (DRON-84) github fix pr write permission issue with bitbucket server, from [@eoinmcafee0](https://github.com/eoinmcafee00). See [#114](https://github.com/drone/go-scm/pull/114).

## 1.15.0
### Added
- add delete file for github and gitlab, from [@tphoney](https://github.com/tphoney). See [#110](https://github.com/drone/go-scm/pull/110).

## 1.14.1
### Fixed
- fix gitlab repo encoding in commits from pr request, from [@aradisavljevic](https://github.com/aradisavljevic). See [#109](https://github.com/drone/go-scm/pull/109).

## 1.14.0
### Added
- Added ListCommits in pull request api, from [@aradisavljevic](https://github.com/aradisavljevic). See [#106](https://github.com/drone/go-scm/pull/106).

## 1.13.1
### Fixed
- gitlab, content.find return last_commit_id not commit_id, from [@tphoney](https://github.com/tphoney). See [#104](https://github.com/drone/go-scm/pull/104).

## 1.13.0
### Added
- Create branch functionality, from [@tphoney](https://github.com/tphoney). See [#103](https://github.com/drone/go-scm/pull/103).

## 1.12.0
### Added
- return sha/blob_id for content.list, from [@tphoney](https://github.com/tphoney). See [#102](https://github.com/drone/go-scm/pull/102).

## 1.11.0
### Added
- normalise sha in content, add bitbucket create/update, from [@tphoney](https://github.com/tphoney). See [#101](https://github.com/drone/go-scm/pull/101).

## 1.10.0
### Added
- return hash/object_id for files changed in github, from [@tphoney](https://github.com/tphoney). See [#99](https://github.com/drone/go-scm/pull/99).

## 1.9.0
### Added
- Added issue_comment parsing for github webhook, from [@aman-harness](https://github.com/aman-harness). See [#91](https://github.com/drone/go-scm/pull/91).
- Added Pr in issue, from [@aman-harness](https://github.com/aman-harness). See [#93](https://github.com/drone/go-scm/pull/93).
- gitlab contents. Find returns hash/blob, from [@tphoney](https://github.com/tphoney). See [#97](https://github.com/drone/go-scm/pull/97).
- add ListCommits for gitea and stash, from [@tphoney](https://github.com/tphoney). See [#98](https://github.com/drone/go-scm/pull/98).

### Changed
- retry with event subset for legacy stash versions, from [@bakito](https://github.com/bakito). See [#90](https://github.com/drone/go-scm/pull/90).

## 1.8.0
### Added
- Support for GitLab visibility attribute, from [@bradrydzewski](https://github.com/bradrydzewski). See [79951ad](https://github.com/drone/go-scm/commit/79951ad7a0d0b1989ea84d99be31fcb9320ae348).
- Support for GitHub visibility attribute, from [@bradrydzewski](https://github.com/bradrydzewski). See [5141b8e](https://github.com/drone/go-scm/commit/5141b8e1db921fe2101c12594c5159b9ffffebc3).

### Changed
- Support for parsing unknown pull request events, from [@bradrydzewski](https://github.com/bradrydzewski). See [ffa46d9](https://github.com/drone/go-scm/commit/ffa46d955454baa609975eebbe9fdfc4b0a9f7e9).

## 1.7.2
### Added
- Support for finding and listing repository tags in GitHub driver, from [@chhsia0](https://github.com/chhsia0). See [#79](https://github.com/drone/go-scm/pull/79).
- Support for finding and listing repository tags in Gitea driver, from [@bradyrdzewski](https://github.com/bradyrdzewski). See [427b8a8](https://github.com/drone/go-scm/commit/427b8a85897c892148801824760bc66d3a3cdcdb).
- Support for git object hashes in GitHub, from from [@bradyrdzewski](https://github.com/bradyrdzewski). See [5230330](https://github.com/drone/go-scm/commit/523033025a7ee875fcfb156f4c660b37e269b1a8).
- Support for before and after commit sha in Stash driver, from [@jlehtimaki](https://github.com/jlehtimaki). See [#82](https://github.com/drone/go-scm/pull/82).
- Support for before and after commit sha in GitLab and Bitbucket driver, from [@shubhag](https://github.com/shubhag). See [#85](https://github.com/drone/go-scm/pull/85).

## 1.7.1
### Added
- Support for skip verification in Bitbucket webhook creation, from [@chhsia0](https://github.com/chhsia0). See [#63](https://github.com/drone/go-scm/pull/63).
- Support for Gitea pagination, from [@CirnoT](https://github.com/CirnoT). See [#66](https://github.com/drone/go-scm/pull/66).
- Support for labels in pull request resources, from [@takirala](https://github.com/takirala). See [#67](https://github.com/drone/go-scm/pull/67).
- Support for updating webhooks, from [@chhsia0](https://github.com/chhsia0). See [#71](https://github.com/drone/go-scm/pull/71).

### Fixed
- Populate diff links in pull request resources, from [@shubhag](https://github.com/shubhag). See [#75](https://github.com/drone/go-scm/pull/75).
- Filter Bitbucket repository search by project, from [@bradrydzewski](https://github.com/bradrydzewski).

## 1.7.0
### Added
- Improve status display text in new bitbucket pull request screen, from [@bradrydzewski](https://github.com/bradrydzewski). See [#27](https://github.com/drone/go-scm/issues/27).
- Implement timestamp value for GitHub push webhooks, from [@bradrydzewski](https://github.com/bradrydzewski).
- Implement deep link to branch.
- Implement git compare function to compare two separate commits, from [@chhsia0](https://github.com/chhsia0).
- Implement support for creating and updating GitLab and GitHub repository contents, from [@zhuxiaoyang](https://github.com/zhuxiaoyang).
- Capture Repository link for Gitea, Gogs and Gitlab, from [@chhsia0](https://github.com/chhsia0).

### Fixed
- Fix issue with GitHub enterprise deep link including API prefix, from [@bradrydzewski](https://github.com/bradrydzewski).
- Fix issue with GitHub deploy hooks for commits having an invalid reference, from [@bradrydzewski](https://github.com/bradrydzewski).
- Support for Skipping SSL verification for GitLab webhooks. See [#40](https://github.com/drone/go-scm/pull/40).
- Support for Skipping SSL verification for GitHub webhooks. See [#44](https://github.com/drone/go-scm/pull/40).
- Fix issue with handling slashes in Bitbucket branch names. See [#7](https://github.com/drone/go-scm/pull/47).
- Fix incorrect Gitea tag link. See [#52](https://github.com/drone/go-scm/pull/52).
- Encode ref when making Gitea API calls. See [#61](https://github.com/drone/go-scm/pull/61).

## [1.6.0]
### Added
- Support Head and Base sha for GitHub pull request, from [@bradrydzewski](https://github.com/bradrydzewski).
- Support Before sha for Bitbucket, from [@jkdev81](https://github.com/jkdev81).
- Support for creating GitHub deployment hooks, from [@bradrydzewski](https://github.com/bradrydzewski).
- Endpoint to get organization membership for GitHub, from [@bradrydzewski](https://github.com/bradrydzewski).
- Functions to generate deep links to git resources, from [@bradrydzewski](https://github.com/bradrydzewski).

### Fixed
- Fix issue getting a GitLab commit by ref, from [@bradrydzewski](https://github.com/bradrydzewski).

## [1.5.0]
### Added

- Fix missing sha for Gitea tag hooks, from [@techknowlogick](https://github.com/techknowlogick). See [#22](https://github.com/drone/go-scm/pull/22).
- Support for Gitea webhook signature verification, from [@techknowlogick](https://github.com/techknowlogick).

## [1.4.0]
### Added

- Fix issues base64 decoding GitLab content, from [@bradrydzewski](https://github.com/bradrydzewski).

## [1.3.0]
### Added

- Fix missing avatar in Gitea commit from [@jgeek1011](https://github.com/geek1011).
- Implement GET commit endpoint for Gogs from [@ogarcia](https://github.com/ogarcia).


\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
