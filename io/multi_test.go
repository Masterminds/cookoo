package io

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestMultiWrite(t *testing.T) {
	sha1 := sha1.New()
	sink := new(bytes.Buffer)
	mw := NewMultiWriter()
	mw.(*MultiWriter).AddWriter("sha1", sha1)
	mw.(*MultiWriter).AddWriter("sink", sink)

	sourceString := "My input text."
	source := strings.NewReader(sourceString)
	written, err := io.Copy(mw, source)

	if written != int64(len(sourceString)) {
		t.Errorf("short write of %d, not %d", written, len(sourceString))
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	sha1hex := fmt.Sprintf("%x", sha1.Sum(nil))
	if sha1hex != "01cb303fa8c30a64123067c5aa6284ba7ec2d31b" {
		t.Error("incorrect sha1 value")
	}

	if sink.String() != sourceString {
		t.Errorf("expected %q; got %q", sourceString, sink.String())
	}
}

func TestMultiWriterCRUD(t *testing.T) {
	sha1 := sha1.New()
	mw := NewMultiWriter()
	mw.(*MultiWriter).AddWriter("sha1", sha1)

	sha1a, found := mw.(*MultiWriter).Writer("sha1")

	if found == false {
		t.Error("Did not find sha1 as expected.")
	}
	if sha1a != sha1 {
		t.Error("Expected sha1 returned from MultiWriter to be what was set. They were different.")
	}

	mw.(*MultiWriter).RemoveWriter("sha1")
	_, found = mw.(*MultiWriter).Writer("sha1")

	if found == true {
		t.Error("Expected sha1 to be removed from MultiWriter but it was not.")
	}

}
