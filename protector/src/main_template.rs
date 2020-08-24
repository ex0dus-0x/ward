use goblin::Object;
use goblin::elf::Elf;

use std::fs;
use std::env;
use std::path::Path;

fn main() {
    let args = env::args().collect::<Vec<String>>();

    // set program to parse as this current one
    let prog: &Path = Path::new(&args[0]);
    let buffer = fs::read(prog)?;

    // parse out Elf binary format
    let elf = match Elf::parse(&buffer) {
        Ok(bin) => bin,
        Err(e) => {
            panic!(e);
        }
    };
}


