// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

type furnitureFamily uint32

const (
	furnitureFamily_Bed furnitureFamily = 0
	furnitureFamily_Table furnitureFamily = 1
	furnitureFamily_Shoe furnitureFamily = 2
)

var _ bebop.Record = &furniture{}

type furniture struct {
	name string
	price uint32
	family furnitureFamily
}

func (bbp *furniture) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.name)))
	copy(buf[at+4:at+4+len(bbp.name)], []byte(bbp.name))
	at += 4 + len(bbp.name)
	iohelp.WriteUint32Bytes(buf[at:], bbp.price)
	at += 4
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.family))
	at += 4
	return at
}

func (bbp *furniture) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	bbp.name, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil {
		return err
	}
	at += 4 + len(bbp.name)
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.price = iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.family = furnitureFamily(iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	return nil
}

func (bbp *furniture) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.name = iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.name)
	bbp.price = iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.family = furnitureFamily(iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
}

func (bbp *furniture) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.name)))
	w.Write([]byte(bbp.name))
	iohelp.WriteUint32(w, bbp.price)
	iohelp.WriteUint32(w, uint32(bbp.family))
	return w.Err
}

func (bbp *furniture) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.name = iohelp.ReadString(r)
	bbp.price = iohelp.ReadUint32(r)
	bbp.family = furnitureFamily(iohelp.ReadUint32(r))
	return r.Err
}

func (bbp *furniture) Size() int {
	bodyLen := 0
	bodyLen += 4 + len(bbp.name)
	bodyLen += 4
	bodyLen += 4
	return bodyLen
}

func (bbp *furniture) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makefurniture(r iohelp.ErrorReader) (furniture, error) {
	v := furniture{}
	err := v.DecodeBebop(r)
	return v, err
}

func makefurnitureFromBytes(buf []byte) (furniture, error) {
	v := furniture{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakefurnitureFromBytes(buf []byte) furniture {
	v := furniture{}
	v.MustUnmarshalBebop(buf)
	return v
}

func (bbp *furniture) Getname() string {
	return bbp.name
}

func (bbp *furniture) Getprice() uint32 {
	return bbp.price
}

func (bbp *furniture) Getfamily() furnitureFamily {
	return bbp.family
}

func newfurniture(
		name string,
		price uint32,
		family furnitureFamily,
) furniture {
	return furniture{
		name: name,
		price: price,
		family: family,
	}
}

const requestResponseOpCode = 0x31323334

var _ bebop.Record = &requestResponse{}

type requestResponse struct {
	availableFurniture []furniture
}

func (bbp *requestResponse) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.availableFurniture)))
	at += 4
	for _, v1 := range bbp.availableFurniture {
		(v1).MarshalBebopTo(buf[at:])
		at += (v1).Size()
	}
	return at
}

func (bbp *requestResponse) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.availableFurniture = make([]furniture, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.availableFurniture {
		(bbp.availableFurniture)[i1], err = makefurnitureFromBytes(buf[at:])
		if err != nil {
			return err
		}
		at += ((bbp.availableFurniture)[i1]).Size()
	}
	return nil
}

func (bbp *requestResponse) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.availableFurniture = make([]furniture, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.availableFurniture {
		(bbp.availableFurniture)[i1] = mustMakefurnitureFromBytes(buf[at:])
		at += ((bbp.availableFurniture)[i1]).Size()
	}
}

func (bbp *requestResponse) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.availableFurniture)))
	for _, elem := range bbp.availableFurniture {
		err = (elem).EncodeBebop(w)
		if err != nil {
			return err
		}
	}
	return w.Err
}

func (bbp *requestResponse) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.availableFurniture = make([]furniture, iohelp.ReadUint32(r))
	for i1 := range bbp.availableFurniture {
		((bbp.availableFurniture[i1])), err = makefurniture(r)
		if err != nil {
			return err
		}
	}
	return r.Err
}

func (bbp *requestResponse) Size() int {
	bodyLen := 0
	bodyLen += 4
	for _, elem := range bbp.availableFurniture {
		bodyLen += (elem).Size()
	}
	return bodyLen
}

func (bbp *requestResponse) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makerequestResponse(r iohelp.ErrorReader) (requestResponse, error) {
	v := requestResponse{}
	err := v.DecodeBebop(r)
	return v, err
}

func makerequestResponseFromBytes(buf []byte) (requestResponse, error) {
	v := requestResponse{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakerequestResponseFromBytes(buf []byte) requestResponse {
	v := requestResponse{}
	v.MustUnmarshalBebop(buf)
	return v
}

func (bbp *requestResponse) GetavailableFurniture() []furniture {
	return bbp.availableFurniture
}

func newrequestResponse(
		availableFurniture []furniture,
) requestResponse {
	return requestResponse{
		availableFurniture: availableFurniture,
	}
}

const requestCatalogOpCode = 0x41454b49

var _ bebop.Record = &requestCatalog{}

type requestCatalog struct {
	family *furnitureFamily
	// Deprecated: Nobody react to what I'm about to say...
	secretTunnel *string
}

func (bbp requestCatalog) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.family != nil {
		buf[at] = 1
		at++
		iohelp.WriteUint32Bytes(buf[at:], uint32(*bbp.family))
		at += 4
	}
	return at
}

func (bbp *requestCatalog) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.family = new(furnitureFamily)
			(*bbp.family) = furnitureFamily(iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
		case 2:
			at += 1
			bbp.secretTunnel = new(string)
			(*bbp.secretTunnel), err = iohelp.ReadStringBytes(buf[at:])
			if err != nil {
				return err
			}
			at += 4 + len((*bbp.secretTunnel))
		default:
			return nil
		}
	}
}

func (bbp *requestCatalog) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.family = new(furnitureFamily)
			(*bbp.family) = furnitureFamily(iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
		case 2:
			at += 1
			bbp.secretTunnel = new(string)
			(*bbp.secretTunnel) = iohelp.MustReadStringBytes(buf[at:])
			at += 4 + len((*bbp.secretTunnel))
		default:
			return
		}
	}
}

func (bbp requestCatalog) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.family != nil {
		w.Write([]byte{1})
		iohelp.WriteUint32(w, uint32(*bbp.family))
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *requestCatalog) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.family = new(furnitureFamily)
			*bbp.family = furnitureFamily(iohelp.ReadUint32(r))
		case 2:
			bbp.secretTunnel = new(string)
			*bbp.secretTunnel = iohelp.ReadString(r)
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp requestCatalog) Size() int {
	bodyLen := 5
	if bbp.family != nil {
		bodyLen += 1
		bodyLen += 4
	}
	return bodyLen
}

func (bbp *requestCatalog) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makerequestCatalog(r iohelp.ErrorReader) (requestCatalog, error) {
	v := requestCatalog{}
	err := v.DecodeBebop(r)
	return v, err
}

func makerequestCatalogFromBytes(buf []byte) (requestCatalog, error) {
	v := requestCatalog{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakerequestCatalogFromBytes(buf []byte) requestCatalog {
	v := requestCatalog{}
	v.MustUnmarshalBebop(buf)
	return v
}

