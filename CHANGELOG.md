# Changelog

## [1.0.3](https://github.com/elct9620/pdf64/compare/v1.0.2...v1.0.3) (2025-03-08)


### Miscellaneous Chores

* release 1.0.3 ([25b81ff](https://github.com/elct9620/pdf64/commit/25b81ff663aa470b173f3e9a5414e2213b85ffc3))

## [1.0.2](https://github.com/elct9620/pdf64/compare/v1.0.1...v1.0.2) (2025-03-08)


### Miscellaneous Chores

* release 1.0.2 ([60ed123](https://github.com/elct9620/pdf64/commit/60ed123132f4561c61ffa39331e893683de570b0))

## [1.0.1](https://github.com/elct9620/pdf64/compare/v1.0.0...v1.0.1) (2025-03-08)


### Miscellaneous Chores

* release 1.0.1 ([edd5a72](https://github.com/elct9620/pdf64/commit/edd5a72a18c8bcf36ca755bc9e0844a39cf30d48))

## 1.0.0 (2025-03-08)


### Features

* Add 'latest' tag for main branch Docker image ([7a72aef](https://github.com/elct9620/pdf64/commit/7a72aef6bb04060469da595bdf0f8bb2dfd754f7))
* Add conversion utilities for v1 API package ([cfbbba8](https://github.com/elct9620/pdf64/commit/cfbbba8bc36cbfc9cca49f658644ea518e783659))
* Add density and quality options to image conversion process ([1de4d85](https://github.com/elct9620/pdf64/commit/1de4d85d8db4a62669078943ff292832ae8e5456))
* Add devbox configuration with PDF and image processing tools ([7cdc866](https://github.com/elct9620/pdf64/commit/7cdc86670c37132750430d0336b2d3e7564eb087))
* Add Dockerfile and update go.mod for deployment ([90ed320](https://github.com/elct9620/pdf64/commit/90ed320edcb3b7bb047e6f8883035e920b1808f9))
* Add encryption status and methods to File entity ([3ed6466](https://github.com/elct9620/pdf64/commit/3ed64664d263162d8b3229772b9ea1d9a306333f))
* Add error handling package for v1 API ([cfeeaa7](https://github.com/elct9620/pdf64/commit/cfeeaa70996b00d7ae892f7af75e02c39cb6b753))
* Add fallback to 'convert' command for PDF image conversion in ImageMagick ([d8cac21](https://github.com/elct9620/pdf64/commit/d8cac21ccbb2dd5837acf3f24f7065137e183d3d))
* Add File entity and FileBuilder interface ([5c9af30](https://github.com/elct9620/pdf64/commit/5c9af30f7daa951afc860e9cd7f125deebc168b5))
* Add file upload handling with temporary file storage for ConvertUsecase ([9d5c812](https://github.com/elct9620/pdf64/commit/9d5c812a91e551efd5e4af324a40367772c032f0))
* Add file.go to internal/builder package ([31e22f0](https://github.com/elct9620/pdf64/commit/31e22f0579204d3b15e81843793528ea64772403))
* Add FileBuilder implementation with BuildFromPath method ([14d88b5](https://github.com/elct9620/pdf64/commit/14d88b5900a44fdec0d7e576acd66873542a6387))
* Add HTTP server with PDF conversion endpoint using Chi router ([c6f70f4](https://github.com/elct9620/pdf64/commit/c6f70f4038bf4a91c98a0bbce327c60f0e5311d3))
* Add ImageConvertService interface to define image conversion usecase ([f957418](https://github.com/elct9620/pdf64/commit/f9574185b3d4882b5d93f1db7e0c99ca62b98ee0))
* Add ImageMagick PDF conversion service with base64 encoding ([b3090cf](https://github.com/elct9620/pdf64/commit/b3090cf6e29e22731539a17cfda26ffefb7668a8))
* Add initial implementation of PDF conversion usecase and controller ([614969f](https://github.com/elct9620/pdf64/commit/614969f3da4766bbeb9d6d24a5e9498ac5e4b234))
* Add jmespath library for JSON response validation in tests ([744a051](https://github.com/elct9620/pdf64/commit/744a051e1bed6ea321c6f59b4ca3afef4e21ca3c))
* Add optional password parameter to convert request ([446cb52](https://github.com/elct9620/pdf64/commit/446cb52420fa68d33b4948481680c8900ce85d0a))
* Add PDF decryption support with password validation ([6b7609b](https://github.com/elct9620/pdf64/commit/6b7609b2c15fa5eba23ca298c7f765f88d65ff56))
* Add structured logging with httplog and enhance error logging ([3f4d3c3](https://github.com/elct9620/pdf64/commit/3f4d3c3f7d097f45f2cc0306a6c63cfab1a6d375))
* Implement `PostConvert` handler with form input validation and file upload ([76e07b8](https://github.com/elct9620/pdf64/commit/76e07b8c21efc3992d70743115da4d902a57355e))
* Implement PDF decryption service using qpdf ([9b02bf6](https://github.com/elct9620/pdf64/commit/9b02bf6f87edf9e3ea6c2b4368d5971cb2ac3f37))
* Initialize Go project with basic Hello World application ([724f669](https://github.com/elct9620/pdf64/commit/724f669e3b71934e902980ef87112a1c402063e5))


### Bug Fixes

* Add --allow-weak-crypto flag to qpdf command in test for encrypted PDF ([36c238d](https://github.com/elct9620/pdf64/commit/36c238dd8b78da5f71da48737087ef5b64d6f0c5))
* Add missing ImageConvertService parameter in NewConvertUsecase ([c303a73](https://github.com/elct9620/pdf64/commit/c303a73d7d57208c61facadd0c08bf0efd13cf0f))
* Add MockFileBuilder to support encrypted file testing ([57c2154](https://github.com/elct9620/pdf64/commit/57c2154f54d8b9463035388136be3e9039e1a7b4))
* Add qpdf dependency to Dockerfile for PDF decryption support ([07acf17](https://github.com/elct9620/pdf64/commit/07acf177d3b5aa45bad07f6981a6bc1ee701709d))
* Correct spelling of "whole" in CONVENTIONS.md ([8012726](https://github.com/elct9620/pdf64/commit/80127266fb067d50f8ae40e73082123d14eef16b))
* Correct YAML indentation in GitHub workflow test configuration ([c15009f](https://github.com/elct9620/pdf64/commit/c15009f981a88c881a4e061e0963ac3222116314))
* Handle password required error and use UUID in tests ([3bd1139](https://github.com/elct9620/pdf64/commit/3bd113914ca1a4c0f39e5c0984c6f938561cc7ba))
* Remove redundant err variable declaration in file_test.go ([aad7ada](https://github.com/elct9620/pdf64/commit/aad7ada020fbb96af50f126a831041675cd59283))
* Resolve linting errors in api_convert_test.go ([8057b98](https://github.com/elct9620/pdf64/commit/8057b98591c68c919136ff7a80f50caec71b1531))
* Resolve PDF file path in ImageMagick convert service test ([1b989ca](https://github.com/elct9620/pdf64/commit/1b989caeac46b2a9c8d96ab556b3b58c68c4dbba))
* Resolve test failures in API conversion tests ([4123305](https://github.com/elct9620/pdf64/commit/41233051d4720766836e0e34a8f8d1deb0cbf35f))
* Update GitHub Actions to install and configure ImageMagick 7 for PDF conversion ([73a3a6b](https://github.com/elct9620/pdf64/commit/73a3a6bcb92edc7585332f0cd2952a443acc18d5))
* Update ImageMagick verification command to use magick --version ([a6a6ff0](https://github.com/elct9620/pdf64/commit/a6a6ff05c0bf9196acdc9c8bb3fbf5e9885255e7))
* Update ImageMagickConvertService test to validate file paths instead of base64 images ([e84e73f](https://github.com/elct9620/pdf64/commit/e84e73f292c0ba2d88f57e032ed7b0c621c31833))
* Update qpdf encryption test to handle different versions and fallback scenarios ([c9a4554](https://github.com/elct9620/pdf64/commit/c9a45548d54410bab75619c53198ba1372b16b78))
* Update test to use project root fixtures/dummy.pdf ([ef7375c](https://github.com/elct9620/pdf64/commit/ef7375c7a6ae6f2eb04f28207f551cc1a5d8f20c))
