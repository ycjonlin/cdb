package encoding

import (
	"bytes"
	"errors"
)

// Err ...
var (
	ErrTODO = errors.New("TODO")
)

// Data ...
type Data interface {
	SetKey() *Bytes
	GetKey() *Bytes

	Value() *Bytes
	Len() *Bytes
	End() *Bytes
	Next() bool

	Data(*Bytes) Data
}

// BytesData ...
type BytesData Bytes

// NewBytesData ...
func NewBytesData() Data {
	return &BytesData{}
}

// String ...
func (b *BytesData) String() string { return (*Bytes)(b).String() }

// SetKey ...
func (b *BytesData) SetKey() *Bytes { return (*Bytes)(b) }

// GetKey ...
func (b *BytesData) GetKey() *Bytes { return (*Bytes)(b) }

// Value ...
func (b *BytesData) Value() *Bytes { return (*Bytes)(b) }

// Len ...
func (b *BytesData) Len() *Bytes { return (*Bytes)(b) }

// End ...
func (b *BytesData) End() *Bytes { return (*Bytes)(b) }

// Next ...
func (b *BytesData) Next() bool { return true }

// Data ...
func (*BytesData) Data(b *Bytes) Data { return (*BytesData)(b) }

type pair struct{ key, value Bytes }

func (p *pair) String() string {
	return p.key.String() + ":" + p.value.String()
}

type node struct {
	pair
	nodes []node
}

func (n *node) String() string {
	var s string
	s += n.key.String() + ":" + n.value.String()
	if len(n.nodes) != 0 {
		s += "{"
		for i, np := range n.nodes {
			if i != 0 {
				s += ", "
			}
			s += np.String()
		}
		s += "}"
	}
	return s
}

// MapData ...
type MapData struct {
	path  []*pair
	keys  []Bytes
	pairs []pair
}

// NewMapData ...
func NewMapData() Data {
	return &MapData{}
}

// String ...
func (m *MapData) String() string {
	var s string
	s += "{"
	for i, p := range m.pairs {
		if i != 0 {
			s += ", "
		}
		s += p.key.String() + ":" + p.value.String()
	}
	s += "}"
	return s
}

// SetKey ...
func (m *MapData) SetKey() *Bytes {
	var k Bytes
	if len(m.path) > 0 {
		k = append(k, m.path[len(m.path)-1].key...)
	}
	m.pairs = append(m.pairs, pair{k, nil})
	p := &m.pairs[len(m.pairs)-1]
	m.path = append(m.path, p)
	m.keys = append(m.keys, nil)
	return &p.key
}

// GetKey ...
func (m *MapData) GetKey() *Bytes {
	var k Bytes
	if len(m.path) > 0 {
		k = m.keys[len(m.keys)-1]
	}
	p := &m.pairs[0]
	m.pairs = m.pairs[1:]
	m.path = append(m.path, p)
	m.keys = append(m.keys, p.key)
	p.key = p.key[len(k):]
	return &p.key
}

// Value ...
func (m *MapData) Value() *Bytes {
	p := m.path[len(m.path)-1]
	m.path = m.path[:len(m.path)-1]
	m.keys = m.keys[:len(m.keys)-1]
	return &p.value
}

// Len ...
func (m *MapData) Len() *Bytes { return nil }

// End ...
func (m *MapData) End() *Bytes { return nil }

// Next ...
func (m *MapData) Next() bool {
	if len(m.pairs) == 0 {
		return false
	}
	p := &m.pairs[0]
	k := m.keys[len(m.keys)-1]
	return bytes.HasPrefix(p.key, k)
}

// Data ...
func (*MapData) Data(b *Bytes) Data { return (*BytesData)(b) }

// TreeData ...
type TreeData struct {
	path []*node
}

// NewTreeData ...
func NewTreeData() Data {
	n := &node{}
	return &TreeData{path: []*node{n}}
}

// String ...
func (t *TreeData) String() string {
	var s string
	for _, n := range t.path {
		s += n.String() + " "
	}
	return s
}

// SetKey ...
func (t *TreeData) SetKey() *Bytes {
	n := t.path[len(t.path)-1]
	n.nodes = append(n.nodes, node{})
	np := &n.nodes[len(n.nodes)-1]
	t.path = append(t.path, np)
	return &np.key
}

// GetKey ...
func (t *TreeData) GetKey() *Bytes {
	n := t.path[len(t.path)-1]
	np := &n.nodes[0]
	n.nodes = n.nodes[1:]
	t.path = append(t.path, np)
	return &np.key
}

// Value ...
func (t *TreeData) Value() *Bytes {
	n := t.path[len(t.path)-1]
	t.path = t.path[:len(t.path)-1]
	return &n.value
}

// Len ...
func (t *TreeData) Len() *Bytes { return nil }

// End ...
func (t *TreeData) End() *Bytes { return nil }

// Next ...
func (t *TreeData) Next() bool {
	n := t.path[len(t.path)-1]
	return len(n.nodes) != 0
}

// Data ...
func (*TreeData) Data(b *Bytes) Data { return (*BytesData)(b) }
