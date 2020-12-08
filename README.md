# ward

An experimental security-hardened notary for Linux binaries. We convert a commonly weaponized ELF infection technique into a defensive mechanism for runtime application self-protection, without the need of whole-system security provenance. Inspired by application notarization on macOS.

## Technique

1. Compile a protector runtime app, which dynamically loads and executes the
