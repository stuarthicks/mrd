# soup
Decoding the primordial soup of video DRM metadata

Inspired by https://walkergriggs.com/2024/10/16/pssh_primordial_soup_of_secure-ish_headers/

## Usage

### Widevine PSSH

From a file:

    soup -pssh widevine_pssh.bin

From stdin:

    cat widevine_pssh.bin | soup -pssh -

### PlayReady Object

From a file:

    soup -pro pro.bin

From stdin:

    cat pro.bin | soup -pro -

> [!NOTE]  
> Currently only PlayReady Objects that contain a single PlayReady Object Record are supported. If multiple records are present, the parsing behaviour is undefined and will likely fail.
