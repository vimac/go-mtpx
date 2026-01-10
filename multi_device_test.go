package mtpx

import (
	"testing"

	"github.com/ganeshrvel/go-mtpfs/mtp"
	. "github.com/smartystreets/goconvey/convey"
)

func TestListDevices(t *testing.T) {
	Convey("Testing ListDevices", t, func() {
		devices, err := ListDevices()

		So(err, ShouldBeNil)
		So(devices, ShouldNotBeNil)

		// If devices are connected, verify the structure
		if len(devices) > 0 {
			device := devices[0]
			So(device.VendorId, ShouldBeGreaterThan, 0)
			So(device.ProductId, ShouldBeGreaterThan, 0)
		}
	})
}

func TestInitializeByDeviceId(t *testing.T) {
	var dev *mtp.Device
	var vendorId, productId uint16
	var serialNumber string

	Convey("Testing InitializeByDeviceId with valid device", t, func() {
		// First, list devices to get valid device identifiers
		devices, err := ListDevices()
		So(err, ShouldBeNil)

		if len(devices) == 0 {
			SkipSo(nil, ShouldBeNil) // Skip if no devices connected
			return
		}

		vendorId = devices[0].VendorId
		productId = devices[0].ProductId
		serialNumber = devices[0].SerialNumber

		// Initialize by the device identifiers
		d, err := InitializeByDeviceId(vendorId, productId, serialNumber, Init{})
		dev = d

		So(err, ShouldBeNil)
		So(d, ShouldNotBeNil)
		So(d.Timeout, ShouldBeGreaterThan, 1)

		// Get device info to verify connection
		info, err := FetchDeviceInfo(dev)
		So(err, ShouldBeNil)
		So(info, ShouldNotBeNil)

		Dispose(dev)
	})

	Convey("Testing InitializeByDeviceId with invalid device", t, func() {
		_, err := InitializeByDeviceId(0xFFFF, 0xFFFF, "", Init{})

		So(err, ShouldNotBeNil)
		_, ok := err.(DeviceNotFoundError)
		So(ok, ShouldBeTrue)
	})
}

func TestGetDeviceIdentifiers(t *testing.T) {
	Convey("Testing GetDeviceIdentifiers", t, func() {
		// First initialize a device
		dev, err := Initialize(Init{})
		if err != nil {
			SkipSo(nil, ShouldBeNil) // Skip if no devices connected
			return
		}

		// Get the device identifiers
		vendorId, productId, serialNumber, err := GetDeviceIdentifiers(dev)
		So(err, ShouldBeNil)
		So(vendorId, ShouldBeGreaterThan, 0)
		So(productId, ShouldBeGreaterThan, 0)
		// serialNumber may be empty on some devices
		_ = serialNumber

		Dispose(dev)
	})
}

func TestMultiDeviceWorkflow(t *testing.T) {
	Convey("Testing multi-device workflow", t, func() {
		// List all devices
		devices, err := ListDevices()
		So(err, ShouldBeNil)

		if len(devices) == 0 {
			SkipSo(nil, ShouldBeNil) // Skip if no devices connected
			return
		}

		// For each device, initialize and verify
		for _, deviceInfo := range devices {
			Convey("Device: "+deviceInfo.Manufacturer+" "+deviceInfo.Model, func() {
				dev, err := InitializeByDeviceId(deviceInfo.VendorId, deviceInfo.ProductId, deviceInfo.SerialNumber, Init{})
				So(err, ShouldBeNil)
				So(dev, ShouldNotBeNil)

				// Verify device identifiers match
				vendorId, productId, serialNumber, err := GetDeviceIdentifiers(dev)
				So(err, ShouldBeNil)
				So(vendorId, ShouldEqual, deviceInfo.VendorId)
				So(productId, ShouldEqual, deviceInfo.ProductId)
				So(serialNumber, ShouldEqual, deviceInfo.SerialNumber)

				// Get storages
				storages, err := FetchStorages(dev)
				So(err, ShouldBeNil)
				So(len(storages), ShouldBeGreaterThan, 0)

				Dispose(dev)
			})
		}
	})
}
