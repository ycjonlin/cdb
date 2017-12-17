package schema

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	par "github.com/ycjonlin/cdb/pkg/parser"
)

// Err ...
var (
	ErrTODO          = errors.New("TODO")
	ErrIllegalTag    = errors.New("illegal tag")
	ErrIllegalName   = errors.New("illegal name")
	ErrIllegalSymbol = errors.New("illegal symbol")
	ErrIllegalType   = errors.New("illegal type")

	ErrExpectEmpty     = errors.New("expect to be empty")
	ErrExpectType      = errors.New("expect a type")
	ErrExpectSpace     = errors.New("expect a space")
	ErrExpectSemicolon = errors.New("expect a semicolon")
	ErrExpectDot       = errors.New("expect a dot")
	ErrExpectTag       = errors.New("expect a tag")
	ErrExpectName      = errors.New("expect a name")
)

var parser = &par.Parser{
	Noun: regexp.MustCompile("^(?s)" +
		"[\\pL\\pN]+|\"(?:[^\"\\\\]|\\\\.)*\"" +
		""),
	Verb: regexp.MustCompile("^(?s)" +
		("(?:" +
			"[\n(),;[\\]{}]" + "|" +
			"[!$%&'*+\\-./:<=>?@\\\\^_`|~]+" +
			")") +
		""),
	Space: regexp.MustCompile("^(?s)" +
		("(?:" +
			"[\t\v\f\r ]*" + "|" +
			"##.*?(?:##|$)" + "|" +
			"#[^\n]*" +
			")*") +
		""),
	Brackets: par.MakeBrackets([][2]string{
		{"{", "}"},
		{"[", "]"},
		{"(", ")"},
	}),
	Associativities: par.MakeAssociativities(
		"\n",
		"}", "]", ")",
		"{", "[", "(",
		"",
	),
	Priorities: par.MakePriorities(
		[]string{"[", "]", "{", "}", "(", ")", "."},
		[]string{""},
		[]string{"\n"},
	),
}

// Err ...
var (
	BoolType   = NewPrimitiveType("bool")
	IntType    = NewPrimitiveType("int")
	UintType   = NewPrimitiveType("uint")
	FloatType  = NewPrimitiveType("float")
	BytesType  = NewPrimitiveType("bytes")
	StringType = NewPrimitiveType("string")

	BuiltInTypes = map[Name]Type{
		"bool":      BoolType,
		"int":       IntType,
		"uint":      UintType,
		"float":     FloatType,
		"bytes":     BytesType,
		"string":    StringType,
		"timestamp": IntType,
		"duration":  IntType,
		"geopoint":  UintType,
	}
)

// Schema ...
type Schema struct {
	Name  string
	Paths map[string]struct{}
	Types *CompositeType

	refs []ref
}

type ref struct {
	r *ReferenceType
	n *par.Node
}

// NewSchema ...
func NewSchema(name string) *Schema {
	return &Schema{
		Name:  name,
		Paths: map[string]struct{}{},
		Types: NewCompositeType(nil, UnionKind),
	}
}

// Parse ...
func (s *Schema) Parse(path string) error {
	// skip read files
	if _, ok := s.Paths[path]; ok {
		return nil
	}
	s.Paths[path] = struct{}{}
	// read file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	source, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	// parse file
	ast, err := parser.Parse(source, path)
	if err != nil {
		return err
	}
	return s.parseSchema(ast)
}

func (s *Schema) parseSchema(n *par.Node) error {
	err := s.parseCompositeTypeBody(s.Types, n)
	if err != nil {
		return err
	}
	for _, ref := range s.refs {
		err = s.parseReferenceType(ref.r, ref.n)
		if err != nil {
			return err
		}
	}
	s.refs = nil
	return nil
}

func (s *Schema) parseType(r *ReferenceType, n *par.Node) error {
	if n.IsNoun() {
		for _, ref := range s.refs {
			if r == ref.r {
				panic("r")
			}
			if n == ref.n {
				panic("n")
			}
		}
		s.refs = append(s.refs, ref{r, n})
		return nil
	}
	if ok, err := n.IsBracket("{", "}"); ok {
		if err != nil {
			return err
		}
		return s.parseCompositeType(r, n)
	}
	if ok, err := n.IsInnerBracket("[", "]"); ok {
		if err != nil {
			return err
		}
		return s.parseContainerType(r, n)
	}
	return n.NewError(ErrIllegalSymbol)
}

func (s *Schema) parseReferenceType(r *ReferenceType, n *par.Node) error {
	name, err := parseName(n)
	if err != nil {
		return err
	}
	bt := BuiltInTypes[name]
	if bt != nil {
		err = r.SetType(bt)
		if err != nil {
			return n.NewError(err)
		}
		return nil
	}
	rt := s.Types.GetByName(name)
	if rt != nil {
		err = r.SetType(rt)
		if err != nil {
			return n.NewError(err)
		}
		return nil
	}
	return n.NewError(ErrIllegalName)
}

func (s *Schema) parseCompositeType(r *ReferenceType, n *par.Node) error {
	var kind Kind
	switch {
	case n.Node0.IsKeyword("enum"):
		kind = EnumKind
	case n.Node0.IsKeyword("bitfield"):
		kind = BitfieldKind
	case n.Node0.IsKeyword("union"):
		kind = UnionKind
	case n.Node0.IsKeyword("struct"):
		kind = StructKind
	default:
		return n.Node0.NewError(ErrIllegalName)
	}
	n = n.Node1
	if n := n.Node1; !n.IsEmpty() {
		return n.NewError(ErrExpectEmpty)
	}
	c := NewCompositeType(r, kind)
	return s.parseCompositeTypeBody(c, n.Node0)
}

func (s *Schema) parseCompositeTypeBody(c *CompositeType, n *par.Node) error {
	isTyped := c.Kind == UnionKind || c.Kind == StructKind
	return parseLines(n, func(n *par.Node) error {
		var np *par.Node
		if isTyped {
			np = n.Node1
			n = n.Node0
		}
		f, err := s.parseCompositeTypeField(c, n)
		if err != nil {
			return err
		}
		if isTyped {
			return s.parseType(f, np)
		}
		return nil
	})
}

func (s *Schema) parseCompositeTypeField(c *CompositeType, n *par.Node) (*ReferenceType, error) {
	tag, name, err := s.parseTagAndName(n)
	if err != nil {
		return nil, err
	}
	f, err := c.Put(tag, name)
	if err != nil {
		if err == ErrDuplicatedTypeTag {
			return nil, n.Node0.NewError(err)
		}
		if err == ErrDuplicatedTypeName {
			return nil, n.Node1.NewError(err)
		}
		return nil, n.NewError(err)
	}
	return f, nil
}

func (s *Schema) parseContainerType(r *ReferenceType, n *par.Node) error {
	var kind Kind
	switch {
	case n.Node0.IsKeyword("array"):
		kind = ArrayKind
	case n.Node0.IsKeyword("map"):
		kind = MapKind
	case n.Node0.IsKeyword("set"):
		kind = SetKind
	default:
		return n.Node0.NewError(ErrIllegalName)
	}
	c := NewContainerType(r, kind)
	return s.parseContainerTypeBody(c, n.Node1)
}

func (s *Schema) parseContainerTypeBody(c *ContainerType, n *par.Node) error {
	isKeyed := c.Kind == MapKind || c.Kind == SetKind
	isTyped := c.Kind == ArrayKind || c.Kind == MapKind
	if n := n.Node0; isKeyed {
		err := s.parseType(c.Key, n)
		if err != nil {
			return err
		}
	} else {
		if !n.IsEmpty() {
			return n.NewError(ErrExpectEmpty)
		}
	}
	if n := n.Node1; isTyped {
		err := s.parseType(c.Elem, n)
		if err != nil {
			return err
		}
	} else {
		if !n.IsEmpty() {
			return n.NewError(ErrExpectEmpty)
		}
	}
	return nil
}

func parseLines(n *par.Node, callback func(n *par.Node) error) error {
	return n.Repeat("\n", callback)
}

func (s *Schema) parseTagAndName(n *par.Node) (Tag, Name, error) {
	if !n.IsVerb(".") {
		return 0, "", n.NewError(ErrExpectDot)
	}
	tag, err := parseTag(n.Node0)
	if err != nil {
		return 0, "", err
	}
	name, err := parseName(n.Node1)
	if err != nil {
		return 0, "", err
	}
	return tag, name, nil
}

func parseTag(n *par.Node) (Tag, error) {
	if !n.IsNoun() {
		return 0, n.NewError(ErrIllegalTag)
	}
	text := n.Text
	base := 10
	switch text[0] {
	case 'b':
		text, base = text[1:], 2
	case 'o':
		text, base = text[1:], 8
	case 'd':
		text, base = text[1:], 10
	case 'h':
		text, base = text[1:], 16
	}
	tag, err := strconv.ParseUint(text, base, 31)
	if err != nil {
		return 0, n.NewError(ErrExpectTag)
	}
	if tag <= 0 {
		return 0, n.NewError(ErrIllegalTag)
	}
	return Tag(tag), nil
}

func parseName(n *par.Node) (Name, error) {
	if !n.IsNoun() {
		return "", n.NewError(ErrExpectName)
	}
	return Name(n.Text), nil
}
