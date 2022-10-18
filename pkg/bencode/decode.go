package bencode

import "io"

// module for decoding bencode

type bdata interface{}

type context struct {
	io  io.Reader
	pos int
}

func (self *context) decode(s string) bdata {

}
