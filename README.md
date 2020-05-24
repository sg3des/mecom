# MeCom

Go implementation of the MeCom protocol by Meerstetter, to control TEC devices.

## Install

```sh
go get github.com/sg3des/mecom
```

## Usage

```go
ctrl, err := mecom.Dial("/dev/ttyUSB0")
// handle error

// get current temperature
temp, err := ctrl.GetObjectTemperature()

// set temperature
err := ctrl.SetTemperature(37.5)

// etc

```
