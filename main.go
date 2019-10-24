package main

import "github.com/kalkspace/orno-exporter/orno"

func main() {
	err := orno.Read("/dev/cu.usbserial-AM00EBGZ")
	if err != nil {
		panic(err)
	}
}
