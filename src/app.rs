//! Defines the `App` interface for consuming a single binary executable. Implements the following
//! hardening workflow:
//!
//!     - compress the executable and embed within the PT_NOTE segment
//!     - generate a file signature for the given finalized binary

use std::path::Path;


/// Implements an `App`, which encapsulates the functionality to protect a single executable
/// consumed. A user can instantiate a new `App` to protect an input, or consume an
/// already-protected application to validate it against a given signature.
#[derive(Debug)]
pub struct WardApp {
    path: Path,
    binbytes: Vec<u8>,
}

impl WardApp {

    /// initializes a new protector application, which parses itself as an ELF binary, and reads
    /// the compressed code hidden in the PT_NOTE-based code cave.
    fn create_protector() -> ! {

    }

    pub fn init(filepath: Path) -> Self {
        // compress the given binary into a
    }
}
