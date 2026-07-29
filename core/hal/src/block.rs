use crate::{HalError, HalResult};

#[derive(Debug, Clone, serde::Serialize)]
pub struct BlockDeviceInfo {
    pub path: String,
    pub size_bytes: u64,
    pub model: Option<String>,
    pub serial: Option<String>,
    pub is_rotational: bool,
}

pub trait BlockDeviceHal {
    fn enumerate_devices(&self) -> HalResult<Vec<BlockDeviceInfo>>;

    fn read(&self, path: &str, offset: u64, buffer: &mut [u8]) -> HalResult<usize>;

    fn write(&self, path: &str, offset: u64, data: &[u8]) -> HalResult<usize>;

    fn device_size(&self, path: &str) -> HalResult<u64>;
}

pub struct DefaultBlockDeviceHal;

impl BlockDeviceHal for DefaultBlockDeviceHal {
    fn enumerate_devices(&self) -> HalResult<Vec<BlockDeviceInfo>> {
        Err(HalError::NotSupported(
            "block device enumeration not available in this context".into(),
        ))
    }

    fn read(&self, _path: &str, _offset: u64, _buffer: &mut [u8]) -> HalResult<usize> {
        Err(HalError::NotSupported(
            "block device read not available in this context".into(),
        ))
    }

    fn write(&self, _path: &str, _offset: u64, _data: &[u8]) -> HalResult<usize> {
        Err(HalError::NotSupported(
            "block device write not available in this context".into(),
        ))
    }

    fn device_size(&self, _path: &str) -> HalResult<u64> {
        Err(HalError::NotSupported(
            "device size query not available in this context".into(),
        ))
    }
}
