use crate::error::HalError;

pub type HalResult<T> = Result<T, HalError>;

pub fn enumerate() -> HalResult<Vec<String>> {
    Err(HalError::NotSupported)
}

pub fn read(device: &str, offset: u64, buf: &mut [u8]) -> HalResult<usize> {
    let _ = (device, offset, buf);
    Err(HalError::NotSupported)
}

pub fn write(device: &str, offset: u64, buf: &[u8]) -> HalResult<usize> {
    let _ = (device, offset, buf);
    Err(HalError::NotSupported)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_enumerate_not_supported() {
        assert!(matches!(enumerate(), Err(HalError::NotSupported)));
    }
}
