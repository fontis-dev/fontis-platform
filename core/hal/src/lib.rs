#![deny(unsafe_code)]

pub mod block;
pub mod error;
pub mod ffi;
pub mod power;
pub mod tpm;

pub use error::HalError;
