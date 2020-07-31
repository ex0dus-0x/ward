use goblin::Object;

use std::fs;
use std::path::Path;
use std::error::Error;

fn run() -> Result<(), Box<dyn Error>> {
    let path = Path::new("./target/bin/ward");
    let buffer = fs::read(path)?;
    let bin = match Object::parse(&buffer)? {
        Object::Elf(elf) => elf,
        _ => {
            panic!("unsupported");
        }
    };
    Ok(())
}


fn main() {
    match run() {
        Err(e) => eprintln!("{}", e),
        _ => {},
    }
}
