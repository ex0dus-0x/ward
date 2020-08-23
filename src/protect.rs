//! Defines the `WardApp` interface for consuming a single binary executable. Implements the following
//! hardening workflow:
//!
//!     - compress the executable and embed within the PT_NOTE segment
//!     - generate a file signature for the given finalized binary

use std::path::PathBuf;

use goblin::elf::Elf;
use goblin::elf::program_header;


/// Implements the `Protector` trait for goblin Elf objects, which extend their functionality
/// to include code injection and extraction.
trait Protector {
    fn inject(&mut self, payload: Vec<u8>) -> ();
    fn extract(&self) -> ();
}

impl Protector for Elf<'_> {
    fn inject(&mut self, payload: Vec<u8>) -> () {

        // iterate over the program header and change the PT_NOTE segment
        //  - change to a PT_LOAD segment since it won't be vacant
        //  - allow read and executable (TODO: maybe not necessary)
        for ph in self.program_headers.iter_mut() {
            if ph.p_type == program_header::PT_NOTE {
                ph.p_type = program_header::PT_LOAD;
                ph.p_flags = program_header::PF_R | program_header::PF_X;
                ph.p_vaddr = (0xc000000 + payload.len()) as u64;
            }
        }

        // once done, inject a new section at the end of the binary with the payload
    }


    fn extract(&self) -> () {
        todo!();
    }
}


/// Implements an `WardApp`, which encapsulates the functionality to protect a single executable
/// consumed.
#[derive(Debug)]
pub struct WardApp {
    filepath: PathBuf,
    binbytes: Vec<u8>,
    //signature:
}

impl WardApp {

    /// initializes a new protector application, which parses itself as an ELF binary, and reads
    /// the compressed code hidden in the PT_NOTE-based code cave.
    fn _create_protector() -> () {
        // mutate source in `protector` executable given the

        todo!()
    }

    pub fn init(filepath: PathBuf) -> Self {
        // compress the given binary into data that we can embed and re-extract later

        // initialize a new protector binary and open as ELF

        // inject the compressed binary into the protector binary

        // generate a signature given a derived keypair
        /*
        Self {
            path,
        }
        */
        todo!()
    }
}
