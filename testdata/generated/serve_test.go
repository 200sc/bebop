package generated

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

type testServeMux struct {
	printed string
	added   int32
	handled int
}

func (tsm *testServeMux) HandlePrint(p Print) {
	fmt.Println(p.Printout)
	tsm.printed = p.Printout
	tsm.handled++
}

func (tsm *testServeMux) HandleAdd(a Add) AddResponse {
	c := a.A + a.B
	tsm.added = c
	tsm.handled++
	return AddResponse{
		C: c,
	}
}

func TestServe(t *testing.T) {
	mx := &testServeMux{}
	s := Server{}
	s.Mux = mx

	lhs, rhs := net.Pipe()

	go func() {
		err := s.Serve(rhs)
		if err != nil {
			panic(err)
		}
	}()

	printBytes := PrintRequest{
		Print: &Print{
			Printout: "Hello World",
		},
	}.MarshalBebop()
	err := binary.Write(lhs, binary.LittleEndian, int32(PrintRequestOpCode))
	if err != nil {
		t.Fatal(err)
	}
	_, err = lhs.Write(printBytes)
	if err != nil {
		t.Fatal(err)
	}

	addBytes := AddRequest{
		Add: &Add{
			A: 42,
			B: 42,
		},
	}.MarshalBebop()
	err = binary.Write(lhs, binary.LittleEndian, int32(AddRequestOpCode))
	if err != nil {
		t.Fatal(err)
	}
	_, err = lhs.Write(addBytes)
	if err != nil {
		t.Fatal(err)
	}
	retries := 9
	for mx.handled != 2 && retries > 0 {
		time.Sleep(10 * time.Millisecond)
		retries--
	}
	if mx.handled != 2 {
		t.Fatal("mux did not handle our messages")
	}

	resp := AddResponse{}
	ln := resp.Size()
	respBytes := make([]byte, ln)
	_, err = io.ReadFull(lhs, respBytes)
	if err != nil {
		t.Fatal(err)
	}
	err = resp.UnmarshalBebop(respBytes)
	if err != nil {
		t.Fatal(err)
	}
	if resp.C != 84 {
		t.Fatal("C was not 84, was:", resp.C)
	}
	if mx.added != resp.C {
		t.Fatal("Mx's C was not 84, was:", mx.added)
	}
	if mx.printed != "Hello World" {
		t.Fatal("Mx's printed was not 'Hello World', was:", mx.printed)
	}
}
