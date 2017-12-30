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
	StringType = NewPrimitiveType("string")
	BytesType  = NewPrimitiveType("bytes")

	BuiltInTypes = map[string]Type{
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
}

// NewSchema ...
func NewSchema(name string) *Schema {
	return &Schema{
		Name:  name,
		Paths: map[string]struct{}{},
		Types: NewCompositeType(nil, UnionKind, nil),
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
	d, err := s.parseFields(n)
	if err != nil {
		return err
	}
	t, err := Compile(UnionDef(d))
	if err != nil {
		return err
	}
	s.Types = t
	return nil
}

func (s *Schema) parseType(n *par.Node) (TypeDef, error) {
	if n.IsNoun() {
		t := BuiltInTypes[n.Text]
		if t != nil {
			return &BuiltIn{t}, nil
		}
		return NameRef(n.Text), nil
	}
	if ok, err := n.IsBracket("{", "}"); ok {
		if err != nil {
			return nil, err
		}
		return s.parseCompositeType(n)
	}
	if ok, err := n.IsInnerBracket("[", "]"); ok {
		if err != nil {
			return nil, err
		}
		return s.parseContainerType(n)
	}
	return nil, n.NewError(ErrIllegalSymbol)
}

func (s *Schema) parseCompositeType(n *par.Node) (TypeDef, error) {
	switch {
	case n.Node0.IsKeyword("enum"):
		n = n.Node1
		d, err := s.parseConsts(n.Node0)
		if err != nil {
			return nil, err
		}
		if n := n.Node1; !n.IsEmpty() {
			return nil, n.NewError(ErrExpectEmpty)
		}
		return EnumDef(d), nil
	case n.Node0.IsKeyword("union"):
		n = n.Node1
		d, err := s.parseFields(n.Node0)
		if err != nil {
			return nil, err
		}
		if n := n.Node1; !n.IsEmpty() {
			return nil, n.NewError(ErrExpectEmpty)
		}
		return UnionDef(d), nil
	case n.Node0.IsKeyword("bitfield"):
		n = n.Node1
		d, err := s.parseConsts(n.Node0)
		if err != nil {
			return nil, err
		}
		if n := n.Node1; !n.IsEmpty() {
			return nil, n.NewError(ErrExpectEmpty)
		}
		return BitfieldDef(d), nil
	case n.Node0.IsKeyword("struct"):
		n = n.Node1
		d, err := s.parseFields(n.Node0)
		if err != nil {
			return nil, err
		}
		if n := n.Node1; !n.IsEmpty() {
			return nil, n.NewError(ErrExpectEmpty)
		}
		return StructDef(d), nil
	default:
		return nil, n.Node0.NewError(ErrIllegalName)
	}
}

func (s *Schema) parseConsts(n *par.Node) ([]ConstDef, error) {
	d := []ConstDef{}
	err := parseLines(n, func(n *par.Node) error {
		tag, name, err := s.parseTagAndName(n)
		if err != nil {
			return err
		}
		d = append(d, ConstDef{tag, name})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Schema) parseFields(n *par.Node) ([]FieldDef, error) {
	d := []FieldDef{}
	err := parseLines(n, func(n *par.Node) error {
		tag, name, err := s.parseTagAndName(n.Node0)
		if err != nil {
			return err
		}
		f, err := s.parseType(n.Node1)
		if err != nil {
			return err
		}
		d = append(d, FieldDef{tag, name, nil, f})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Schema) parseContainerType(n *par.Node) (TypeDef, error) {
	switch {
	case n.Node0.IsKeyword("array"):
		n = n.Node1
		if n := n.Node0; !n.IsEmpty() {
			return nil, n.NewError(ErrExpectEmpty)
		}
		e, err := s.parseType(n.Node1)
		if err != nil {
			return nil, err
		}
		return &ArrayDef{e}, nil
	case n.Node0.IsKeyword("map"):
		n = n.Node1
		k, err := s.parseType(n.Node0)
		if err != nil {
			return nil, err
		}
		e, err := s.parseType(n.Node1)
		if err != nil {
			return nil, err
		}
		return &MapDef{k, e}, nil
	case n.Node0.IsKeyword("set"):
		n = n.Node1
		k, err := s.parseType(n.Node0)
		if err != nil {
			return nil, err
		}
		if n := n.Node1; !n.IsEmpty() {
			return nil, n.NewError(ErrExpectEmpty)
		}
		return &SetDef{k}, nil
	default:
		return nil, n.Node0.NewError(ErrIllegalName)
	}
}

func parseLines(n *par.Node, callback func(n *par.Node) error) error {
	return n.Repeat("\n", callback)
}

func (s *Schema) parseTagAndName(n *par.Node) (int, string, error) {
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

func parseTag(n *par.Node) (int, error) {
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
	return int(tag), nil
}

func parseName(n *par.Node) (string, error) {
	if !n.IsNoun() {
		return "", n.NewError(ErrExpectName)
	}
	return n.Text, nil
}
