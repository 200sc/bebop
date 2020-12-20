package generated

import (
	"fmt"
	"io"

	"github.com/200sc/bebop/iohelp"
)

type ServeMux interface {
	HandlePrint(Print)
	HandleAdd(Add) AddResponse
}

type Server struct {
	Mux ServeMux
}

func (s Server) Serve(rw io.ReadWriter) error {
	er := iohelp.NewErrorReader(rw)
	for {
		opCode := iohelp.ReadUint32(er)
		if er.Err != nil {
			fmt.Println("DEBUG: opcode read error", er.Err)
			continue
		}
		bodyLen := iohelp.ReadUint32(er)
		if er.Err != nil {
			fmt.Println("DEBUG: bodyLen read error", er.Err)
			continue
		}
		body := make([]byte, bodyLen)
		_, err := io.ReadFull(er, body)
		if err != nil {
			fmt.Println("DEBUG: body read error", err)
			continue
		}
		if body[0] != 1 {
			fmt.Println("DEBUG: body invalid error")
			continue
		}
		switch opCode {
		case PrintRequestOpCode:
			p := Print{}
			parsingBody := body[1:]
			err = p.UnmarshalBebop(parsingBody)
			if err != nil {
				fmt.Println("DEBUG: unmarshal print", err)
			}
			if s.Mux != nil {
				s.Mux.HandlePrint(p)
			}
		case AddRequestOpCode:
			a := Add{}
			err = a.UnmarshalBebop(body[1:])
			if err != nil {
				fmt.Println("DEBUG: unmarshal add", err)
			}
			if s.Mux != nil {
				resp := s.Mux.HandleAdd(a)
				respBytes := resp.MarshalBebop()
				if _, err := rw.Write(respBytes); err != nil {
					fmt.Println("DEBUG: write add response", err)
				}
			}
		}
	}
}
