package parser

import (
	"errors"
	"regexp"
)

// Err ...
var (
	ErrUnrecognizedVerb  = errors.New("parser: unrecognized verb")
	ErrUnpairedBracket   = errors.New("parser: unpaired bracket")
	ErrTrailedBracket    = errors.New("parser: trailed bracket")
	ErrNotMakingProgress = errors.New("parser: not making progress")
)

// BracketConfig ...
type BracketConfig struct {
	Openings map[string]map[string]struct{}
	Closings map[string]map[string]struct{}
}

// Parser ...
type Parser struct {
	Noun  *regexp.Regexp
	Verb  *regexp.Regexp
	Space *regexp.Regexp

	Keywords        map[string]struct{}
	Brackets        BracketConfig
	Priorities      map[string]int
	Associativities map[string]struct{}
}

// MakePriorities ...
func MakePriorities(tiers ...[]string) map[string]int {
	priorities := map[string]int{}
	for i, tier := range tiers {
		for _, name := range tier {
			priorities[name] = i + 1
		}
	}
	return priorities
}

// MakeKeywords ...
func MakeKeywords(names ...string) map[string]struct{} {
	keywords := map[string]struct{}{}
	for _, name := range names {
		keywords[name] = struct{}{}
	}
	return keywords
}

// MakeAssociativities ...
func MakeAssociativities(names ...string) map[string]struct{} {
	associativities := map[string]struct{}{}
	for _, name := range names {
		associativities[name] = struct{}{}
	}
	return associativities
}

// MakeBrackets ...
func MakeBrackets(brackets [][2]string) BracketConfig {
	openings := map[string]map[string]struct{}{}
	closings := map[string]map[string]struct{}{}
	for _, bracket := range brackets {
		opening := bracket[0]
		closing := bracket[1]
		if _, ok := openings[closing]; !ok {
			openings[closing] = map[string]struct{}{}
		}
		openings[closing][opening] = struct{}{}
		if _, ok := closings[opening]; !ok {
			closings[opening] = map[string]struct{}{}
		}
		closings[opening][closing] = struct{}{}
	}
	return BracketConfig{
		Closings: closings,
		Openings: openings,
	}
}

type verb struct {
	node          *Node
	priority      int
	associativity bool
}

type bracket struct {
	node  *Node
	index int
}

type cursor struct {
	source     []byte
	offset     int
	path       string
	lineno     int
	charno     int
	lastOffset int

	nouns    []*Node
	verbs    []verb
	brackets []bracket
}

func (c *cursor) pushNoun(node *Node) (err error) {
	c.nouns = append(c.nouns, node)
	return nil
}

func (c *cursor) pushVerb(node *Node, priority int, associativity bool) (err error) {
	lastBracket := c.brackets[len(c.brackets)-1]
	for len(c.verbs) > lastBracket.index {
		lastVerb := &c.verbs[len(c.verbs)-1]
		if lastVerb.priority > priority ||
			(lastVerb.priority == priority && associativity) {
			break
		}
		c.popVerb()
	}
	c.verbs = append(c.verbs, verb{node, priority, associativity})
	return nil
}

func (c *cursor) popVerb() {
	n := len(c.nouns)
	v := len(c.verbs)

	noun := c.verbs[v-1].node
	noun.Node0 = c.nouns[n-2]
	noun.Node1 = c.nouns[n-1]
	c.nouns[n-2] = noun

	c.nouns = c.nouns[:n-1]
	c.verbs = c.verbs[:v-1]
}

func (c *cursor) pushBracket(node *Node) (err error) {
	newBracket := bracket{node: node, index: len(c.verbs)}
	c.brackets = append(c.brackets, newBracket)
	return nil
}

func (c *cursor) popBracket() {
	lastBracket := c.brackets[len(c.brackets)-1]
	for len(c.verbs) > lastBracket.index {
		c.popVerb()
	}
	c.brackets = c.brackets[:len(c.brackets)-1]
}

func (c *cursor) node(segment string) *Node {
	return &Node{segment, c.path, c.offset, c.lineno, c.charno, nil, nil}
}

func (c *cursor) advance(segment []byte) {
	c.offset += len(segment)
	for _, ch := range string(segment) {
		c.charno++
		if ch == '\n' {
			c.lineno++
			c.charno = 0
		}
	}
}

func (c *cursor) peek(pattern *regexp.Regexp) *Node {
	source := c.source[c.offset:]
	loc := pattern.FindIndex(source)
	if loc == nil || loc[0] != 0 || loc[0] == loc[1] {
		return nil
	}
	segment := source[loc[0]:loc[1]]
	node := c.node(string(segment))
	c.advance(segment)
	return node
}

func (c *cursor) hasFinished() bool {
	return c.offset == len(c.source)
}

func (c *cursor) hasAdvenced() bool {
	ok := c.lastOffset != c.offset
	c.lastOffset = c.offset
	return ok
}

func (p *Parser) parseNoun(verb, noun *Node, c *cursor) (nextVerb, nextNoun *Node, err error) {
	// parse
	if noun == nil {
		c.peek(p.Space)
		noun = c.peek(p.Noun)
		if noun == nil {
			// epsilon
			noun = c.node("")
		} else if _, ok := p.Keywords[noun.Text]; ok {
			// keyword
			verb = noun
			noun = c.node("")
		}
	}
	// push
	err = c.pushNoun(noun)
	noun = nil
	return verb, noun, err
}

func (p *Parser) parseVerb(verb, noun *Node, c *cursor) (nextVerb, nextNoun *Node, err error) {
	// parse
	if verb == nil {
		c.peek(p.Space)
		verb = c.peek(p.Verb)
		if verb == nil {
			noun = c.peek(p.Noun)
			if noun == nil {
				// lambda
				verb = c.node("")
			} else if _, ok := p.Keywords[noun.Text]; ok {
				// keyword
				verb = noun
				noun = nil
			} else {
				// lambda
				verb = c.node("")
			}
		}
	}
	// push
	err = p.pushVerb(verb, noun, c)
	verb = nil
	return verb, noun, err
}

func (p *Parser) pushVerb(verb, noun *Node, c *cursor) (err error) {
	// closing bracket
	if openings, ok := p.Brackets.Openings[verb.Text]; ok {
		lastBracket := c.brackets[len(c.brackets)-1]
		if _, ok := openings[lastBracket.node.Text]; ok {
			c.popBracket()
		}
	}
	// verb
	priority, ok := p.Priorities[verb.Text]
	if !ok {
		return verb.NewError(ErrUnrecognizedVerb)
	}
	_, associativity := p.Associativities[verb.Text]
	err = c.pushVerb(verb, priority, associativity)
	if err != nil {
		return err
	}
	// opening bracket
	if _, ok := p.Brackets.Closings[verb.Text]; ok {
		err = c.pushBracket(verb)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) popAll(c *cursor) (*Node, error) {
	if len(c.brackets) > 1 {
		lastBracket := c.brackets[len(c.brackets)-1]
		return nil, lastBracket.node.NewError(ErrUnpairedBracket)
	}
	for len(c.verbs) > 0 {
		c.popVerb()
	}
	return c.nouns[0], nil
}

// Parse ...
func (p *Parser) Parse(source []byte, path string) (*Node, error) {
	c := &cursor{source: source, path: path, lastOffset: -1}
	err := c.pushBracket(c.node(""))
	if err != nil {
		return nil, err
	}
	var verb, noun *Node
	for {
		// parse noun
		verb, noun, err = p.parseNoun(verb, noun, c)
		if err != nil {
			return nil, err
		}
		if c.hasFinished() {
			break // reach the end
		}
		if !c.hasAdvenced() {
			// parser stocked
			return nil, c.node("").NewError(ErrNotMakingProgress)
		}
		// parse verb
		verb, noun, err = p.parseVerb(verb, noun, c)
		if err != nil {
			return nil, err
		}
	}
	return p.popAll(c)
}
