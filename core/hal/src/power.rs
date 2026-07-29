use crate::error::HalError;

pub type HalResult<T> = Result<T, HalError>;

pub fn shutdown() -> HalResult<()> {
    Err(HalError::NotSupported)
}

pub fn reboot() -> HalResult<()> {
    Err(HalError::NotSupported)
}

pub fn suspend() -> HalResult<()> {
    Err(HalError::NotSupported)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_shutdown_not_supported() {
        assert!(matches!(shutdown(), Err(HalError::NotSupported)));
    }

    #[test]
    fn test_reboot_not_supported() {
        assert!(matches!(reboot(), Err(HalError::NotSupported)));
    }
}
