// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &s{}

type s struct {
	x int32
	y int32
}

func (bbp s) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteInt32Bytes(buf[at:], bbp.x)
	at += 4
	iohelp.WriteInt32Bytes(buf[at:], bbp.y)
	at += 4
	return at
}

func (bbp *s) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.x = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.y = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	return nil
}

func (bbp *s) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.x = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	bbp.y = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
}

func (bbp s) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteInt32(w, bbp.x)
	iohelp.WriteInt32(w, bbp.y)
	return w.Err
}

func (bbp *s) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.x = iohelp.ReadInt32(r)
	bbp.y = iohelp.ReadInt32(r)
	return r.Err
}

func (bbp s) Size() int {
	bodyLen := 0
	bodyLen += 4
	bodyLen += 4
	return bodyLen
}

func (bbp s) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makes(r iohelp.ErrorReader) (s, error) {
	v := s{}
	err := v.DecodeBebop(r)
	return v, err
}

func makesFromBytes(buf []byte) (s, error) {
	v := s{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakesFromBytes(buf []byte) s {
	v := s{}
	v.MustUnmarshalBebop(buf)
	return v
}

func (bbp s) Getx() int32 {
	return bbp.x
}

func (bbp s) Gety() int32 {
	return bbp.y
}

func news(
		x int32,
		y int32,
) s {
	return s{
		x: x,
		y: y,
	}
}

var _ bebop.Record = &someMaps{}

type someMaps struct {
	m1 map[bool]bool
	m2 map[string]map[string]string
	m3 []map[int32][]map[bool]s
	m4 []map[string][]float32
	m5 map[[16]byte]m
}

func (bbp someMaps) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.m1)))
	at += 4
	for k1, v1 := range bbp.m1 {
		iohelp.WriteBoolBytes(buf[at:], k1)
		at += 1
		iohelp.WriteBoolBytes(buf[at:], v1)
		at += 1
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.m2)))
	at += 4
	for k1, v1 := range bbp.m2 {
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(k1)))
		copy(buf[at+4:at+4+len(k1)], []byte(k1))
		at += 4 + len(k1)
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(v1)))
		at += 4
		for k2, v2 := range v1 {
			iohelp.WriteUint32Bytes(buf[at:], uint32(len(k2)))
			copy(buf[at+4:at+4+len(k2)], []byte(k2))
			at += 4 + len(k2)
			iohelp.WriteUint32Bytes(buf[at:], uint32(len(v2)))
			copy(buf[at+4:at+4+len(v2)], []byte(v2))
			at += 4 + len(v2)
		}
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.m3)))
	at += 4
	for _, v1 := range bbp.m3 {
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(v1)))
		at += 4
		for k2, v2 := range v1 {
			iohelp.WriteInt32Bytes(buf[at:], k2)
			at += 4
			iohelp.WriteUint32Bytes(buf[at:], uint32(len(v2)))
			at += 4
			for _, v3 := range v2 {
				iohelp.WriteUint32Bytes(buf[at:], uint32(len(v3)))
				at += 4
				for k4, v4 := range v3 {
					iohelp.WriteBoolBytes(buf[at:], k4)
					at += 1
					(v4).MarshalBebopTo(buf[at:])
					tmp7051 := (v4); at += tmp7051.Size()
				}
			}
		}
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.m4)))
	at += 4
	for _, v1 := range bbp.m4 {
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(v1)))
		at += 4
		for k2, v2 := range v1 {
			iohelp.WriteUint32Bytes(buf[at:], uint32(len(k2)))
			copy(buf[at+4:at+4+len(k2)], []byte(k2))
			at += 4 + len(k2)
			iohelp.WriteUint32Bytes(buf[at:], uint32(len(v2)))
			at += 4
			for _, v3 := range v2 {
				iohelp.WriteFloat32Bytes(buf[at:], v3)
				at += 4
			}
		}
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.m5)))
	at += 4
	for k1, v1 := range bbp.m5 {
		iohelp.WriteGUIDBytes(buf[at:], k1)
		at += 16
		(v1).MarshalBebopTo(buf[at:])
		tmp7097 := (v1); at += tmp7097.Size()
	}
	return at
}

func (bbp *someMaps) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	ln1 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.m1 = make(map[bool]bool,ln1)
	for i := uint32(0); i < ln1; i++ {
		if len(buf[at:]) < 1 {
			return io.ErrUnexpectedEOF
		}
		k1 := iohelp.ReadBoolBytes(buf[at:])
		at += 1
		if len(buf[at:]) < 1 {
			return io.ErrUnexpectedEOF
		}
		(bbp.m1)[k1] = iohelp.ReadBoolBytes(buf[at:])
		at += 1
	}
	ln2 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.m2 = make(map[string]map[string]string,ln2)
	for i := uint32(0); i < ln2; i++ {
		k1, err := iohelp.ReadStringBytes(buf[at:])
		if err != nil {
			return err
		}
		at += 4 + len(k1)
		ln3 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.m2)[k1] = make(map[string]string,ln3)
		for i := uint32(0); i < ln3; i++ {
			k2, err := iohelp.ReadStringBytes(buf[at:])
			if err != nil {
				return err
			}
			at += 4 + len(k2)
			((bbp.m2)[k1])[k2], err = iohelp.ReadStringBytes(buf[at:])
			if err != nil {
				return err
			}
			at += 4 + len(((bbp.m2)[k1])[k2])
		}
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.m3 = make([]map[int32][]map[bool]s, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.m3 {
		ln4 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.m3)[i1] = make(map[int32][]map[bool]s,ln4)
		for i := uint32(0); i < ln4; i++ {
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			k2 := iohelp.ReadInt32Bytes(buf[at:])
			at += 4
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			((bbp.m3)[i1])[k2] = make([]map[bool]s, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			for i3 := range ((bbp.m3)[i1])[k2] {
				ln5 := iohelp.ReadUint32Bytes(buf[at:])
				at += 4
				(((bbp.m3)[i1])[k2])[i3] = make(map[bool]s,ln5)
				for i := uint32(0); i < ln5; i++ {
					if len(buf[at:]) < 1 {
						return io.ErrUnexpectedEOF
					}
					k4 := iohelp.ReadBoolBytes(buf[at:])
					at += 1
					((((bbp.m3)[i1])[k2])[i3])[k4], err = makesFromBytes(buf[at:])
					if err != nil {
						return err
					}
					tmp7194 := (((((bbp.m3)[i1])[k2])[i3])[k4]); at += tmp7194.Size()
				}
			}
		}
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.m4 = make([]map[string][]float32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.m4 {
		ln6 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.m4)[i1] = make(map[string][]float32,ln6)
		for i := uint32(0); i < ln6; i++ {
			k2, err := iohelp.ReadStringBytes(buf[at:])
			if err != nil {
				return err
			}
			at += 4 + len(k2)
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			((bbp.m4)[i1])[k2] = make([]float32, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			if len(buf[at:]) < len(((bbp.m4)[i1])[k2])*4 {
				return io.ErrUnexpectedEOF
			}
			for i3 := range ((bbp.m4)[i1])[k2] {
				(((bbp.m4)[i1])[k2])[i3] = iohelp.ReadFloat32Bytes(buf[at:])
				at += 4
			}
		}
	}
	ln7 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.m5 = make(map[[16]byte]m,ln7)
	for i := uint32(0); i < ln7; i++ {
		if len(buf[at:]) < 16 {
			return io.ErrUnexpectedEOF
		}
		k1 := iohelp.ReadGUIDBytes(buf[at:])
		at += 16
		(bbp.m5)[k1], err = makemFromBytes(buf[at:])
		if err != nil {
			return err
		}
		tmp7248 := ((bbp.m5)[k1]); at += tmp7248.Size()
	}
	return nil
}

func (bbp *someMaps) MustUnmarshalBebop(buf []byte) {
	at := 0
	ln8 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.m1 = make(map[bool]bool,ln8)
	for i := uint32(0); i < ln8; i++ {
		k1 := iohelp.ReadBoolBytes(buf[at:])
		at += 1
		(bbp.m1)[k1] = iohelp.ReadBoolBytes(buf[at:])
		at += 1
	}
	ln9 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.m2 = make(map[string]map[string]string,ln9)
	for i := uint32(0); i < ln9; i++ {
		k1 := iohelp.MustReadStringBytes(buf[at:])
		at += 4 + len(k1)
		ln10 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.m2)[k1] = make(map[string]string,ln10)
		for i := uint32(0); i < ln10; i++ {
			k2 := iohelp.MustReadStringBytes(buf[at:])
			at += 4 + len(k2)
			((bbp.m2)[k1])[k2] = iohelp.MustReadStringBytes(buf[at:])
			at += 4 + len(((bbp.m2)[k1])[k2])
		}
	}
	bbp.m3 = make([]map[int32][]map[bool]s, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.m3 {
		ln11 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.m3)[i1] = make(map[int32][]map[bool]s,ln11)
		for i := uint32(0); i < ln11; i++ {
			k2 := iohelp.ReadInt32Bytes(buf[at:])
			at += 4
			((bbp.m3)[i1])[k2] = make([]map[bool]s, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			for i3 := range ((bbp.m3)[i1])[k2] {
				ln12 := iohelp.ReadUint32Bytes(buf[at:])
				at += 4
				(((bbp.m3)[i1])[k2])[i3] = make(map[bool]s,ln12)
				for i := uint32(0); i < ln12; i++ {
					k4 := iohelp.ReadBoolBytes(buf[at:])
					at += 1
					((((bbp.m3)[i1])[k2])[i3])[k4] = mustMakesFromBytes(buf[at:])
					tmp7295 := (((((bbp.m3)[i1])[k2])[i3])[k4]); at += tmp7295.Size()
				}
			}
		}
	}
	bbp.m4 = make([]map[string][]float32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.m4 {
		ln13 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.m4)[i1] = make(map[string][]float32,ln13)
		for i := uint32(0); i < ln13; i++ {
			k2 := iohelp.MustReadStringBytes(buf[at:])
			at += 4 + len(k2)
			((bbp.m4)[i1])[k2] = make([]float32, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			for i3 := range ((bbp.m4)[i1])[k2] {
				(((bbp.m4)[i1])[k2])[i3] = iohelp.ReadFloat32Bytes(buf[at:])
				at += 4
			}
		}
	}
	ln14 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.m5 = make(map[[16]byte]m,ln14)
	for i := uint32(0); i < ln14; i++ {
		k1 := iohelp.ReadGUIDBytes(buf[at:])
		at += 16
		(bbp.m5)[k1] = mustMakemFromBytes(buf[at:])
		tmp7326 := ((bbp.m5)[k1]); at += tmp7326.Size()
	}
}

func (bbp someMaps) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.m1)))
	for k1, v1 := range bbp.m1 {
		iohelp.WriteBool(w, k1)
		iohelp.WriteBool(w, v1)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.m2)))
	for k1, v1 := range bbp.m2 {
		iohelp.WriteUint32(w, uint32(len(k1)))
		w.Write([]byte(k1))
		iohelp.WriteUint32(w, uint32(len(v1)))
		for k2, v2 := range v1 {
			iohelp.WriteUint32(w, uint32(len(k2)))
			w.Write([]byte(k2))
			iohelp.WriteUint32(w, uint32(len(v2)))
			w.Write([]byte(v2))
		}
	}
	iohelp.WriteUint32(w, uint32(len(bbp.m3)))
	for _, elem := range bbp.m3 {
		iohelp.WriteUint32(w, uint32(len(elem)))
		for k2, v2 := range elem {
			iohelp.WriteInt32(w, k2)
			iohelp.WriteUint32(w, uint32(len(v2)))
			for _, elem := range v2 {
				iohelp.WriteUint32(w, uint32(len(elem)))
				for k4, v4 := range elem {
					iohelp.WriteBool(w, k4)
					err = (v4).EncodeBebop(w)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	iohelp.WriteUint32(w, uint32(len(bbp.m4)))
	for _, elem := range bbp.m4 {
		iohelp.WriteUint32(w, uint32(len(elem)))
		for k2, v2 := range elem {
			iohelp.WriteUint32(w, uint32(len(k2)))
			w.Write([]byte(k2))
			iohelp.WriteUint32(w, uint32(len(v2)))
			for _, elem := range v2 {
				iohelp.WriteFloat32(w, elem)
			}
		}
	}
	iohelp.WriteUint32(w, uint32(len(bbp.m5)))
	for k1, v1 := range bbp.m5 {
		iohelp.WriteGUID(w, k1)
		err = (v1).EncodeBebop(w)
		if err != nil {
			return err
		}
	}
	return w.Err
}

func (bbp *someMaps) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	ln1 := iohelp.ReadUint32(r)
	bbp.m1 = make(map[bool]bool, ln1)
	for i1 := uint32(0); i1 < ln1; i1++ {
		k1 := iohelp.ReadBool(r)
		(bbp.m1[k1]) = iohelp.ReadBool(r)
	}
	ln1 = iohelp.ReadUint32(r)
	bbp.m2 = make(map[string]map[string]string, ln1)
	for i1 := uint32(0); i1 < ln1; i1++ {
		k1 := iohelp.ReadString(r)
		ln2 := iohelp.ReadUint32(r)
		(bbp.m2[k1]) = make(map[string]string, ln2)
		for i2 := uint32(0); i2 < ln2; i2++ {
			k2 := iohelp.ReadString(r)
			((bbp.m2[k1])[k2]) = iohelp.ReadString(r)
		}
	}
	bbp.m3 = make([]map[int32][]map[bool]s, iohelp.ReadUint32(r))
	for i1 := range bbp.m3 {
		ln2 := iohelp.ReadUint32(r)
		(bbp.m3[i1]) = make(map[int32][]map[bool]s, ln2)
		for i2 := uint32(0); i2 < ln2; i2++ {
			k2 := iohelp.ReadInt32(r)
			((bbp.m3[i1])[k2]) = make([]map[bool]s, iohelp.ReadUint32(r))
			for i3 := range ((bbp.m3[i1])[k2]) {
				ln4 := iohelp.ReadUint32(r)
				(((bbp.m3[i1])[k2])[i3]) = make(map[bool]s, ln4)
				for i4 := uint32(0); i4 < ln4; i4++ {
					k4 := iohelp.ReadBool(r)
					(((((bbp.m3[i1])[k2])[i3])[k4])), err = makes(r)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	bbp.m4 = make([]map[string][]float32, iohelp.ReadUint32(r))
	for i1 := range bbp.m4 {
		ln2 := iohelp.ReadUint32(r)
		(bbp.m4[i1]) = make(map[string][]float32, ln2)
		for i2 := uint32(0); i2 < ln2; i2++ {
			k2 := iohelp.ReadString(r)
			((bbp.m4[i1])[k2]) = make([]float32, iohelp.ReadUint32(r))
			for i3 := range ((bbp.m4[i1])[k2]) {
				(((bbp.m4[i1])[k2])[i3]) = iohelp.ReadFloat32(r)
			}
		}
	}
	ln1 = iohelp.ReadUint32(r)
	bbp.m5 = make(map[[16]byte]m, ln1)
	for i1 := uint32(0); i1 < ln1; i1++ {
		k1 := iohelp.ReadGUID(r)
		((bbp.m5[k1])), err = makem(r)
		if err != nil {
			return err
		}
	}
	return r.Err
}

func (bbp someMaps) Size() int {
	bodyLen := 0
	bodyLen += 4
	for range bbp.m1 {
		bodyLen += 1
		bodyLen += 1
	}
	bodyLen += 4
	for k1, v1 := range bbp.m2 {
		bodyLen += 4 + len(k1)
		bodyLen += 4
		for k2, v2 := range v1 {
			bodyLen += 4 + len(k2)
			bodyLen += 4 + len(v2)
		}
	}
	bodyLen += 4
	for _, elem := range bbp.m3 {
		bodyLen += 4
		for _, v2 := range elem {
			bodyLen += 4
			bodyLen += 4
			for _, elem := range v2 {
				bodyLen += 4
				for _, v4 := range elem {
					bodyLen += 1
					tmp8614 := (v4); bodyLen += tmp8614.Size()
				}
			}
		}
	}
	bodyLen += 4
	for _, elem := range bbp.m4 {
		bodyLen += 4
		for k2, v2 := range elem {
			bodyLen += 4 + len(k2)
			bodyLen += 4
			bodyLen += len(v2) * 4
		}
	}
	bodyLen += 4
	for _, v1 := range bbp.m5 {
		bodyLen += 16
		tmp8644 := (v1); bodyLen += tmp8644.Size()
	}
	return bodyLen
}

func (bbp someMaps) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makesomeMaps(r iohelp.ErrorReader) (someMaps, error) {
	v := someMaps{}
	err := v.DecodeBebop(r)
	return v, err
}

func makesomeMapsFromBytes(buf []byte) (someMaps, error) {
	v := someMaps{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakesomeMapsFromBytes(buf []byte) someMaps {
	v := someMaps{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &m{}

type m struct {
	a *float32
	b *float64
}

func (bbp m) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.a != nil {
		buf[at] = 1
		at++
		iohelp.WriteFloat32Bytes(buf[at:], *bbp.a)
		at += 4
	}
	if bbp.b != nil {
		buf[at] = 2
		at++
		iohelp.WriteFloat64Bytes(buf[at:], *bbp.b)
		at += 8
	}
	return at
}

func (bbp *m) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.a = new(float32)
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.a) = iohelp.ReadFloat32Bytes(buf[at:])
			at += 4
		case 2:
			at += 1
			bbp.b = new(float64)
			if len(buf[at:]) < 8 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.b) = iohelp.ReadFloat64Bytes(buf[at:])
			at += 8
		default:
			return nil
		}
	}
}

func (bbp *m) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.a = new(float32)
			(*bbp.a) = iohelp.ReadFloat32Bytes(buf[at:])
			at += 4
		case 2:
			at += 1
			bbp.b = new(float64)
			(*bbp.b) = iohelp.ReadFloat64Bytes(buf[at:])
			at += 8
		default:
			return
		}
	}
}

func (bbp m) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.a != nil {
		w.Write([]byte{1})
		iohelp.WriteFloat32(w, *bbp.a)
	}
	if bbp.b != nil {
		w.Write([]byte{2})
		iohelp.WriteFloat64(w, *bbp.b)
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *m) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.a = new(float32)
			*bbp.a = iohelp.ReadFloat32(r)
		case 2:
			bbp.b = new(float64)
			*bbp.b = iohelp.ReadFloat64(r)
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp m) Size() int {
	bodyLen := 5
	if bbp.a != nil {
		bodyLen += 1
		bodyLen += 4
	}
	if bbp.b != nil {
		bodyLen += 1
		bodyLen += 8
	}
	return bodyLen
}

func (bbp m) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makem(r iohelp.ErrorReader) (m, error) {
	v := m{}
	err := v.DecodeBebop(r)
	return v, err
}

func makemFromBytes(buf []byte) (m, error) {
	v := m{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakemFromBytes(buf []byte) m {
	v := m{}
	v.MustUnmarshalBebop(buf)
	return v
}

