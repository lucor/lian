# Changelog - lian

## 1.1.0 - 16 August 2025

- feat: add support for replace directives in go.mod files

## 1.0.0 - 11 February 2025

This release primarily updates dependencies to their latest available versions that support Go 1.18.

- update dependencies
  - golang.org/x/mod v0.20.0
  - golang.org/x/sys v0.30.0

## 0.4.1 - 11 February 2024

- fix duplicated `COPYING.MD` should be `COPYING.TXT`

## 0.4.0 - 22 December 2022

- added the `excluded` option (`--excluded`) to exclude specific version of a repository from from the licenses check #1

## 0.3.2 - 21 September 2022

- fix dual license check

## 0.3.1 - 12 February 2022

- action:  missed LICENSE-THIRD-PARTY in release binaries

## 0.3.0 - 12 February 2022

- add license check github action
- add release binary github action
- rename to lian to avoid name conflict with other packages

## 0.2.0 - 08 February 2022

- update to display the license report by default
- the license report is now displayed as a table
- added the `dump` option (`--dump`) to dump all license files

## 0.1.0 - 07 February 2022

- First release
