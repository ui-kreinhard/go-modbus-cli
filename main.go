package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"github.com/goburrow/modbus"
	"fmt"
	"log"
	"math"
	"strconv"
)

func prepareFlags() (*string, *int, *string, *int, *int) {
	host := flag.String("host", "", "Defines the host to be addressed")
	port := flag.Int("port", 502, "Can be used to override the port")
	operation := flag.String("operation", "", "Defines which operation should be performed")
	address := flag.Int("address", -1, "Address to be used")
	size := flag.Int("size", -1, "Size for reading how many coils/registers etc")

	return host, port, operation, address, size
}

func convertToByteArray(values []string) ([]byte, error) {
	valuesAsByte := []byte{}
	for _, value := range values {
		valueInt, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return nil, err
		}
		valuesAsByte = append(valuesAsByte, byte(valueInt))
	}
	return valuesAsByte, nil
}

func main() {
	host, port, operation, address, size := prepareFlags()
	flag.Parse()
	client, err := initConnection(*host, *port)
	valuesAsString := flag.Args()

	if err != nil {
		log.Println(err)
		return
	}
	resp := []byte{}
	switch *operation {
	case "writeSingleCoil":
		resp, err = writeSingleCoil(client, *address, valuesAsString)
		break
	case "writeMultipleRegisters":
		resp, err = writeMultipleRegisters(client, *address, valuesAsString)
		break
	case "writeMultipleCoils":
		resp, err = writeMultipleCoils(client, *address, valuesAsString)
		break
	case "writeSingleRegister":
		resp, err = writeSingleRegisters(client, *address, valuesAsString)
		break
	case "readCoils":
		resp, err = readCoils(client, *address, *size)
	default:
		fmt.Println("Unknown operation - aborting")
		fmt.Println("Known operations:")
		fmt.Println(" writeSingleCoil")
		fmt.Println(" writeMultipleRegisters")
	 	fmt.Println(" writeMultipleCoils")
	 	fmt.Println(" writeSingleRegister")
	 	fmt.Println(" readCoils")
	 }
	fmt.Println("response", resp)
	if err != nil {
		fmt.Println("err", err)
	}
}

func initConnection(host string, port int) (modbus.Client, error) {
	hostPort := host + ":" + strconv.Itoa(port)
	handler := modbus.NewTCPClientHandler(hostPort)
	err := handler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(handler)

	return client, nil
}

func convertStringArrayToInt(arrayToConvert []string) (int, error) {
	ret := 0
	for i, toConvert := range arrayToConvert {
		value, err := strconv.ParseInt(toConvert, 10, 2)
		if err != nil {
			return -1, err
		}
		elem := value << i
		ret = ret | int(elem)
	}
	return ret, nil
}

func writeSingleCoil(client modbus.Client, address int, valueAsString []string) ([]byte, error) {
	values, err := convertToByteArray(valueAsString)
	if err != nil {
		return nil, err
	}
	if values[0] >= 1 {
		return client.WriteSingleCoil(uint16(address), 0xFF00)
	}
	return client.WriteSingleCoil(uint16(address), 0x0000)
}

func convertUInt16TArray(values []string, bitsize int) ([]byte, error) {
	var valuesAsInt []uint64
	for _, value := range values {
		parsedValue, err := strconv.ParseUint(value, 10, bitsize)
		if err != nil {
			return nil, err
		}
		valuesAsInt = append(valuesAsInt, parsedValue)
	}

	byteBuffer := new(bytes.Buffer)
	for _, valueAsInt := range valuesAsInt {
		err := binary.Write(byteBuffer, binary.BigEndian, uint16(valueAsInt))
		if err != nil {
			return nil, err
		}
	}
	return byteBuffer.Bytes(), nil
}

func writeMultipleRegisters(client modbus.Client, address int, values []string) ([]byte, error) {
	valuesBytes, err := convertUInt16TArray(values, 16)
	if err != nil {
		return nil, err
	}
	return client.WriteMultipleRegisters(uint16(address), uint16(len(valuesBytes)/2), valuesBytes)
}

func writeSingleRegisters(client modbus.Client, address int, values []string) ([]byte, error) {
	valuesBytes, err := convertUInt16TArray(values, 16)
	if err != nil {
		return nil, err
	}
	return client.WriteMultipleRegisters(uint16(address), uint16(len(valuesBytes)/2), valuesBytes)
}

func reverse(a []byte) []byte {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}

func getLen(length int) int {
	return int(math.Ceil(float64(length) / 8.0))
}

func writeMultipleCoils(client modbus.Client, address int, values []string) ([]byte, error) {
	valuesAsInt, err := convertStringArrayToInt(values)
	if err != nil {
		return nil, err
	}
	byteBuffer := new(bytes.Buffer)

	err = binary.Write(byteBuffer, binary.BigEndian, uint64(valuesAsInt))
	if err != nil {
		return nil, err
	}

	modbusBytes := reverse(byteBuffer.Bytes())[0:getLen(len(values))]

	fmt.Println("byte", modbusBytes[0:getLen(len(values))])

	return client.WriteMultipleCoils(uint16(address), uint16(len(values)), modbusBytes)
}

func readCoils(client modbus.Client, address int, size int) (results []byte, err error) {
	return client.ReadCoils(uint16(address), uint16(size))
}
