// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &S{}

type S struct {
	x int32
	y int32
}

func (bbp *S) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteInt32Bytes(buf[at:], bbp.x)
	at += 4
	iohelp.WriteInt32Bytes(buf[at:], bbp.y)
	at += 4
	return at
}

func (bbp *S) UnmarshalBebop(buf []byte) (err error) {
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

func (bbp *S) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.x = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
	bbp.y = iohelp.ReadInt32Bytes(buf[at:])
	at += 4
}

func (bbp *S) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteInt32(w, bbp.x)
	iohelp.WriteInt32(w, bbp.y)
	return w.Err
}

func (bbp *S) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.x = iohelp.ReadInt32(r)
	bbp.y = iohelp.ReadInt32(r)
	return r.Err
}

func (bbp *S) Size() int {
	bodyLen := 0
	bodyLen += 4
	bodyLen += 4
	return bodyLen
}

func (bbp *S) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeS(r *iohelp.ErrorReader) (S, error) {
	v := S{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeSFromBytes(buf []byte) (S, error) {
	v := S{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeSFromBytes(buf []byte) S {
	v := S{}
	v.MustUnmarshalBebop(buf)
	return v
}

func (bbp *S) GetX() int32 {
	return bbp.x
}

func (bbp *S) GetY() int32 {
	return bbp.y
}

func NewS(
		x int32,
		y int32,
) S {
	return S{
		x: x,
		y: y,
	}
}

var _ bebop.Record = &SomeMaps{}

type SomeMaps struct {
	M1 map[bool]bool
	M2 map[string]map[string]string
	M3 []map[int32][]map[bool]S
	M4 []map[string][]float32
	M5 map[[16]byte]M
}

func (bbp *SomeMaps) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.M1)))
	at += 4
	for k1, v1 := range bbp.M1 {
		iohelp.WriteBoolBytes(buf[at:], k1)
		at += 1
		iohelp.WriteBoolBytes(buf[at:], v1)
		at += 1
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.M2)))
	at += 4
	for k1, v1 := range bbp.M2 {
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
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.M3)))
	at += 4
	for _, v1 := range bbp.M3 {
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
					{
						tmp := (v4)
						at += tmp.Size()
					}
					
				}
			}
		}
	}
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.M4)))
	at += 4
	for _, v1 := range bbp.M4 {
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
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.M5)))
	at += 4
	for k1, v1 := range bbp.M5 {
		iohelp.WriteGUIDBytes(buf[at:], k1)
		at += 16
		(v1).MarshalBebopTo(buf[at:])
		{
			tmp := (v1)
			at += tmp.Size()
		}
		
	}
	return at
}

func (bbp *SomeMaps) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	ln1 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.M1 = make(map[bool]bool,ln1)
	for i := uint32(0); i < ln1; i++ {
		if len(buf[at:]) < 1 {
			return io.ErrUnexpectedEOF
		}
		k1 := iohelp.ReadBoolBytes(buf[at:])
		at += 1
		if len(buf[at:]) < 1 {
			return io.ErrUnexpectedEOF
		}
		(bbp.M1)[k1] = iohelp.ReadBoolBytes(buf[at:])
		at += 1
	}
	ln2 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.M2 = make(map[string]map[string]string,ln2)
	for i := uint32(0); i < ln2; i++ {
		k1, err := iohelp.ReadStringBytes(buf[at:])
		if err != nil {
			return err
		}
		at += 4 + len(k1)
		ln3 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.M2)[k1] = make(map[string]string,ln3)
		for i := uint32(0); i < ln3; i++ {
			k2, err := iohelp.ReadStringBytes(buf[at:])
			if err != nil {
				return err
			}
			at += 4 + len(k2)
			((bbp.M2)[k1])[k2], err = iohelp.ReadStringBytes(buf[at:])
			if err != nil {
				return err
			}
			at += 4 + len(((bbp.M2)[k1])[k2])
		}
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.M3 = make([]map[int32][]map[bool]S, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.M3 {
		ln4 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.M3)[i1] = make(map[int32][]map[bool]S,ln4)
		for i := uint32(0); i < ln4; i++ {
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			k2 := iohelp.ReadInt32Bytes(buf[at:])
			at += 4
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			((bbp.M3)[i1])[k2] = make([]map[bool]S, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			for i3 := range ((bbp.M3)[i1])[k2] {
				ln5 := iohelp.ReadUint32Bytes(buf[at:])
				at += 4
				(((bbp.M3)[i1])[k2])[i3] = make(map[bool]S,ln5)
				for i := uint32(0); i < ln5; i++ {
					if len(buf[at:]) < 1 {
						return io.ErrUnexpectedEOF
					}
					k4 := iohelp.ReadBoolBytes(buf[at:])
					at += 1
					((((bbp.M3)[i1])[k2])[i3])[k4], err = MakeSFromBytes(buf[at:])
					if err != nil {
						return err
					}
					{
						tmp := (((((bbp.M3)[i1])[k2])[i3])[k4])
						at += tmp.Size()
					}
					
				}
			}
		}
	}
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.M4 = make([]map[string][]float32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.M4 {
		ln6 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.M4)[i1] = make(map[string][]float32,ln6)
		for i := uint32(0); i < ln6; i++ {
			k2, err := iohelp.ReadStringBytes(buf[at:])
			if err != nil {
				return err
			}
			at += 4 + len(k2)
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			((bbp.M4)[i1])[k2] = make([]float32, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			if len(buf[at:]) < len(((bbp.M4)[i1])[k2])*4 {
				return io.ErrUnexpectedEOF
			}
			for i3 := range ((bbp.M4)[i1])[k2] {
				(((bbp.M4)[i1])[k2])[i3] = iohelp.ReadFloat32Bytes(buf[at:])
				at += 4
			}
		}
	}
	ln7 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.M5 = make(map[[16]byte]M,ln7)
	for i := uint32(0); i < ln7; i++ {
		if len(buf[at:]) < 16 {
			return io.ErrUnexpectedEOF
		}
		k1 := iohelp.ReadGUIDBytes(buf[at:])
		at += 16
		(bbp.M5)[k1], err = MakeMFromBytes(buf[at:])
		if err != nil {
			return err
		}
		{
			tmp := ((bbp.M5)[k1])
			at += tmp.Size()
		}
		
	}
	return nil
}

func (bbp *SomeMaps) MustUnmarshalBebop(buf []byte) {
	at := 0
	ln8 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.M1 = make(map[bool]bool,ln8)
	for i := uint32(0); i < ln8; i++ {
		k1 := iohelp.ReadBoolBytes(buf[at:])
		at += 1
		(bbp.M1)[k1] = iohelp.ReadBoolBytes(buf[at:])
		at += 1
	}
	ln9 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.M2 = make(map[string]map[string]string,ln9)
	for i := uint32(0); i < ln9; i++ {
		k1 := iohelp.MustReadStringBytes(buf[at:])
		at += 4 + len(k1)
		ln10 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.M2)[k1] = make(map[string]string,ln10)
		for i := uint32(0); i < ln10; i++ {
			k2 := iohelp.MustReadStringBytes(buf[at:])
			at += 4 + len(k2)
			((bbp.M2)[k1])[k2] = iohelp.MustReadStringBytes(buf[at:])
			at += 4 + len(((bbp.M2)[k1])[k2])
		}
	}
	bbp.M3 = make([]map[int32][]map[bool]S, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.M3 {
		ln11 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.M3)[i1] = make(map[int32][]map[bool]S,ln11)
		for i := uint32(0); i < ln11; i++ {
			k2 := iohelp.ReadInt32Bytes(buf[at:])
			at += 4
			((bbp.M3)[i1])[k2] = make([]map[bool]S, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			for i3 := range ((bbp.M3)[i1])[k2] {
				ln12 := iohelp.ReadUint32Bytes(buf[at:])
				at += 4
				(((bbp.M3)[i1])[k2])[i3] = make(map[bool]S,ln12)
				for i := uint32(0); i < ln12; i++ {
					k4 := iohelp.ReadBoolBytes(buf[at:])
					at += 1
					((((bbp.M3)[i1])[k2])[i3])[k4] = MustMakeSFromBytes(buf[at:])
					{
						tmp := (((((bbp.M3)[i1])[k2])[i3])[k4])
						at += tmp.Size()
					}
					
				}
			}
		}
	}
	bbp.M4 = make([]map[string][]float32, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.M4 {
		ln13 := iohelp.ReadUint32Bytes(buf[at:])
		at += 4
		(bbp.M4)[i1] = make(map[string][]float32,ln13)
		for i := uint32(0); i < ln13; i++ {
			k2 := iohelp.MustReadStringBytes(buf[at:])
			at += 4 + len(k2)
			((bbp.M4)[i1])[k2] = make([]float32, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			for i3 := range ((bbp.M4)[i1])[k2] {
				(((bbp.M4)[i1])[k2])[i3] = iohelp.ReadFloat32Bytes(buf[at:])
				at += 4
			}
		}
	}
	ln14 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.M5 = make(map[[16]byte]M,ln14)
	for i := uint32(0); i < ln14; i++ {
		k1 := iohelp.ReadGUIDBytes(buf[at:])
		at += 16
		(bbp.M5)[k1] = MustMakeMFromBytes(buf[at:])
		{
			tmp := ((bbp.M5)[k1])
			at += tmp.Size()
		}
		
	}
}

func (bbp *SomeMaps) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.M1)))
	for k1, v1 := range bbp.M1 {
		iohelp.WriteBool(w, k1)
		iohelp.WriteBool(w, v1)
	}
	iohelp.WriteUint32(w, uint32(len(bbp.M2)))
	for k1, v1 := range bbp.M2 {
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
	iohelp.WriteUint32(w, uint32(len(bbp.M3)))
	for _, elem := range bbp.M3 {
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
	iohelp.WriteUint32(w, uint32(len(bbp.M4)))
	for _, elem := range bbp.M4 {
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
	iohelp.WriteUint32(w, uint32(len(bbp.M5)))
	for k1, v1 := range bbp.M5 {
		iohelp.WriteGUID(w, k1)
		err = (v1).EncodeBebop(w)
		if err != nil {
			return err
		}
	}
	return w.Err
}

func (bbp *SomeMaps) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	ln1 := iohelp.ReadUint32(r)
	bbp.M1 = make(map[bool]bool, ln1)
	for i1 := uint32(0); i1 < ln1; i1++ {
		k1 := iohelp.ReadBool(r)
		(bbp.M1[k1]) = iohelp.ReadBool(r)
	}
	ln1 = iohelp.ReadUint32(r)
	bbp.M2 = make(map[string]map[string]string, ln1)
	for i1 := uint32(0); i1 < ln1; i1++ {
		k1 := iohelp.ReadString(r)
		ln2 := iohelp.ReadUint32(r)
		(bbp.M2[k1]) = make(map[string]string, ln2)
		for i2 := uint32(0); i2 < ln2; i2++ {
			k2 := iohelp.ReadString(r)
			((bbp.M2[k1])[k2]) = iohelp.ReadString(r)
		}
	}
	bbp.M3 = make([]map[int32][]map[bool]S, iohelp.ReadUint32(r))
	for i1 := range bbp.M3 {
		ln2 := iohelp.ReadUint32(r)
		(bbp.M3[i1]) = make(map[int32][]map[bool]S, ln2)
		for i2 := uint32(0); i2 < ln2; i2++ {
			k2 := iohelp.ReadInt32(r)
			((bbp.M3[i1])[k2]) = make([]map[bool]S, iohelp.ReadUint32(r))
			for i3 := range ((bbp.M3[i1])[k2]) {
				ln4 := iohelp.ReadUint32(r)
				(((bbp.M3[i1])[k2])[i3]) = make(map[bool]S, ln4)
				for i4 := uint32(0); i4 < ln4; i4++ {
					k4 := iohelp.ReadBool(r)
					(((((bbp.M3[i1])[k2])[i3])[k4])), err = MakeS(r)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	bbp.M4 = make([]map[string][]float32, iohelp.ReadUint32(r))
	for i1 := range bbp.M4 {
		ln2 := iohelp.ReadUint32(r)
		(bbp.M4[i1]) = make(map[string][]float32, ln2)
		for i2 := uint32(0); i2 < ln2; i2++ {
			k2 := iohelp.ReadString(r)
			((bbp.M4[i1])[k2]) = make([]float32, iohelp.ReadUint32(r))
			for i3 := range ((bbp.M4[i1])[k2]) {
				(((bbp.M4[i1])[k2])[i3]) = iohelp.ReadFloat32(r)
			}
		}
	}
	ln1 = iohelp.ReadUint32(r)
	bbp.M5 = make(map[[16]byte]M, ln1)
	for i1 := uint32(0); i1 < ln1; i1++ {
		k1 := iohelp.ReadGUID(r)
		((bbp.M5[k1])), err = MakeM(r)
		if err != nil {
			return err
		}
	}
	return r.Err
}

func (bbp *SomeMaps) Size() int {
	bodyLen := 0
	bodyLen += 4
	for range bbp.M1 {
		bodyLen += 1
		bodyLen += 1
	}
	bodyLen += 4
	for k1, v1 := range bbp.M2 {
		bodyLen += 4 + len(k1)
		bodyLen += 4
		for k2, v2 := range v1 {
			bodyLen += 4 + len(k2)
			bodyLen += 4 + len(v2)
		}
	}
	bodyLen += 4
	for _, elem := range bbp.M3 {
		bodyLen += 4
		for _, v2 := range elem {
			bodyLen += 4
			bodyLen += 4
			for _, elem := range v2 {
				bodyLen += 4
				for _, v4 := range elem {
					bodyLen += 1
					{
						tmp := (v4)
						bodyLen += tmp.Size()
					}
					
				}
			}
		}
	}
	bodyLen += 4
	for _, elem := range bbp.M4 {
		bodyLen += 4
		for k2, v2 := range elem {
			bodyLen += 4 + len(k2)
			bodyLen += 4
			bodyLen += len(v2) * 4
		}
	}
	bodyLen += 4
	for _, v1 := range bbp.M5 {
		bodyLen += 16
		{
			tmp := (v1)
			bodyLen += tmp.Size()
		}
		
	}
	return bodyLen
}

func (bbp *SomeMaps) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeSomeMaps(r *iohelp.ErrorReader) (SomeMaps, error) {
	v := SomeMaps{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeSomeMapsFromBytes(buf []byte) (SomeMaps, error) {
	v := SomeMaps{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeSomeMapsFromBytes(buf []byte) SomeMaps {
	v := SomeMaps{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &M{}

type M struct {
	A *float32
	B *float64
}

func (bbp *M) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.A != nil {
		buf[at] = 1
		at++
		iohelp.WriteFloat32Bytes(buf[at:], *bbp.A)
		at += 4
	}
	if bbp.B != nil {
		buf[at] = 2
		at++
		iohelp.WriteFloat64Bytes(buf[at:], *bbp.B)
		at += 8
	}
	return at
}

func (bbp *M) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.A = new(float32)
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.A) = iohelp.ReadFloat32Bytes(buf[at:])
			at += 4
		case 2:
			at += 1
			bbp.B = new(float64)
			if len(buf[at:]) < 8 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.B) = iohelp.ReadFloat64Bytes(buf[at:])
			at += 8
		default:
			return nil
		}
	}
}

func (bbp *M) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.A = new(float32)
			(*bbp.A) = iohelp.ReadFloat32Bytes(buf[at:])
			at += 4
		case 2:
			at += 1
			bbp.B = new(float64)
			(*bbp.B) = iohelp.ReadFloat64Bytes(buf[at:])
			at += 8
		default:
			return
		}
	}
}

func (bbp *M) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.A != nil {
		w.Write([]byte{1})
		iohelp.WriteFloat32(w, *bbp.A)
	}
	if bbp.B != nil {
		w.Write([]byte{2})
		iohelp.WriteFloat64(w, *bbp.B)
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *M) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.A = new(float32)
			*bbp.A = iohelp.ReadFloat32(r)
		case 2:
			bbp.B = new(float64)
			*bbp.B = iohelp.ReadFloat64(r)
		default:
			r.Drain()
			return r.Err
		}
	}
}

func (bbp *M) Size() int {
	bodyLen := 5
	if bbp.A != nil {
		bodyLen += 1
		bodyLen += 4
	}
	if bbp.B != nil {
		bodyLen += 1
		bodyLen += 8
	}
	return bodyLen
}

func (bbp *M) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeM(r *iohelp.ErrorReader) (M, error) {
	v := M{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeMFromBytes(buf []byte) (M, error) {
	v := M{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeMFromBytes(buf []byte) M {
	v := M{}
	v.MustUnmarshalBebop(buf)
	return v
}

