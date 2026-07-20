package main_test

import "testing"

var deviceRoutes = []route{}

func TestApiDevice(t *testing.T) {
	testAuthRoutes(t, "Device", deviceRoutes)

	login(t)

	t.Run("CreateDevice", testCreateDevice)
	t.Run("GetDevices", testGetDevices)
	t.Run("GetDevice", testGetDevice)
	t.Run("DeleteDevice", testDeleteDevice)
}

func testCreateDevice(t *testing.T) {

}

func testGetDevices(t *testing.T) {

}

func testGetDevice(t *testing.T) {

}

func testDeleteDevice(t *testing.T) {

}
