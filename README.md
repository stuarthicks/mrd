# MRD

mrd ("drm" backwards), pronounced /mɛː(ɹ)d/

Decoding the primordial soup of video DRM metadata. Inspired by https://walkergriggs.com/2024/10/16/pssh_primordial_soup_of_secure-ish_headers/

## Install

Using Homebrew:

    brew install stuarthicks/brews/mrd

Using Go:

    go install github.com/stuarthicks/mrd@latest

## Usage

References:

- https://github.com/shaka-project/shaka-packager/blob/main/packager/media/base/widevine_pssh_data.proto
- https://learn.microsoft.com/en-us/playready/specifications/playready-header-specification

Either pipe input via STDIN, or use `-input [file]` to specify a filename. Accepts mp4 and raw pssh/pro paylods (with or without base64 encoding).

Use `-verbose` for debug logging.

Use `-pretty=false` to disable pretty-printing (one line per header object printed).

> [!NOTE]
> Currently only PlayReady Objects that contain a single PlayReady Object Record are supported. If multiple records are present, the parsing behaviour is undefined and will likely fail.
