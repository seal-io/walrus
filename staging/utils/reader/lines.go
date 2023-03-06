package reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// Lines represents line reader in a given 'FileName' immediately before the given 'LineNum'.
// Extraction happens when 'Condition' is met and being processed by 'Parser' function.
type Lines struct {
	FileName  string
	LineNum   int
	Condition func(line string) bool
	Parser    func(line string) (string, bool)
}

// Read reads lines in given file and based on the provided condition.
func (l *Lines) Read() ([]string, error) {
	f, err := os.Open(l.FileName)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Size() == 0 {
		return []string{}, nil
	}
	defer func() {
		_ = f.Close()
	}()
	return l.extract(f)
}

func (l *Lines) extract(r io.Reader) ([]string, error) {
	bf := bufio.NewReader(r)
	var lines = make([]string, 0)
	for lnum := 0; ; lnum++ {
		if lnum >= l.LineNum-1 {
			break
		}
		line, err := bf.ReadString('\n')
		if errors.Is(err, io.EOF) && line == "" {
			switch lnum {
			case 0:
				return nil, errors.New("no lines in file")
			case 1:
				return nil, errors.New("only 1 line")
			default:
				return nil, fmt.Errorf("only %d lines", lnum)
			}
		}

		if l.Condition(line) {
			if extracted, capture := l.Parser(line); capture {
				lines = append(lines, extracted)
			}
		} else {
			lines = nil
		}
	}
	return lines, nil
}
