package encoding

// Buffer ...
type Buffer struct {
	Bytes
	offset int
}

// Get ...
func (b *Buffer) Get() []byte {
	v := b.Bytes[b.offset:]
	b.offset = len(b.Bytes)
	return v
}

// Path ...
type Path struct {
	Bytes
	offsets []int
}

// Push ...
func (p *Path) Push() {
	p.offsets = append(p.offsets, len(p.Bytes))
}

// Pop ...
func (p *Path) Pop() {
	n := len(p.offsets) - 1
	o := p.offsets[n]
	p.offsets = p.offsets[:n]
	p.Bytes = p.Bytes[:o]
}

// Pair ...
type Pair struct {
	Key, Value Bytes
}

// Map ...
type Map struct {
	Data  Buffer
	Path  Path
	Pairs []Pair
}

// Put ...
func (m *Map) Put() {
	k := m.Path.Bytes
	v := m.Data.Get()
	m.Pairs = append(m.Pairs, Pair{k, v})
}
