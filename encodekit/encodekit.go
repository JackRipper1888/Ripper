package encodekit

import (
	"bytes"
	"encoding/binary"
)

func BinaryWrite(t interface{}) []byte {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, t)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func BinaryRead(b []byte, t interface{}) error {
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.BigEndian, t)
	if err != nil {
		panic(err)
	}
	return err
}

func BinaryVarint(buf []byte) int64 {
	body, n := binary.Varint(buf)
	if n != len(buf) {
		//panic("Varint did not consume all of in")
		panic(body)
	}
	return body
}

func BytesToInt32(b []byte) int32 {
    bytesBuffer := bytes.NewBuffer(b)
    var x int32
	err := binary.Read(bytesBuffer, binary.BigEndian, &x)
	if err != nil{
		panic(err)
	}
    return  x
}
func BytesToInt64(b []byte) int64 {
    bytesBuffer := bytes.NewBuffer(b)
    var x int64
	err := binary.Read(bytesBuffer, binary.BigEndian, &x)
	if err != nil{
		panic(err)
	}
    return  x
}