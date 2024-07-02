package generated

import (
	"os"
	"testing"
)

func TestInlineMessage(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "test.bbp")
	if err != nil {
		t.Fatal(err)
	}
	onOff := true
	inline := MessageInline{
		OnOff: &onOff,
	}
	snapshotLE := MessageInlineWrapper{
		Bla: &inline,
	}
	err = snapshotLE.EncodeBebop(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	reader, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1024)
	n, err := reader.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	buf = buf[:n]
	snapshotLE = MessageInlineWrapper{}
	err = snapshotLE.UnmarshalBebop(buf)
	if err != nil {
		t.Fatal(err)
	}
	if *snapshotLE.Bla.OnOff != true {
		t.Fatal("expected true")
	}
	reader.Close()

	reader, err = os.Open(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	snapshotLE = MessageInlineWrapper{}
	err = snapshotLE.DecodeBebop(reader)
	if *snapshotLE.Bla.OnOff != true {
		t.Fatal("expected true")
	}
	if err != nil {
		t.Fatal(err)
	}
}
