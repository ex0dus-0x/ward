# ward

An experimental security-hardened notary for Linux binaries. We convert commonly weaponized ELF infection techniques into a defensive mechanism for runtime application self-protection, without the need of whole-system security provenance. Inspired by application notarization on macOS.

## Technique

1. Compile a protector runtime app, which employs several checks to ensure implants aren't attempting to inject themselves into the app

2. `ward` application uses the `PT_NOTE` code cave injection technique to hide a compressed blob of the original target binary

3. The protector app will read from itself during runtime, decompress the blob and use `memfd_create` to execute the original executable in-memory.

## Disclaimer

* Most stuff you try to protect will probably break
* Probably not the most resilient mitigation against adversarials.
