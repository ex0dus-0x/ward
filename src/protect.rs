//! Defines the `WardApp` interface for consuming a single binary executable. Implements the following
//! hardening workflow:
//!
//!     - compress the executable and embed within the PT_NOTE segment
//!     - generate a file signature for the given finalized binary

use std::fs::{self, File, OpenOptions};
use std::io::{self, Seek, SeekFrom};
use std::io::{Read, Write};
use std::path::PathBuf;

use goblin::elf::program_header;
use goblin::elf::Elf;

use flate2::read::DeflateEncoder;
use flate2::Compression;

/// Implements the `Protector` trait for goblin Elf objects, which extend their functionality
/// to include code injection and extraction.
trait Protector {
    fn inject(&mut self, hostpath: PathBuf, payload: Vec<u8>) -> io::Result<()>;
    fn extract(&self) -> ();
}

impl Protector for Elf<'_> {
    /// When called with a target input file to protect, `inject()` will convert it into compressed
    /// bytes for writing, and manipulate hte .note.ABI-tag section header and PT_NOTE segment to
    /// point to it for later recovery.
    fn inject(&mut self, hostpath: PathBuf, payload: Vec<u8>) -> io::Result<()> {
        // open protector file and read out contents, we want to write to the end of it
        // with the standard filesystem facilities.
        let mut f = OpenOptions::new().read(true).append(true).open(hostpath)?;
        let mut buffer = Vec::new();
        f.read_to_end(&mut buffer)?;

        // save file offset of file end for later
        let offset: u64 = f.seek(SeekFrom::End(0))?;

        // given the compressed payload, append to end of file
        f.write(&payload)?;

        // overwrite the .note.ABI-tag section header

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
        Ok(())
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
    signature: Vec<u8>,
}

impl WardApp {
    pub fn new(filepath: PathBuf) -> io::Result<Self> {
        // open target file and read as bytes
        let mut f = File::open(&filepath)?;
        let mut buffer = Vec::new();
        f.read_to_end(&mut buffer);

        // compress the contents into a payload
        // TODO: encrypt symmetrically with password if configured
        let mut _payload = Vec::new();
        let mut deflater = DeflateEncoder::new(&buffer[..], Compression::fast());
        let count = deflater.read(&mut _payload)?;
        let binbytes: Vec<u8> = _payload[0..count].to_vec();

        Ok(Self {
            filepath,
            binbytes,
            signature: Vec::new(),
        })
    }

    /// initializes a new protector application, which parses itself as an ELF binary, and reads
    /// the compressed code hidden in the PT_NOTE-based code cave.
    fn _init_protector_app(&self) -> PathBuf {
        todo!()
    }

    /// given a path to a target binary, create a protector app to encapsulate it and inject the
    /// binary into a PT_NOTE-based code cave for recovery and re-execution under a protected environment.
    pub fn protect(&self) {
        // initialize a new protector binary and open as ELF
        let protector: PathBuf = self._init_protector_app();
        let buffer = fs::read(&protector).unwrap();
        let mut elf = match Elf::parse(&buffer) {
            Ok(bin) => bin,
            Err(e) => {
                panic!(e);
            }
        };

        // inject the compressed binary into the protector binary
        elf.inject(protector, self.binbytes.clone());
    }
}
