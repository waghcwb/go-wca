// +build !windows

package main

import (
	"github.com/go-ole/go-ole"
)

func aevRegisterControlChangeNotify() (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevUnregisterControlChangeNotify() (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevGetChannelCount(aev *IAudioEndpointVolume, channelCount *uint32) (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevSetMasterVolumeLevel() (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevSetMasterVolumeLevelScalar(aev *IAudioEndpointVolume, level float32, eventContextGUID *ole.GUID) (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevGetMasterVolumeLevel() (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevGetMasterVolumeLevelScalar(aev *IAudioEndpointVolume, level *float32) (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevSetChannelVolumeLevel() (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevSetChannelVolumeLevelScalar() (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevGetChannelVolumeLevel() (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevGetChannelVolumeLevelScalar(aev *IAudioEndpointVolume, channel uint32, level *float32) (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevSetMute(aev *IAudioEndpointVolume, mute bool, eventContextGUID *ole.GUID) (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevGetMute(aev *IAudioEndpointVolume, mute *bool) (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevGetVolumeStepInfo(aev *IAudioEndpointVolume, step, stepCount *uint32) (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevVolumeStepUp(aev *IAudioEndpointVolume, eventContextGUID *ole.GUID) (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevVolumeStepDown(aev *IAudioEndpointVolume, eventContextGUID *ole.GUID) (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevQueryHardwareSupport(aev *IAudioEndpointVolume, hardwareSupportMask *uint32) (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}

func aevGetVolumeRange(aev *IAudioEndpointVolume, minDB, maxDB, incrementDB *float32) (err error) {
	return ole.NewError(ole.E_NOTIMPL)
}
