package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/DazFather/brush"
)

var (
	selfTab      = regexp.MustCompile(`^[\.+\-]|[^\.+\-][\.+\-]$`)
	spacer       string
	handleOutput = func(fpath string, content []byte) error {
		return os.WriteFile(fpath, content, 0666)
	}
)

func init() {
	flag.BoolVar(&brush.Disable, "no-color", false, "force disable of colored output")
	flag.BoolVar(&brush.Disable, "nc", false, "shorthand for 'no-color'")

	flag.StringVar(&spacer, "spacer", "\t", "unit of the indentation spacer")
	flag.StringVar(&spacer, "s", "\t", "shorthand for spacer")

	endln := flag.String("eol", eol, "line ending characters or LF/CRLF. Defaults to the OS")

	flag.Func("indent", "regexp pattern for indenting a single line", func(val string) error {
		var err error
		selfTab, err = regexp.Compile(val)
		return err
	})

	parseOutputFlag := func(val string) error {
		switch val {
		case "":
			// skip
		case "-":
			handleOutput = func(fpath string, content []byte) error {
				_, err := fmt.Printf("%s\n\n%s\n", magenta.Paint(fpath), content)
				return err
			}
		default:
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

	flag.Func("output", "specify another output directory for the files (leave blank for overwite, - for stdout)", parseOutputFlag)
	flag.Func("o", "shorthand for output", parseOutputFlag)

	echo := flag.Bool("echo", false, "print output to stdout instead of writing to files")

	flag.Parse()

	switch strings.ToLower(*endln) {
	case "lf":
		eol = "\n"
	case "crlf":
		eol = "\r\n"
	default:
		eol = *endln
	}

	if *echo {
		parseOutputFlag("-")
	}
}

func main() {
	var (
		files, err = loadFiles()
		wg         sync.WaitGroup
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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

			if err = handleOutput(fpath, []byte(Indent(f, &logs))); err != nil {
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

			filepath.WalkDir(fpath, func(path string, info fs.DirEntry, e error) error {
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
