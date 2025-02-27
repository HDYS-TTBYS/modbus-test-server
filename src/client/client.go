package main

import (
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	var client *modbus.ModbusClient
	var err error

	// for a TCP endpoint
	// (see examples/tls_client.go for TLS usage and options)
	client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://localhost:5502",
		Timeout: 1 * time.Second,
	})

	if err != nil {
		println(err)
	}

	err = client.Open()
	if err != nil {
		println(err)
	}

	// DO roop test
	now_begin := time.Now()
	println(now_begin.Format("2006-01-02 15:04:05"))
	for i := 0; i < 2000; i++ {
		_, _ = client.ReadCoil(uint16(i))
	}
	now_end := time.Now()
	println(now_end.Format("2006-01-02 15:04:05"))
	after := now_end.Sub(now_begin)
	println("ReadCoil roop time used:" + after.String())

	// DO batch test
	now_begin = time.Now()
	println(now_begin.Format("2006-01-02 15:04:05"))
	_, _ = client.ReadCoils(0, 2000)
	now_end = time.Now()
	println(now_end.Format("2006-01-02 15:04:05"))
	after = now_end.Sub(now_begin)
	println("ReadCoils batch time used:" + after.String())

	// DI roop test
	now_begin = time.Now()
	println(now_begin.Format("2006-01-02 15:04:05"))
	for i := 0; i < 2000; i++ {
		_, _ = client.ReadDiscreteInput(uint16(i))
	}
	now_end = time.Now()
	println(now_end.Format("2006-01-02 15:04:05"))
	after = now_end.Sub(now_begin)
	println("ReadDiscreteInput roop time used:" + after.String())

	// DI batch test
	now_begin = time.Now()
	println(now_begin.Format("2006-01-02 15:04:05"))
	_, _ = client.ReadDiscreteInputs(0, 2000)
	now_end = time.Now()
	println(now_end.Format("2006-01-02 15:04:05"))
	after = now_end.Sub(now_begin)
	println("ReadDiscreteInputs batch time used:" + after.String())

	// AO roop test
	now_begin = time.Now()
	println(now_begin.Format("2006-01-02 15:04:05"))
	for i := 0; i < 1000; i++ {
		_, _ = client.ReadRegister(uint16(i), modbus.HOLDING_REGISTER)
	}
	now_end = time.Now()
	println(now_end.Format("2006-01-02 15:04:05"))
	after = now_end.Sub(now_begin)
	println("ReadHoldingRegister roop time used:" + after.String())

	// AO batch test
	now_begin = time.Now()
	println(now_begin.Format("2006-01-02 15:04:05"))
	_, _ = client.ReadRegisters(0, 100, modbus.HOLDING_REGISTER)
	_, _ = client.ReadRegisters(100, 100, modbus.HOLDING_REGISTER)
	_, _ = client.ReadRegisters(200, 100, modbus.HOLDING_REGISTER)
	_, _ = client.ReadRegisters(300, 100, modbus.HOLDING_REGISTER)
	_, _ = client.ReadRegisters(400, 100, modbus.HOLDING_REGISTER)
	_, _ = client.ReadRegisters(500, 100, modbus.HOLDING_REGISTER)
	_, _ = client.ReadRegisters(600, 100, modbus.HOLDING_REGISTER)
	_, _ = client.ReadRegisters(700, 100, modbus.HOLDING_REGISTER)
	_, _ = client.ReadRegisters(800, 100, modbus.HOLDING_REGISTER)
	_, _ = client.ReadRegisters(900, 100, modbus.HOLDING_REGISTER)
	now_end = time.Now()
	println(now_end.Format("2006-01-02 15:04:05"))
	after = now_end.Sub(now_begin)
	println("ReadHoldingRegisters batch time used:" + after.String())

	// AI roop test
	now_begin = time.Now()
	println(now_begin.Format("2006-01-02 15:04:05"))
	for i := 0; i < 1000; i++ {
		_, _ = client.ReadRegister(uint16(i), modbus.INPUT_REGISTER)
	}
	now_end = time.Now()
	println(now_end.Format("2006-01-02 15:04:05"))
	after = now_end.Sub(now_begin)
	println("ReadInputRegister roop time used:" + after.String())

	// AO batch test
	now_begin = time.Now()
	println(now_begin.Format("2006-01-02 15:04:05"))
	_, _ = client.ReadRegisters(0, 100, modbus.INPUT_REGISTER)
	_, _ = client.ReadRegisters(100, 100, modbus.INPUT_REGISTER)
	_, _ = client.ReadRegisters(200, 100, modbus.INPUT_REGISTER)
	_, _ = client.ReadRegisters(300, 100, modbus.INPUT_REGISTER)
	_, _ = client.ReadRegisters(400, 100, modbus.INPUT_REGISTER)
	_, _ = client.ReadRegisters(500, 100, modbus.INPUT_REGISTER)
	_, _ = client.ReadRegisters(600, 100, modbus.INPUT_REGISTER)
	_, _ = client.ReadRegisters(700, 100, modbus.INPUT_REGISTER)
	_, _ = client.ReadRegisters(800, 100, modbus.INPUT_REGISTER)
	_, _ = client.ReadRegisters(900, 100, modbus.INPUT_REGISTER)
	now_end = time.Now()
	println(now_end.Format("2006-01-02 15:04:05"))
	after = now_end.Sub(now_begin)
	println("ReadInputRegister batch time used:" + after.String())

	// close the TCP connection/serial port
	err = client.Close()
	if err != nil {
		println(err)
	}
}
