#![feature(plugin)]
#![plugin(rocket_codegen)]
extern crate rocket;

#[get("/")]
fn index() -> &'static str {
    "Hello, Rocket, from Cloud Foundry!"
}
