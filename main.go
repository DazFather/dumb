package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/DazFather/brush"
)

var (
	spacer       string
	handleOutput = func(fpath string, content []byte) error {
		return os.WriteFile(fpath, content, 0666)
	}
)

func init() {
	flag.BoolVar(&brush.Disable, "no-color", false, "force disable of colored output'")
	flag.BoolVar(&brush.Disable, "nc", false, "shorthand for 'no-color'")

	flag.StringVar(&spacer, "spacer", "\t", "unit of the indentation spacer")
	flag.StringVar(&spacer, "s", "\t", "shorthand for spacer")

	endln := flag.String("eol", eol, "line ending characters or LF/CRLF. Defaults to the OS")

	parseOutputFlag := func(val string) error {
		if val != "" {
			val = filepath.Clean(val)
			if err := os.MkdirAll(val, 0755); err != nil {
				return err
			}
			handleOutput = func(fpath string, content []byte) error {
				fpath = filepath.Join(val, fpath)
				err := os.MkdirAll(filepath.Dir(fpath), 0755)
				if err == nil {
					err = os.WriteFile(fpath, content, 0666)
				}
				return err
			}
		}
		return nil
	}
	flag.Func("output", "specify another output directory for the files (leave blank for overwite)", parseOutputFlag)
	flag.Func("o", "shorthand for output", parseOutputFlag)

	flag.Parse()

	switch strings.ToLower(*endln) {
	case "lf":
		eol = "\n"
	case "crlf":
		eol = "\r\n"
	default:
		eol = *endln
	}
}

/*
	v asdpkapsdkpaskdpad
	! sdladlmasld/saddasdas
		-asdasdadad
		-saddadada
		-adadsadad
	x asd/asd/asdad
*/

func main() {
	var (
		files, err = loadFiles()
		wg         sync.WaitGroup
	)

	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
		return
	}

	for fpath := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()

			logs := []string{}

			f, err := os.Open(fpath)
			if err != nil {
				fmt.Println(collect(fpath, danger(err)))
				return
			}
			defer f.Close()

			if err = handleOutput(fpath, []byte(indent(f, &logs))); err != nil {
				logs = append(logs, danger(err))
			}
			fmt.Println(collect(fpath, logs...))
		}()
	}
	wg.Wait()
}

func loadFiles() (<-chan string, error) {
	var rawpaths []string

	if args := flag.Args(); len(args) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		rawpaths = []string{dir}
	} else {
		rawpaths = args
	}

	var files = make(chan string)

	go func() {
		defer close(files)

		for _, fpath := range rawpaths {
			fpath = filepath.Clean(fpath)

			info, err := os.Stat(fpath)
			if err != nil {
				return
			}

			if !info.IsDir() {
				files <- fpath
				continue
			}

			// TODO: Handle error
			err = filepath.WalkDir(fpath, func(path string, info fs.DirEntry, e error) error {
				if e != nil {
					return e
				}

				if !info.IsDir() {
					files <- path
				}
				return nil
			})
		}
	}()

	return files, nil
}

func indent(rd io.Reader, logs *[]string) string {
	var (
		scanner  = bufio.NewScanner(rd)
		toclose  rune
		brackets queue[rune]
		parsed   = NewTree()
		txt      strings.Builder
	)

	for ln := 1; scanner.Scan(); ln++ {
		closing, escaped, inlineclosing := false, false, queue[rune]{}

		for _, ch := range strings.TrimSpace(scanner.Text()) {
			txt.WriteRune(ch)
			switch ch {
			case '\\':
				escaped = true
			case '(':
				if escaped {
					escaped = false
					continue
				}
				toclose = ch + 1
			case '{', '[':
				if escaped {
					escaped = false
					continue
				}
				toclose = ch + 2
			case toclose:
				if escaped {
					escaped = false
					continue
				}
				if item := brackets.pop(); item != nil {
					if size := len(brackets); size > 0 {
						toclose = brackets[size-1]
					} else {
						toclose = 0
					}
					closing = closing || inlineclosing.pop() == nil
				}
				continue
			case ')', ']', '}':
				if escaped {
					escaped = false
					continue
				}
				msg := ""
				if toclose == 0 {
					msg = fmt.Sprintf("Closing unopened bracket '%c' at line: %d\n", ch, ln)
				} else {
					msg = fmt.Sprintf("Mismatch bracket, closing '%c' but expected '%c', at line: %d\n", ch, toclose, ln)
				}
				line := txt.String()
				*logs = append(*logs, warn(msg+caret(line, txt.Len())))
				continue
			default:
				escaped = false
				continue
			}
			brackets.push(toclose)
			inlineclosing.push(ch)
		}

		if len(inlineclosing) > 0 {
			if closing {
				parsed.close("")
			}
			parsed.open(txt.String())
		} else if closing {
			parsed.close(txt.String())
		} else {
			parsed.add(txt.String())
		}
		txt.Reset()
	}

	if err := scanner.Err(); err != nil {
		*logs = append(*logs, danger(err))
		return ""
	}

	if len(brackets) > 0 {
		*logs = append(*logs, warn("Unclosed brackets: "+string(brackets)))
	}

	return parsed.Root().Indent(0, spacer)
}
