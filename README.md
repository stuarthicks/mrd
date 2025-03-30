# soup
Decoding the primordial soup of video DRM metadata

Inspired by https://walkergriggs.com/2024/10/16/pssh_primordial_soup_of_secure-ish_headers/

## Usage

### Widevine PSSH

Reference: https://github.com/shaka-project/shaka-packager/blob/main/packager/media/base/widevine_pssh_data.proto

Decode from a file:

    soup -pssh widevine_pssh.bin

Decode from STDIN:

    cat widevine_pssh.bin | soup -pssh -

### PlayReady Object

Reference: https://learn.microsoft.com/en-us/playready/specifications/playready-header-specification

Decode from a file:

    soup -pro pro.bin

Decode from STDIN:

    cat pro.bin | soup -pro -

> [!NOTE]
> Currently only PlayReady Objects that contain a single PlayReady Object Record are supported. If multiple records are present, the parsing behaviour is undefined and will likely fail.
