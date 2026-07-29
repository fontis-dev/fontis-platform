use std::ffi::CString;
use std::os::raw::c_char;

use crate::block::{BlockDeviceHal, DefaultBlockDeviceHal};
use crate::tpm::{DefaultTpmHal, TpmHal};

static BLOCK_HAL: DefaultBlockDeviceHal = DefaultBlockDeviceHal;
static TPM_HAL: DefaultTpmHal = DefaultTpmHal;

#[no_mangle]
pub unsafe extern "C" fn hal_block_enumerate(
    result: *mut *mut c_char,
    result_len: *mut usize,
) -> i32 {
    match BLOCK_HAL.enumerate_devices() {
        Ok(devices) => {
            let json = serde_json::to_string(&devices).unwrap_or_default();
            let c_str = CString::new(json).unwrap_or_default();
            unsafe {
                *result = c_str.into_raw();
                *result_len = 0;
            }
            0
        }
        Err(e) => {
            let err_str = CString::new(e.to_string()).unwrap_or_default();
            unsafe {
                *result = err_str.into_raw();
                *result_len = 0;
            }
            -1
        }
    }
}

#[no_mangle]
pub extern "C" fn hal_tpm_is_present() -> i32 {
    if TPM_HAL.is_present() {
        1
    } else {
        0
    }
}

#[no_mangle]
pub unsafe extern "C" fn hal_tpm_read_pcr(index: u32, result: *mut *mut c_char) -> i32 {
    match TPM_HAL.read_pcr(index) {
        Ok(pcr) => {
            let json = serde_json::to_string(&pcr).unwrap_or_default();
            let c_str = CString::new(json).unwrap_or_default();
            unsafe {
                *result = c_str.into_raw();
            }
            0
        }
        Err(e) => {
            let err_str = CString::new(e.to_string()).unwrap_or_default();
            unsafe {
                *result = err_str.into_raw();
            }
            -1
        }
    }
}

#[no_mangle]
pub unsafe extern "C" fn hal_free_string(s: *mut c_char) {
    if !s.is_null() {
        unsafe {
            let _ = CString::from_raw(s);
        }
    }
}
