// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &readOnlyMap{}

type readOnlyMap struct {
	vals *map[string]uint8
}

func (bbp readOnlyMap) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.vals != nil {
		buf[at] = 1
		at++
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(*bbp.vals)))
		at += 4
		for k2, v2 := range *bbp.vals {
			iohelp.WriteUint32Bytes(buf[at:], uint32(len(k2)))
			copy(buf[at+4:at+4+len(k2)], []byte(k2))
			at += 4 + len(k2)
			iohelp.WriteUint8Bytes(buf[at:], v2)
			at += 1
		}
	}
	return at
}

func (bbp *readOnlyMap) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.vals = new(map[string]uint8)
			ln1 := iohelp.ReadUint32Bytes(buf[at:])
			at += 4
			(*bbp.vals) = make(map[string]uint8,ln1)
			for i := uint32(0); i < ln1; i++ {
				k3, err := iohelp.ReadStringBytes(buf[at:])
				if err != nil {
					return err
				}
				at += 4 + len(k3)
				if len(buf[at:]) < 1 {
					return io.ErrUnexpectedEOF
				}
				((*bbp.vals))[k3] = iohelp.ReadUint8Bytes(buf[at:])
				at += 1
			}
		default:
			return nil
		}
	}
}

func (bbp *readOnlyMap) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.vals = new(map[string]uint8)
			ln2 := iohelp.ReadUint32Bytes(buf[at:])
			at += 4
			(*bbp.vals) = make(map[string]uint8,ln2)
			for i := uint32(0); i < ln2; i++ {
				k3 := iohelp.MustReadStringBytes(buf[at:])
				at += 4 + len(k3)
				((*bbp.vals))[k3] = iohelp.ReadUint8Bytes(buf[at:])
				at += 1
			}
		default:
			return
		}
	}
}

func (bbp readOnlyMap) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.vals != nil {
		w.Write([]byte{1})
		iohelp.WriteUint32(w, uint32(len(*bbp.vals)))
		for k2, v2 := range *bbp.vals {
			iohelp.WriteUint32(w, uint32(len(k2)))
			w.Write([]byte(k2))
			iohelp.WriteUint8(w, v2)
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *readOnlyMap) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.vals = new(map[string]uint8)
			ln3 := iohelp.ReadUint32(r)
			*bbp.vals = make(map[string]uint8)
			for i := uint32(0); i < ln3; i++ {
				k3 := iohelp.ReadString(r)
				(*bbp.vals)[k3] = iohelp.ReadUint8(r)
			}
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp readOnlyMap) Size() int {
	bodyLen := 5
	if bbp.vals != nil {
		bodyLen += 1
		bodyLen += 4
		for k2 := range *bbp.vals {
			bodyLen += 4 + len(k2)
			bodyLen += 1
		}
	}
	return bodyLen
}

func (bbp readOnlyMap) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makereadOnlyMap(r iohelp.ErrorReader) (readOnlyMap, error) {
	v := readOnlyMap{}
	err := v.DecodeBebop(r)
	return v, err
}

func makereadOnlyMapFromBytes(buf []byte) (readOnlyMap, error) {
	v := readOnlyMap{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakereadOnlyMapFromBytes(buf []byte) readOnlyMap {
	v := readOnlyMap{}
	v.MustUnmarshalBebop(buf)
	return v
}
