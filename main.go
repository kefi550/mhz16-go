package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	port = kingpin.Arg("port", "portName").Required().ExistingFile()
	zero = kingpin.Flag("zero", "zero calibration mode").Bool()
)

func main() {
	kingpin.Parse()
	options := &serial.Config{
		Name:        *port,
		Baud:        9600,
		ReadTimeout: 0,
	}
	port, err := serial.OpenPort(options)
	if err != nil {
		log.Fatalln("serial port cannot open: %w", err)
	}
	time.Sleep(time.Second * 1)
	defer port.Close()

	if *zero {
		zeroCalibration(*port)
		return
	}
	co2 := getCo2(*port)
	fmt.Println(co2)
}

func getCo2(port serial.Port) int {
	command := [9]byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
	n, err := port.Write(command[:])
	if err != nil {
		log.Fatalln("serial write error: %w", err)
	}
	res := make([]byte, 9)
	readed := 0
	for {
		n, err = port.Read(res[readed:9])
		if err != nil {
			log.Fatalln("serial read error: %w", err)
		}
		if n == 0 {
			break
		}
		readed += n
		if readed >= 9 {
			break
		}
	}
	checksum := 0xff & (^(res[1] + res[2] + res[3] + res[4] + res[5] + res[6] + res[7]) + 1)
	if res[8] != checksum {
		log.Fatalln("checksum not match")
	}
	result := int(res[2]) << 8 + int(res[3])
	return result
}

func zeroCalibration(port serial.Port) error {
	command := [9]byte{0xff, 0x01, 0x87, 0x00, 0x00, 0x00, 0x00, 0x00, 0x78}
	n, err := port.Write(command[:])
	if err != nil || n != 9 {
		log.Fatalln("serial write error: %w", err)
	}
	return err
}
