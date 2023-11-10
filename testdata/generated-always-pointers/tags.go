// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &TaggedStruct{}

type TaggedStruct struct {
	Foo string `json:"foo,omitempty"`
}

func (bbp *TaggedStruct) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.Foo)))
	copy(buf[at+4:at+4+len(bbp.Foo)], []byte(bbp.Foo))
	at += 4 + len(bbp.Foo)
	return at
}

func (bbp *TaggedStruct) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	bbp.Foo, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil {
		return err
	}
	at += 4 + len(bbp.Foo)
	return nil
}

func (bbp *TaggedStruct) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.Foo = iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.Foo)
}

func (bbp *TaggedStruct) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.Foo)))
	w.Write([]byte(bbp.Foo))
	return w.Err
}

func (bbp *TaggedStruct) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.Foo = iohelp.ReadString(r)
	return r.Err
}

func (bbp *TaggedStruct) Size() int {
	bodyLen := 0
	bodyLen += 4 + len(bbp.Foo)
	return bodyLen
}

func (bbp *TaggedStruct) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeTaggedStruct(r *iohelp.ErrorReader) (TaggedStruct, error) {
	v := TaggedStruct{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeTaggedStructFromBytes(buf []byte) (TaggedStruct, error) {
	v := TaggedStruct{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeTaggedStructFromBytes(buf []byte) TaggedStruct {
	v := TaggedStruct{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &TaggedMessage{}

type TaggedMessage struct {
	Bar *uint8 `db:"bar"`
}

func (bbp *TaggedMessage) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.Bar != nil {
		buf[at] = 1
		at++
		iohelp.WriteUint8Bytes(buf[at:], *bbp.Bar)
		at += 1
	}
	return at
}

func (bbp *TaggedMessage) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Bar = new(uint8)
			if len(buf[at:]) < 1 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.Bar) = iohelp.ReadUint8Bytes(buf[at:])
			at += 1
		default:
			return nil
		}
	}
}

func (bbp *TaggedMessage) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Bar = new(uint8)
			(*bbp.Bar) = iohelp.ReadUint8Bytes(buf[at:])
			at += 1
		default:
			return
		}
	}
}

func (bbp *TaggedMessage) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.Bar != nil {
		w.Write([]byte{1})
		iohelp.WriteUint8(w, *bbp.Bar)
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *TaggedMessage) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Bar = new(uint8)
			*bbp.Bar = iohelp.ReadUint8(r)
		default:
			r.Drain()
			return r.Err
		}
	}
}

func (bbp *TaggedMessage) Size() int {
	bodyLen := 5
	if bbp.Bar != nil {
		bodyLen += 1
		bodyLen += 1
	}
	return bodyLen
}

func (bbp *TaggedMessage) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeTaggedMessage(r *iohelp.ErrorReader) (TaggedMessage, error) {
	v := TaggedMessage{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeTaggedMessageFromBytes(buf []byte) (TaggedMessage, error) {
	v := TaggedMessage{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeTaggedMessageFromBytes(buf []byte) TaggedMessage {
	v := TaggedMessage{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &TaggedSubStruct{}

type TaggedSubStruct struct {
	Biz [16]byte `four:"four"`
}

func (bbp *TaggedSubStruct) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteGUIDBytes(buf[at:], bbp.Biz)
	at += 16
	return at
}

func (bbp *TaggedSubStruct) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 16 {
		return io.ErrUnexpectedEOF
	}
	bbp.Biz = iohelp.ReadGUIDBytes(buf[at:])
	at += 16
	return nil
}

func (bbp *TaggedSubStruct) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.Biz = iohelp.ReadGUIDBytes(buf[at:])
	at += 16
}

func (bbp *TaggedSubStruct) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteGUID(w, bbp.Biz)
	return w.Err
}

func (bbp *TaggedSubStruct) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.Biz = iohelp.ReadGUID(r)
	return r.Err
}

func (bbp *TaggedSubStruct) Size() int {
	bodyLen := 0
	bodyLen += 16
	return bodyLen
}

func (bbp *TaggedSubStruct) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeTaggedSubStruct(r *iohelp.ErrorReader) (TaggedSubStruct, error) {
	v := TaggedSubStruct{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeTaggedSubStructFromBytes(buf []byte) (TaggedSubStruct, error) {
	v := TaggedSubStruct{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeTaggedSubStructFromBytes(buf []byte) TaggedSubStruct {
	v := TaggedSubStruct{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &TaggedUnion{}

type TaggedUnion struct {
	TaggedSubStruct *TaggedSubStruct `one:"one" two:"two" boolean`
}

func (bbp *TaggedUnion) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-5))
	at += 4
	if bbp.TaggedSubStruct != nil {
		buf[at] = 1
		at++
		(*bbp.TaggedSubStruct).MarshalBebopTo(buf[at:])
		tmp := (*bbp.TaggedSubStruct)
		at += tmp.Size()
		return at
	}
	return at
}

func (bbp *TaggedUnion) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	if len(buf) == 0 {
		return iohelp.ErrUnpopulatedUnion
	}
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.TaggedSubStruct = new(TaggedSubStruct)
			(*bbp.TaggedSubStruct), err = MakeTaggedSubStructFromBytes(buf[at:])
			if err != nil {
				return err
			}
			tmp := ((*bbp.TaggedSubStruct))
			at += tmp.Size()
			return nil
		default:
			return nil
		}
	}
}

func (bbp *TaggedUnion) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.TaggedSubStruct = new(TaggedSubStruct)
			(*bbp.TaggedSubStruct) = MustMakeTaggedSubStructFromBytes(buf[at:])
			tmp := ((*bbp.TaggedSubStruct))
			at += tmp.Size()
			return
		default:
			return
		}
	}
}

func (bbp *TaggedUnion) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-5))
	if bbp.TaggedSubStruct != nil {
		w.Write([]byte{1})
		err = (*bbp.TaggedSubStruct).EncodeBebop(w)
		if err != nil {
			return err
		}
		return w.Err
	}
	return w.Err
}

func (bbp *TaggedUnion) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R: r.Reader, N: int64(bodyLen) + 1}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.TaggedSubStruct = new(TaggedSubStruct)
			(*bbp.TaggedSubStruct), err = MakeTaggedSubStruct(r)
			if err != nil {
				return err
			}
			r.Drain()
			return r.Err
		default:
			r.Drain()
			return r.Err
		}
	}
}

func (bbp *TaggedUnion) Size() int {
	bodyLen := 4
	if bbp.TaggedSubStruct != nil {
		bodyLen += 1
		tmp := (*bbp.TaggedSubStruct)
		bodyLen += tmp.Size()
		return bodyLen
	}
	return bodyLen
}

func (bbp *TaggedUnion) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeTaggedUnion(r *iohelp.ErrorReader) (TaggedUnion, error) {
	v := TaggedUnion{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeTaggedUnionFromBytes(buf []byte) (TaggedUnion, error) {
	v := TaggedUnion{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeTaggedUnionFromBytes(buf []byte) TaggedUnion {
	v := TaggedUnion{}
	v.MustUnmarshalBebop(buf)
	return v
}

