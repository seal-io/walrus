# Changelog
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
