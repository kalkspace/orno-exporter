package orno

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/goburrow/modbus"
)

func Read(address string) error {
	handler := modbus.NewRTUClientHandler(address)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "E"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Logger = log.New(os.Stdout, "", 0)

	client := modbus.NewClient(handler)
	res, err := client.ReadHoldingRegisters(0xD, 16)
	if err != nil {
		return err
	}

	fmt.Printf("result: %+v\n", res)
	if len(res) >= 2 {
		num := binary.BigEndian.Uint16(res)
		fmt.Printf("num: %d\n", num)
	}
	return nil
}
