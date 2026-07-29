use crate::{HalError, HalResult};

#[derive(Debug, Clone, serde::Serialize)]
pub struct PcrValue {
    pub index: u32,
    pub value: Vec<u8>,
}

pub trait TpmHal {
    fn read_pcr(&self, index: u32) -> HalResult<PcrValue>;

    fn extend_pcr(&self, index: u32, data: &[u8]) -> HalResult<()>;

    fn seal(&self, data: &[u8], pcr_mask: u32) -> HalResult<Vec<u8>>;

    fn unseal(&self, sealed_data: &[u8]) -> HalResult<Vec<u8>>;

    fn is_present(&self) -> bool;
}

pub struct DefaultTpmHal;

impl TpmHal for DefaultTpmHal {
    fn read_pcr(&self, _index: u32) -> HalResult<PcrValue> {
        Err(HalError::NotSupported(
            "TPM not available in this context".into(),
        ))
    }

    fn extend_pcr(&self, _index: u32, _data: &[u8]) -> HalResult<()> {
        Err(HalError::NotSupported(
            "TPM not available in this context".into(),
        ))
    }

    fn seal(&self, _data: &[u8], _pcr_mask: u32) -> HalResult<Vec<u8>> {
        Err(HalError::NotSupported(
            "TPM not available in this context".into(),
        ))
    }

    fn unseal(&self, _sealed_data: &[u8]) -> HalResult<Vec<u8>> {
        Err(HalError::NotSupported(
            "TPM not available in this context".into(),
        ))
    }

    fn is_present(&self) -> bool {
        false
    }
}
