package main

import "strings"

type blocktype uint8

const (
	root blocktype = iota
	literal
	scope
)

type block struct {
	blocktype
	Body     []string
	Open     string
	Contents []*block
	Close    string
}

func (b block) Indent(lv int, space string) string {
	var txt strings.Builder

	if b.blocktype == literal && len(b.Body) > 0 {
		for _, line := range b.Body {
			txt.WriteString(strings.Repeat(space, lv))
			txt.WriteString(line)
			txt.WriteString(eol)
		}
		return txt.String()
	}

	if len(b.Contents) == 0 {
		txt.WriteString(strings.Repeat(space, lv))
		txt.WriteString(b.Open)
		txt.WriteString(b.Close)
		txt.WriteString(eol)
		return txt.String()
	}

	if b.Open != "" {
		txt.WriteString(strings.Repeat(space, lv))
		txt.WriteString(b.Open)
		txt.WriteString(eol)
		lv++
	}

	for _, c := range b.Contents {
		txt.WriteString(c.Indent(lv, space))
	}

	if b.Open != "" {
		lv--
	}

	if b.Close != "" {
		txt.WriteString(strings.Repeat(space, lv))
		txt.WriteString(b.Close)
		txt.WriteString(eol)
	}

	return txt.String()
}

func (b *block) Append(txt string) {
	if b.blocktype == literal {
		b.Body = append(b.Body, txt)
		return
	}
	b.Contents = append(b.Contents, &block{
		Body:      []string{txt},
		blocktype: literal,
	})
}

type tree struct {
	root *block
	Prev []*block
	curr int
}

func NewTree() tree {
	root := new(block)
	return tree{curr: 0, Prev: []*block{root}, root: root}
}

func (t *tree) add(txt string) {
	if t.curr == 0 {
		t.open("")
	}
	t.Prev[t.curr].Append(txt)
}

func (t *tree) open(openstr string) {
	var item = block{Open: openstr, blocktype: scope}

	t.Prev[t.curr].Contents = append(t.Prev[t.curr].Contents, &item)
	t.Prev = append(t.Prev, &item)
	t.curr++
}

func (t *tree) close(closestr string) {
	t.Prev[t.curr].Close = closestr
	t.Prev = t.Prev[:t.curr]
	t.curr--
}

func (t tree) Root() (b block) {
	return *t.root
}
