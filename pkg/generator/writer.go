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

// WriteMarshalling ...
func WriteMarshalling(wr io.Writer, s *sch.Schema, newline, indent string) {
	w := writer{wr, []byte(newline), []byte(indent)}
	var c marshaller = &marshallerImpl{&w, &typerImpl{&w}}
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

func (w *writer) pushIndent(s string) {
	w.putLine(s)
	w.newline = append(w.newline, w.indent...)
}

func (w *writer) pushIndentCont(s string) {
	w.putString(s)
	w.newline = append(w.newline, w.indent...)
}

func (w *writer) popIndent(s string) {
	w.newline = w.newline[:len(w.newline)-len(w.indent)]
	w.putLine(s)
}

func (w *writer) putTag(t int) {
	w.putUint(uint64(t))
}

func (w *writer) putName(n string) {
	w.putString(string(n))
}

func (w *writer) putCompoundName(r *sch.ReferenceType) {
	if r.Super != nil {
		w.putCompoundName(r.Super)
		w.putString("_")
	}
	w.putName(r.Name)
}
