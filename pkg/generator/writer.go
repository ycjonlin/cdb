package generator

import (
	"io"
	"reflect"
	"strconv"
	"unsafe"

	sch "github.com/ycjonlin/cdb/pkg/schema"
)

type writer struct {
	io.Writer
	newline []byte
	indent  []byte
}

// WriteDeclaration ...
func WriteDeclaration(wr io.Writer, s *sch.Schema, newline, indent string) {
	w := writer{wr, []byte(newline), []byte(indent)}
	var c declarator = &declaratorImpl{&w, &typerImpl{&w}}
	c.putSchema(s)
}

// WriteBytesMarshalling ...
func WriteBytesMarshalling(wr io.Writer, s *sch.Schema, newline, indent string) {
	w := writer{wr, []byte(newline), []byte(indent)}
	var t typer = &typerImpl{&w}
	var d marshallerDef = &bytesMarshaller{&w, t}
	var c marshaller = &marshallerImpl{&w, t, d}
	c.putSchema(s)
}

// WriteMapMarshalling ...
func WriteMapMarshalling(wr io.Writer, s *sch.Schema, newline, indent string) {
	w := writer{wr, []byte(newline), []byte(indent)}
	var t typer = &typerImpl{&w}
	var d marshallerDef = &mapMarshaller{&w, t}
	var c marshaller = &marshallerImpl{&w, t, d}
	c.putSchema(s)
}

func (w *writer) putString(s string) {
	if len(s) == 0 {
		return
	}
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	arg := (*[0x7fffffff]byte)(unsafe.Pointer(hdr.Data))[:len(s):len(s)]
	w.Write(arg)
}

func (w *writer) putUint(u uint64) {
	w.Write(strconv.AppendUint(nil, u, 10))
}

func (w *writer) putHexUint(u uint64) {
	w.Write(strconv.AppendUint(nil, u, 16))
}

func (w *writer) putNewline() {
	w.Write(w.newline)
}

func (w *writer) putLine(s string) {
	w.putNewline()
	w.putString(s)
}

func (w *writer) pushIndent() {
	w.newline = append(w.newline, w.indent...)
}

func (w *writer) popIndent() {
	w.newline = w.newline[:len(w.newline)-len(w.indent)]
}

func (w *writer) putTag(t sch.Tag) {
	w.putUint(uint64(t))
}

func (w *writer) putName(n sch.Name) {
	w.putString(string(n))
}

func (w *writer) putCompoundName(r *sch.ReferenceType) {
	if r.Super != nil {
		w.putCompoundName(r.Super)
		w.putString("_")
	}
	w.putName(r.Name)
}
