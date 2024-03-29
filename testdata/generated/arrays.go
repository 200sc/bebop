// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &ArraySamples{}

type ArraySamples struct {
	Bytes [][][]byte
	Bytes2 [][][]byte
}

func (bbp ArraySamples) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.Bytes)))
	at += 4
	for _, v1 := range bbp.Bytes {
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(v1)))
		at += 4
		for _, v2 := range v1 {
			iohelp.WriteUint32Bytes(buf[at:], uint32(len(v2)))
			at += 4
			copy(buf[at:at+len(v2)], v2)
			at += len(v2)
		}
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.Bytes2)))
	at += 4
	for _, v1 := range bbp.Bytes2 {
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(v1)))
		at += 4
		for _, v2 := range v1 {
			iohelp.WriteUint32Bytes(buf[at:], uint32(len(v2)))
			at += 4
			copy(buf[at:at+len(v2)], v2)
			at += len(v2)
		}
	}
	return at
}

func (bbp *ArraySamples) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.Bytes = make([][][]byte, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.Bytes {
		if len(buf[at:]) < 4 {
			return io.ErrUnexpectedEOF
		}
		(bbp.Bytes)[i1] = make([][]byte, iohelp.ReadUint32Bytes(buf[at:]))
		at += 4
		for i2 := range (bbp.Bytes)[i1] {
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			((bbp.Bytes)[i1])[i2] = make([]byte, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			if len(buf[at:]) < len(((bbp.Bytes)[i1])[i2])*1 {
				return io.ErrUnexpectedEOF
			}
			copy(((bbp.Bytes)[i1])[i2], buf[at:at+len(((bbp.Bytes)[i1])[i2])])
			at += len(((bbp.Bytes)[i1])[i2])
		}
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.Bytes2 = make([][][]byte, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.Bytes2 {
		if len(buf[at:]) < 4 {
			return io.ErrUnexpectedEOF
		}
		(bbp.Bytes2)[i1] = make([][]byte, iohelp.ReadUint32Bytes(buf[at:]))
		at += 4
		for i2 := range (bbp.Bytes2)[i1] {
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			((bbp.Bytes2)[i1])[i2] = make([]byte, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			if len(buf[at:]) < len(((bbp.Bytes2)[i1])[i2])*1 {
				return io.ErrUnexpectedEOF
			}
			copy(((bbp.Bytes2)[i1])[i2], buf[at:at+len(((bbp.Bytes2)[i1])[i2])])
			at += len(((bbp.Bytes2)[i1])[i2])
		}
	}
	return nil
}

func (bbp *ArraySamples) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.Bytes = make([][][]byte, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.Bytes {
		(bbp.Bytes)[i1] = make([][]byte, iohelp.ReadUint32Bytes(buf[at:]))
		at += 4
		for i2 := range (bbp.Bytes)[i1] {
			((bbp.Bytes)[i1])[i2] = make([]byte, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			copy(((bbp.Bytes)[i1])[i2], buf[at:at+len(((bbp.Bytes)[i1])[i2])])
			at += len(((bbp.Bytes)[i1])[i2])
		}
	}
	bbp.Bytes2 = make([][][]byte, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.Bytes2 {
		(bbp.Bytes2)[i1] = make([][]byte, iohelp.ReadUint32Bytes(buf[at:]))
		at += 4
		for i2 := range (bbp.Bytes2)[i1] {
			((bbp.Bytes2)[i1])[i2] = make([]byte, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			copy(((bbp.Bytes2)[i1])[i2], buf[at:at+len(((bbp.Bytes2)[i1])[i2])])
			at += len(((bbp.Bytes2)[i1])[i2])
		}
	}
}

func (bbp ArraySamples) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.Bytes)))
	for _, elem := range bbp.Bytes {
		iohelp.WriteUint32(w, uint32(len(elem)))
		for _, elem := range elem {
			iohelp.WriteUint32(w, uint32(len(elem)))
			w.Write(elem)
		}
	}
	iohelp.WriteUint32(w, uint32(len(bbp.Bytes2)))
	for _, elem := range bbp.Bytes2 {
		iohelp.WriteUint32(w, uint32(len(elem)))
		for _, elem := range elem {
			iohelp.WriteUint32(w, uint32(len(elem)))
			w.Write(elem)
		}
	}
	return w.Err
}

func (bbp *ArraySamples) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.Bytes = make([][][]byte, iohelp.ReadUint32(r))
	for i1 := range bbp.Bytes {
		(bbp.Bytes[i1]) = make([][]byte, iohelp.ReadUint32(r))
		for i2 := range (bbp.Bytes[i1]) {
			((bbp.Bytes[i1])[i2]) = make([]byte, iohelp.ReadUint32(r))
			r.Read(((bbp.Bytes[i1])[i2]))
		}
	}
	bbp.Bytes2 = make([][][]byte, iohelp.ReadUint32(r))
	for i1 := range bbp.Bytes2 {
		(bbp.Bytes2[i1]) = make([][]byte, iohelp.ReadUint32(r))
		for i2 := range (bbp.Bytes2[i1]) {
			((bbp.Bytes2[i1])[i2]) = make([]byte, iohelp.ReadUint32(r))
			r.Read(((bbp.Bytes2[i1])[i2]))
		}
	}
	return r.Err
}

func (bbp ArraySamples) Size() int {
	bodyLen := 0
	bodyLen += 4
	for _, elem := range bbp.Bytes {
		bodyLen += 4
		for _, elem := range elem {
			bodyLen += 4
			bodyLen += len(elem) * 1
		}
	}
	bodyLen += 4
	for _, elem := range bbp.Bytes2 {
		bodyLen += 4
		for _, elem := range elem {
			bodyLen += 4
			bodyLen += len(elem) * 1
		}
	}
	return bodyLen
}

func (bbp ArraySamples) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeArraySamples(r *iohelp.ErrorReader) (ArraySamples, error) {
	v := ArraySamples{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeArraySamplesFromBytes(buf []byte) (ArraySamples, error) {
	v := ArraySamples{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeArraySamplesFromBytes(buf []byte) ArraySamples {
	v := ArraySamples{}
	v.MustUnmarshalBebop(buf)
	return v
}

