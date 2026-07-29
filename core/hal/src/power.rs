use crate::{HalError, HalResult};

pub enum PowerState {
    On,
    Off,
    Suspend,
    Hibernate,
    Rebooting,
}

pub trait PowerHal {
    fn shutdown(&self) -> HalResult<()>;

    fn reboot(&self) -> HalResult<()>;

    fn suspend(&self) -> HalResult<()>;

    fn hibernate(&self) -> HalResult<()>;

    fn power_state(&self) -> HalResult<PowerState>;
}

pub struct DefaultPowerHal;

impl PowerHal for DefaultPowerHal {
    fn shutdown(&self) -> HalResult<()> {
        Err(HalError::NotSupported(
            "power management not available in this context".into(),
        ))
    }

    fn reboot(&self) -> HalResult<()> {
        Err(HalError::NotSupported(
            "power management not available in this context".into(),
        ))
    }

    fn suspend(&self) -> HalResult<()> {
        Err(HalError::NotSupported(
            "power management not available in this context".into(),
        ))
    }

    fn hibernate(&self) -> HalResult<()> {
        Err(HalError::NotSupported(
            "power management not available in this context".into(),
        ))
    }

    fn power_state(&self) -> HalResult<PowerState> {
        Ok(PowerState::On)
    }
}
