//! Command line application for `ward`. Interfaces all of the necessary functionality, including
//! binary signature generation and validation, and using `goblin` in order to wrangle and protect
//! target binary inputs

use goblin::Object;

use clap::{App, AppSettings, Arg, ArgMatches, SubCommand};

use std::fs;
use std::path::Path;
use std::error::Error;

fn parse_args<'a>() -> ArgMatches<'a> {
    App::new(env!("CARGO_PKG_NAME"))
        .version(env!("CARGO_PKG_VERSION"))
        .author(env!("CARGO_PKG_AUTHORS"))
        .about(env!("CARGO_PKG_DESCRIPTION"))
        .setting(AppSettings::ArgRequiredElseHelp)

        // `protect` subcommand for static patching input ELFs
        .subcommand(SubCommand::with_name("protect")
            .about("Statically patches input ELF binary/binaries with the protection runtime.")
            .arg(Arg::with_name("BINARY")
                .help("Path to compiled binary or binaries to protect")
                .index(1)
                .multiple(true)
                .required(true),
            )
        )

        // `verify` subcommand for checking if binary is secured
        .subcommand(SubCommand::with_name("verify")
            .about("Check to see if binary/binaries are protected and signature is preserved")
            .arg(Arg::with_name("BINARY")
                .help("Path to compiled binary or binaries to verify for protection")
                .index(1)
                .multiple(true)
                .required(true),
            ),
        )
        .get_matches()
}


fn run() -> Result<(), Box<dyn Error>> {
    let args: ArgMatches = parse_args();

    // parse subcommands
    match args.subcommand() {
        ("protect", Some(subargs)) => {
            let bins: Vec<&Path> = subargs.values_of("BINARY").unwrap()
                .collect::<Vec<&str>>()
                .iter()
                .map(|b| Path::new(b))
                .collect();

        }
        ("verify", Some(subargs)) => {
            todo!()
        },
        _ => todo!()

    }

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
