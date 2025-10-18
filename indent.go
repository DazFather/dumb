package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type indenter struct {
	toclose       rune
	brackets      queue[rune]
	inlineclosing int
	closing       bool
}

type indentFn func(ch rune, at int) (indentFn, error)

func (i *indenter) text(ch rune, at int) (indentFn, error) {
	switch ch {
	case i.toclose:
		if i.inlineclosing == 0 {
			i.closing = true
		} else {
			i.inlineclosing--
		}

		i.brackets.pop()
		if size := len(i.brackets); size > 0 {
			switch i.toclose = i.brackets[size-1]; i.toclose {
			case '\'', '"', '`':
				return i.text, nil
			}
		} else {
			i.toclose = 0
		}
		return i.char, nil
	case '\\':
		return func(nextCh rune, nextAt int) (indentFn, error) {
			return i.text, nil
		}, nil
	}

	return i.text, nil
}

func (i *indenter) slash(ch rune, at int) (indentFn, error) {
	switch ch {
	case '/':
		return i.comment, nil
	case '*':
		return i.multilineComment, nil
	}
	return i.char(ch, at)
}

func (i *indenter) comment(ch rune, at int) (indentFn, error) {
	if at == 0 {
		return i.char(ch, at)
	}
	return i.comment, nil
}

func (i *indenter) multilineComment(ch rune, at int) (indentFn, error) {
	if ch != '*' {
		return i.multilineComment, nil
	}

	return func(nextCh rune, nextAt int) (indentFn, error) {
		if nextCh == '/' && nextAt == at+1 {
			return i.char, nil
		}
		return i.multilineComment, nil
	}, nil
}

func (i *indenter) char(ch rune, at int) (indentFn, error) {
	switch ch {
	case i.toclose:
		if i.inlineclosing == 0 {
			i.closing = true
		} else {
			i.inlineclosing--
		}

		i.brackets.pop()
		if size := len(i.brackets); size > 0 {
			switch i.toclose = i.brackets[size-1]; i.toclose {
			case '\'', '"', '`':
				return i.text, nil
			}
		} else {
			i.toclose = 0
		}
	case '\\':
		return func(ch rune, at int) (indentFn, error) {
			return i.char, nil
		}, nil
	case '/':
		return i.slash, nil
	case '#':
		return i.comment, nil
	case '(':
		i.toclose = ch + 1
		i.brackets.push(i.toclose)
		i.inlineclosing++
	case '[', '{':
		i.toclose = ch + 2
		i.brackets.push(i.toclose)
		i.inlineclosing++
	case ')', ']', '}':
		if i.toclose == 0 {
			return i.char, fmt.Errorf("Closing unopened bracket '%c'", ch)
		}
		return i.char, fmt.Errorf("Mismatch bracket, closing '%c' but expected '%c'", ch, i.toclose)
	case '\'', '"', '`':
		i.toclose = ch
		i.brackets.push(i.toclose)
		i.inlineclosing++
		return i.text, nil
	}

	return i.char, nil
}

func Indent(rd io.Reader, logs *[]string) string {
	var (
		err     error
		scanner = bufio.NewScanner(rd)
		parsed  = NewTree()
		parser  indenter
	)

	for action, ln := parser.char, 1; scanner.Scan(); ln++ {
		line := strings.TrimSpace(scanner.Text())

		for at, ch := range line {
			if action, err = action(ch, at); err != nil {
				*logs = append(*logs, warn(fmt.Sprintln(err, "at line:", ln)+caret(line, at)))
			}
		}

		if parser.inlineclosing > 0 {
			if parser.closing {
				parsed.close("")
			}
			parsed.open(line)
		} else if parser.closing {
			parsed.close(line)
		} else if selfTab.MatchString(line) {
			parsed.add(spacer + line)
		} else {
			parsed.add(line)
		}

		parser.inlineclosing, parser.closing = 0, false
	}

	return parsed.Root().Indent(0, spacer)
}
