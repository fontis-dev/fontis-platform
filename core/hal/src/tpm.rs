use crate::error::HalError;

pub type HalResult<T> = Result<T, HalError>;

pub fn pcr_read(index: u32) -> HalResult<[u8; 32]> {
    let _ = index;
    Err(HalError::TpmError("not implemented".into()))
}

pub fn pcr_extend(index: u32, value: &[u8]) -> HalResult<()> {
    let _ = (index, value);
    Err(HalError::TpmError("not implemented".into()))
}

pub fn seal(key_name: &str, data: &[u8]) -> HalResult<Vec<u8>> {
    let _ = (key_name, data);
    Err(HalError::TpmError("not implemented".into()))
}

pub fn unseal(key_name: &str, sealed: &[u8]) -> HalResult<Vec<u8>> {
    let _ = (key_name, sealed);
    Err(HalError::TpmError("not implemented".into()))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_pcr_read_not_implemented() {
        assert!(matches!(pcr_read(0), Err(HalError::TpmError(_))));
    }
}
