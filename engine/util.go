package engine

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/davecgh/go-spew/spew"
	"github.com/robertkrimen/otto"
)

func (e *Engine) ValueToByteSlice(v otto.Value) []byte {
	valueBytes := []byte{}
	if v.IsNull() || v.IsUndefined() {
		return valueBytes
	}
	if v.IsString() {
		str, err := v.Export()
		if err != nil {
			e.Logger.WithField("trace", "true").Errorf("Cannot convert string to byte slice: %s", v)
			return valueBytes
		}
		valueBytes = []byte(str.(string))
	} else if v.IsNumber() {
		num, err := v.Export()
		if err != nil {
			e.Logger.WithField("trace", "true").Errorf("Cannot convert string to byte slice: %s", v)
			return valueBytes
		}
		buf := new(bytes.Buffer)
		err = binary.Write(buf, binary.LittleEndian, num)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
		valueBytes = buf.Bytes()
	} else if v.Class() == "Array" || v.Class() == "GoArray" {
		arr, err := v.Export()
		if err != nil {
			e.Logger.WithField("trace", "true").Errorf("Cannot convert array to byte slice: %x", v)
			return valueBytes
		}
		switch t := arr.(type) {
		case []uint:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []uint8:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []uint16:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []uint32:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []uint64:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []int:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []int16:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []int32:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []int64:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []float32:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []float64:
			for _, i := range t {
				valueBytes = append(valueBytes, byte(i))
			}
		case []string:
			for _, i := range t {
				for _, c := range i {
					valueBytes = append(valueBytes, byte(c))
				}
			}
		default:
			_ = t
			e.Logger.WithField("trace", "true").Errorf("Failed to cast array to byte slice array=%v", arr)
		}
	} else {
		e.Logger.WithField("trace", "true").Errorf("Unknown class to cast to byte slice")
		spew.Dump(v)
	}

	return valueBytes
}

func GetLocalIPs() []string {
	addresses := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return addresses
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addresses = append(addresses, ipnet.IP.String())
			}
		}
	}
	return addresses
}
