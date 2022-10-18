package bencode

import (
	"bufio"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	inputBencodeStr := `xx`

	reader := bufio.NewReader(strings.NewReader(inputBencodeStr))
	beValue, err := decode(reader)
	if err != nil {
		t.Error(err)
	}
	if beValue == nil {
		t.Fatalf("expected value go nil")
	}
}
