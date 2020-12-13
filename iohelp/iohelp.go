package iohelp

import (
	"encoding/binary"
	"io"
	"time"
)

func ReadString(r io.Reader) string {
	ln := uint32(0)
	binary.Read(r, binary.LittleEndian, &ln)
	data := make([]byte, ln)
	r.Read(data)
	return string(data)
}

func ReadTime(r io.Reader) time.Time {
	tm := int64(0)
	binary.Read(r, binary.LittleEndian, (&tm))
	tm *= 100
	t := time.Time{}
	if tm == 0 {
		return t
	}
	return time.Unix(0, tm).UTC()
}

func ReadGUID(r io.Reader) [16]byte {
	data := make([]byte, 16)
	r.Read(data)
	flipped := [16]byte{
		data[3], data[2], data[1], data[0],
		data[5], data[4],
		data[7], data[6],
		data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15],
	}
	return flipped
}

func WriteGUID(w io.Writer, guid [16]byte) {
	// 3 2 1 0
	// 5 4
	// 7 6
	// 8 9 10 11 12 13 14 15
	flipped := [16]byte{
		guid[3], guid[2], guid[1], guid[0],
		guid[5], guid[4],
		guid[7], guid[6],
		guid[8], guid[9], guid[10], guid[11], guid[12], guid[13], guid[14], guid[15],
	}
	w.Write(flipped[:])
}
