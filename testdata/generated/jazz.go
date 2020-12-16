// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
)

type Instrument uint32

const (
	Instrument_Sax Instrument = 0
	Instrument_Trumpet Instrument = 1
	Instrument_Clarinet Instrument = 2
)

var _ bebop.Record = &Musician{}

type Musician struct {
	name string
	plays Instrument
}

func (bbp Musician) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.name)))
	w.Write([]byte(bbp.name))
	iohelp.WriteUint32(w, uint32(bbp.plays))
	return w.Err
}

func (bbp *Musician) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	{
		bbp.name = iohelp.ReadString(r)
	}
	{
		binary.Read(r, binary.LittleEndian, &bbp.plays)
	}
	return r.Err
}

func (bbp *Musician) bodyLen() uint32 {
	bodyLen := uint32(0)
	bodyLen += 4
	bodyLen += uint32(len(bbp.name))
	bodyLen += 4
	return bodyLen
}

func (bbp Musician) GetName() string {
	return bbp.name
}

func (bbp Musician) GetPlays() Instrument {
	return bbp.plays
}

var _ bebop.Record = &Library{}

type Library struct {
	Songs map[[16]byte]Song
}

func (bbp Library) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.Songs)))
	for k, v := range bbp.Songs {
		iohelp.WriteGUID(w, k)
		err = (v).EncodeBebop(w)
		if err != nil {
			return err
		}
	}
	return w.Err
}

func (bbp *Library) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	{
		ln2 := iohelp.ReadUint32(r)
		bbp.Songs = make(map[[16]byte]Song, ln2)
		for i2 := uint32(0); i2 < ln2; i2++ {
			k := new([16]byte)
			*k = iohelp.ReadGUID(r)
			elem2 := new(Song)
			err = (elem2).DecodeBebop(r)
			if err != nil {
				return err
			}
			(bbp.Songs)[*k] = *elem2
		}
	}
	return r.Err
}

func (bbp *Library) bodyLen() uint32 {
	bodyLen := uint32(0)
	bodyLen += 4
	for _, v := range bbp.Songs {
		bodyLen += 16
		bodyLen += (v).bodyLen()
	}
	return bodyLen
}

var _ bebop.Record = &Song{}

type Song struct {
	Title *string
	Year *uint16
	Performers *[]Musician
}

func (bbp Song) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	binary.Write(w, binary.LittleEndian, bbp.bodyLen())
	if bbp.Title != nil {
		w.Write([]byte{1})
		iohelp.WriteUint32(w, uint32(len(*bbp.Title)))
		w.Write([]byte(*bbp.Title))
	}
	if bbp.Year != nil {
		w.Write([]byte{2})
		iohelp.WriteUint16(w, *bbp.Year)
	}
	if bbp.Performers != nil {
		w.Write([]byte{3})
		binary.Write(w, binary.LittleEndian, uint32(len(*bbp.Performers)))
		for _, elem := range *bbp.Performers {
			err = (elem).EncodeBebop(w)
			if err != nil {
				return err
			}
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *Song) DecodeBebop(ior io.Reader) (err error) {
	er := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(er)
	body := make([]byte, bodyLen)
	er.Read(body)
	r := iohelp.NewErrorReader(bytes.NewReader(body))
	for {
		switch iohelp.ReadByte(r) {
		case 0:
			return er.Err
		case 1:
			bbp.Title = new(string)
			*bbp.Title = iohelp.ReadString(r)
		case 2:
			bbp.Year = new(uint16)
			*bbp.Year = iohelp.ReadUint16(r)
		case 3:
			bbp.Performers = new([]Musician)
			ln3 := uint32(0)
			binary.Read(r, binary.LittleEndian, &ln3)
			for i := uint32(0); i < ln3; i++ {
				elem3 := new(Musician)
				err = (elem3).DecodeBebop(r)
				if err != nil {
					return err
				}
			}
		default:
			return er.Err
		}
	}
	return er.Err
}

func (bbp *Song) bodyLen() uint32 {
	bodyLen := uint32(1)
	if bbp.Title != nil {
		bodyLen += 1
		bodyLen += 4
		bodyLen += uint32(len(*bbp.Title))
	}
	if bbp.Year != nil {
		bodyLen += 1
		bodyLen += 2
	}
	if bbp.Performers != nil {
		bodyLen += 1
		bodyLen += 4
		for _, elem := range *bbp.Performers {
			bodyLen += (elem).bodyLen()
		}
	}
	return bodyLen
}

