//! Command line application for `ward`. Interfaces all of the necessary functionality, including
//! binary signature generation and validation, and using `goblin` in order to wrangle and protect
//! target binary inputs

use goblin::Object;

use clap::{App, AppSettings, Arg, ArgMatches, SubCommand};

use std::fs;
use std::path::Path;
use std::error::Error;

use self::app::WardApp;

fn parse_args<'a>() -> ArgMatches<'a> {
    App::new(env!("CARGO_PKG_NAME"))
        .version(env!("CARGO_PKG_VERSION"))
        .author(env!("CARGO_PKG_AUTHORS"))
        .about(env!("CARGO_PKG_DESCRIPTION"))
        .setting(AppSettings::ArgRequiredElseHelp)

        // `protect` subcommand for incorporating the target ELF(s) with the protection runtime
        .subcommand(SubCommand::with_name("protect")
            .about("Statically injects binary into a protected runtime.")
            .arg(Arg::with_name("BINARY")
                .help("Path to compiled binary or binaries to protect")
                .index(1)
                .multiple(true)
                .required(true),
            )
            .arg(Arg::with_name("keep")
                .short("k")
                .long("keep")
                .help("Keeps the originals that were protected, and renames finalized secured binary.")
                .required(false),
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

    // open or initialize microkv store to interact with file signatures

    // parse subcommands
    match args.subcommand() {
        ("protect", Some(subargs)) => {
            let bins: Vec<&str> = subargs.values_of("BINARY").unwrap().collect();
            let binpaths: Vec<&Path> = bins
                .iter()
                .map(|b| Path::new(b))
                .collect();

            // given path in target binary list, initialize a new `WardApp` to protect
            for path in bins.iter() {
                let app: App = WardApp::init(path);
            }
        },
        ("verify", Some(subargs)) => {
            let bins: Vec<&str> = subargs.values_of("BINARY").unwrap().collect();
            let binpaths: Vec<&Path> = bins
                .iter()
                .map(|b| Path::new(b))
                .collect();

            // given path in target binary list, parse out objects
            for path in bins.iter() {
                let buffer = fs::read(path)?;
                let bin = match Object::parse(&buffer)? {
                    Object::Elf(elf) => elf,
                    _ => {
                        panic!("unsupported");
                    }
                };
            }
        },
        _ => todo!()
    }
    Ok(())
}


fn main() {
    match run() {
        Err(e) => eprintln!("{}", e),
        _ => {},
    }
}
