//! Defines the `WardApp` interface for consuming a single binary executable. Implements the following
//! hardening workflow:
//!
//!     - compress the executable and embed within the PT_NOTE segment
//!     - generate a file signature for the given finalized binary

use std::fs::{self, File, OpenOptions};
use std::io::{self, Seek, SeekFrom};
use std::io::{Read, Write};
use std::path::{Path, PathBuf};

use goblin::elf::program_header;
use goblin::elf::Elf;

use flate2::write::GzEncoder;
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
        // with the standard filesystem facilities instead of libgoblin
        let mut f = OpenOptions::new().read(true).append(true).open(hostpath)?;
        let mut buffer = Vec::new();
        f.read_to_end(&mut buffer)?;

        // save file offset of file end for later
        let offset: u64 = f.seek(SeekFrom::End(0))?;

        // overwrite the .note.ABI-tag section header

        // iterate over the program header and change the PT_NOTE segment
        for ph in self.program_headers.iter_mut() {
            if ph.p_type == program_header::PT_NOTE {

                // change to a PT_LOAD segment since it won't be vacant
                ph.p_type = program_header::PT_LOAD;

                // save and modify the entry point to point to where payload should be
                ph.p_vaddr = (0xc000000 + payload.len()) as u64;
                self.entry = ph.p_vaddr;

                // modify filesz and memsz attributes
                ph.p_filesz += offset;
                ph.p_memsz += offset;

                // file offset change
                ph.p_offset = payload.len() as u64;
            }
        }

        // TODO: write new elf to new path

        // reopen elf, seek to end, and write payload
        let mut protected = OpenOptions::new().read(true).append(true).open(hostpath)?;
        let mut protectedbuf = Vec::new();
        protected.read_to_end(&mut protectedbuf)?;

        // given the compressed payload and new file, append to end of file
        protected.write(&payload)?;
        Ok(())
    }

    fn extract(&self, hostpath: PathBuf) -> () {
        // iterate over program header and identify new .injected section
        // find entry point and file seek
        todo!()
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
        f.read_to_end(&mut buffer)?;

        // TODO: encrypt symmetrically with password if configured
        let mut payload = Vec::new();
        let mut deflater = GzEncoder::new(payload, Compression::fast());
        deflater.write_all(&mut buffer)?;
        let binbytes: Vec<u8> = deflater.finish()?;

        Ok(Self {
            filepath,
            binbytes,
            signature: Vec::new(),
        })
    }


    /// given a path to a target binary, create a protector app to encapsulate it and inject the
    /// binary into a PT_NOTE-based code cave for recovery and re-execution under a protected environment.
    pub fn protect(&self) -> Result<(), ()> {

        // initialize path to protector application
        let protector: &Path = &Path::new("protector/protector");

        // initialize a new protector binary and open as ELF
        // TODO: if not exist, user did not invoke Makefile
        let buffer = fs::read(&protector.to_path_buf()).unwrap();
        let mut elf = match Elf::parse(&buffer) {
            Ok(bin) => bin,
            Err(e) => {
                panic!(e);
            }
        };

        // inject the compressed binary into the protector binary
        elf.inject(protector.to_path_buf(), self.binbytes.clone());

        Ok(())
    }
}
