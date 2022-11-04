# Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 2020-03-21

### Changed

* **Deprecated** `derr.HasAny` has been renamed to `derr.Is` (to fit with `errors.Is`).
* **Deprecated** `derr.FindFirstMatching` has been renamed to `derr.Find`.
* The `Wrap` and `Wrapf` now wraps using standard Golang wrapping behavior.
* Methods `derr.Wrap`, `derr.Is` (and deprecated `derr.HasAny`), `derr.Find` (and deprecated `derr.FindFirstMatching`) and `derr.ToErrorResponse` now supports Golang based wrapped errors.
* `derr.Is` (and deprecated `derr.HasAny`)

### Added

* Added `derr.DebugErrorChain` that makes it easy to get a debug string of the chain of the error.
* Added `derr.Is` that replaces `derr.HasAny` (see deprecation notice).
* Added `derr.Find` that replaces `derr.FindFirstMatching` (see deprecation notice).

* License changed to Apache 2.0
