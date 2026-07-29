#![deny(unsafe_code)]

pub mod block;
pub mod ffi;
pub mod power;
pub mod tpm;

use thiserror::Error;

#[derive(Error, Debug)]
pub enum HalError {
    #[error("device not found: {0}")]
    DeviceNotFound(String),

    #[error("I/O error: {0}")]
    Io(#[from] std::io::Error),

    #[error("TPM error: {0}")]
    TpmError(String),

    #[error("permission denied: {0}")]
    PermissionDenied(String),

    #[error("not supported: {0}")]
    NotSupported(String),
}

pub type HalResult<T> = Result<T, HalError>;
