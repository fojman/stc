package bencode

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
)

const (
	typeEndByte   = 'e'
	listStartByte = 'l'
	dictStartByte = 'd'
)

type encoder struct {
	reader *bufio.Reader
	pos    int
	len    int
}

func newEncoder(r *bufio.Reader) encoder {
	return encoder{
		reader: r,
		pos:    0,
		len:    r.Size(),
	}
}

// decode BE byte stream
func decode(r *bufio.Reader) (data interface{}, err error) {
	e := newEncoder(r)

	return e.decode()
}

func (e *encoder) decode() (interface{}, error) {

	var result interface{}
	var ch byte

	ch, err := e.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch ch {
	case 'i':
		numStr, err := e.readUntil(byte('e'))
		if err != nil {
			return nil, err
		}

		number, err := strconv.Atoi(string(numStr))
		if err != nil {
			return nil, err
		}
		result = number
	case 'l': // list -> l<any_type>e

		result, err = e.decode()     // recursive call -> TODO: better to use only functions func(r *bufio.Reader) so easier to use recursion
		lb, _ := e.reader.ReadByte() // todo: err ignored
		if lb != 'e' {
			return nil, fmt.Errorf("list: [%s] invalid termination symbol", lb)
		}
	case 'd': // dict -> d<any>e
		result, err = e.decode() // recursive call -> TODO: better to use only functions func(r *bufio.Reader) so easier to use recursion
		lb, _ := e.reader.ReadByte()
		if lb != 'e' {
			return nil, fmt.Errorf("dict: [%s] invalid termination symbol", lb)
		}
	default: // string
		// assert `ch` must be number followed by `:` 12:qwer...sdf
		if !isDigit(ch) {
			return nil, fmt.Errorf("str: number expected")
		}
		e.reader.UnreadByte()
		strLenStr, err := e.readUntil(byte(':'))
		if err != nil {
			return nil, err
		}
		len, err := strconv.Atoi(string(strLenStr))
		if err != nil {
			return nil, err
		}

		// `:` already skipped so we need just read len bytes from reader
		strBuff := make([]byte, len) // @evgen: is OK way to create buff?
		n, err := e.reader.Read(strBuff)
		if err != nil || n < len { // @evgen: should we use one cond here?
			return nil, err
		}
		result = string(strBuff)
	}

	return result, nil
}

// read from reader until delimiter
// returns error in case of delim not found or any other related to buff io
func (e *encoder) readUntil(delim byte) ([]byte, error) {
	canReadLen := e.reader.Buffered() // num of available bytes in reader

	var buf []byte
	buf, err := e.reader.Peek(canReadLen)
	if err != nil {
		return nil, err
	}

	if i := bytes.IndexByte(buf, delim); i >= 0 {
		return e.reader.ReadSlice(delim)
	}

	return e.reader.ReadBytes(delim)
}

func isDigit(ch byte) bool { return ch >= '0' && ch <= 9 }
