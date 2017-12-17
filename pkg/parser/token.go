package parser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// Node ...
type Node struct {
	Text   string
	Path   string
	Offset int
	Lineno int
	Charno int
	Node0  *Node
	Node1  *Node
}

// NodeError ...
type NodeError struct {
	Err  error
	Node *Node
}

// NewNode ...
func NewNode(text string) *Node {
	return &Node{Text: text}
}

// ToJSON ...
func (n *Node) ToJSON(w io.Writer) {
	w.Write([]byte(`{"name":`))
	w.Write(strconv.AppendQuote(nil, n.Text))
	if n.Node0 != nil {
		w.Write([]byte(`,"children":[`))
		n.Node0.ToJSON(w)
		w.Write([]byte(`,`))
		n.Node1.ToJSON(w)
		w.Write([]byte(`]`))
	}
	w.Write([]byte(`}`))
}

// ToText ...
func (n *Node) ToText(w io.Writer, indent []byte) {
	if n.Node0 != nil {
		w.Write([]byte("┬"))
	} else {
		w.Write([]byte("─"))
	}
	w.Write(strconv.AppendQuote(nil, n.Text))
	w.Write([]byte("\n"))
	if n.Node0 != nil {
		w.Write(indent)
		w.Write([]byte("├"))
		n.Node0.ToText(w, append(indent, []byte("│")...))
		w.Write(indent)
		w.Write([]byte("└"))
		n.Node1.ToText(w, append(indent, []byte(" ")...))
	}
}

func (n *Node) String() string {
	if n.Node0 == nil {
		return n.Text
	}
	return n.Node0.String() + n.Text + n.Node1.String()
}

// IsVerb ...
func (n *Node) IsVerb(text string) bool {
	return n.Node0 != nil && n.Text == text
}

// IsAnyVerb ...
func (n *Node) IsAnyVerb(texts map[string]struct{}) bool {
	_, ok := texts[n.Text]
	return n.Node0 != nil && ok
}

// IsNoun ...
func (n *Node) IsNoun() bool {
	return n.Node0 == nil && len(n.Text) > 0
}

// IsKeyword ...
func (n *Node) IsKeyword(text string) bool {
	return n.Node0 == nil && n.Text == text
}

// IsAnyKeyword ...
func (n *Node) IsAnyKeyword(texts map[string]struct{}) bool {
	_, ok := texts[n.Text]
	return n.Node0 == nil && ok
}

// IsEmpty ...
func (n *Node) IsEmpty() bool {
	return n.Node0 == nil && len(n.Text) == 0
}

// IsBracket ...
func (n *Node) IsBracket(opening, closing string) (bool, error) {
	if !n.IsVerb(opening) {
		return false, nil
	}
	if !n.Node1.IsVerb(closing) {
		return true, &NodeError{ErrUnpairedBracket, n}
	}
	if !n.Node1.Node1.IsEmpty() {
		return true, &NodeError{ErrTrailedBracket, n}
	}
	return true, nil
}

// IsInnerBracket ...
func (n *Node) IsInnerBracket(opening, closing string) (bool, error) {
	if !n.IsVerb(opening) {
		return false, nil
	}
	if !n.Node1.IsVerb(closing) {
		return true, &NodeError{ErrUnpairedBracket, n}
	}
	return true, nil
}

// Callback ...
type Callback func(*Node) error

// Repeat ...
func (n *Node) Repeat(delimiter string, callback Callback) error {
	for n.IsVerb(delimiter) {
		n0 := n.Node0
		if !n0.IsEmpty() {
			if err := callback(n0); err != nil {
				return err
			}
		}
		n = n.Node1
	}
	n0 := n
	if !n0.IsEmpty() {
		if err := callback(n0); err != nil {
			return err
		}
	}
	return nil
}

// NewError ...
func (n *Node) NewError(err error) error {
	return &NodeError{err, n}
}

// Error ...
func (e *NodeError) Error() string {
	n := e.Node
	b := bytes.NewBuffer(nil)
	fmt.Fprintf(b, "%s - %s:%d:%d\n", e.Err.Error(), n.Path, n.Lineno+1, n.Charno+1)
	n.ToText(b, nil)
	return b.String()
}
