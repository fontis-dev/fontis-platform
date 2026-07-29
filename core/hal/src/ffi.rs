#![allow(unsafe_code)]

use crate::error::HalError;
use std::ffi::{CStr, CString};
use std::os::raw::c_char;

pub type HalResult<T> = Result<T, HalError>;

#[repr(C)]
pub struct HalFfiResult {
    pub success: bool,
    pub error_message: *mut c_char,
}

fn to_ffi_result(res: HalResult<()>) -> HalFfiResult {
    match res {
        Ok(()) => HalFfiResult {
            success: true,
            error_message: std::ptr::null_mut(),
        },
        Err(e) => {
            let msg = CString::new(e.to_string()).unwrap_or_default();
            HalFfiResult {
                success: false,
                error_message: msg.into_raw(),
            }
        }
    }
}

#[no_mangle]
pub unsafe extern "C" fn hal_shutdown() -> HalFfiResult {
    to_ffi_result(crate::power::shutdown())
}

#[no_mangle]
pub unsafe extern "C" fn hal_reboot() -> HalFfiResult {
    to_ffi_result(crate::power::reboot())
}

#[no_mangle]
pub unsafe extern "C" fn hal_free_error_message(ptr: *mut c_char) {
    if !ptr.is_null() {
        let _ = unsafe { CString::from_raw(ptr) };
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_shutdown_ffi_result_is_not_supported() {
        let result = unsafe { hal_shutdown() };
        assert!(!result.success);
        if !result.error_message.is_null() {
            let msg = unsafe { CStr::from_ptr(result.error_message) };
            assert!(!msg.to_str().unwrap().is_empty());
            unsafe { hal_free_error_message(result.error_message) };
        }
    }
}
