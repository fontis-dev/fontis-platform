#[derive(Debug)]
pub enum HalError {
    NotSupported,
    IoError(String),
    TpmError(String),
    PowerError(String),
}

impl core::fmt::Display for HalError {
    fn fmt(&self, f: &mut core::fmt::Formatter<'_>) -> core::fmt::Result {
        match self {
            HalError::NotSupported => write!(f, "operation not supported"),
            HalError::IoError(msg) => write!(f, "I/O error: {}", msg),
            HalError::TpmError(msg) => write!(f, "TPM error: {}", msg),
            HalError::PowerError(msg) => write!(f, "power error: {}", msg),
        }
    }
}

impl std::error::Error for HalError {}
