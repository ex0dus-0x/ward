use goblin::Object;
use goblin::elf::Elf;

use std::fs;
use std::env;
use std::error::Error;
use std::path::Path;

fn main() -> Result<(), Box<dyn Error>> {
    let args = env::args().collect::<Vec<String>>();
    if args.len() > 1 {
        panic!("No other arguments required. Run executable as is.");
    }

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

    Ok(())
}


