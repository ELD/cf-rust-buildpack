#![feature(plugin)]
#![plugin(rocket_codegen)]
extern crate rocket;
extern crate librocketapp;

fn main() {
    rocket::ignite().mount("/", routes![librocketapp::index]).launch();
}
