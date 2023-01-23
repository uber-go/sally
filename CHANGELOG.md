# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0]
### Added
- Add an optional `description` field to packages.

### Changed
- Use a fluid layout for the index page.
  This renders better on narrow screens.

[1.3.0]: https://github.com/uber-go/sally/compare/v1.2.0...HEAD

## [1.2.0] - 2022-05-17
### Added
- Packages now support specifying branches for target repositories with the
  `branch` field.
- Packages can now override the `url` on a per-package basis with the `url`
  field.

### Changed
- Use documentation badges from https://pkg.go.dev.

Thanks to @lucianonooijen, @jpbede, and @sullivtr for their contributions to
this release.

[1.2.0]: https://github.com/uber-go/sally/compare/v1.1.1...v1.2.0

## [1.1.1] - 2020-03-02
### Fixed
- Fixed godoc badge image.

[1.1.1]: https://github.com/uber-go/sally/compare/v1.1.0...v1.1.1

## [1.1.0] - 2020-02-13
### Added
- Support configuring the godoc server used for documentation links.

### Changed
- Updated default godoc server from `https://godoc.org` to `https://pkg.go.dev`.

[1.1.0]: https://github.com/uber-go/sally/compare/v1.0.1...v1.1.0

## [1.0.1] - 2019-01-03
### Fixed
- Templates are now bundled with the binary rather than requiring a copy of the
  sally source.

[1.0.1]: https://github.com/uber-go/sally/compare/v1.0.0...v1.0.1

## 1.0.0 - 2019-01-03

- Initial tagged release.
