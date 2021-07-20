package jiagu

import (
	"os"
	"testing"
)

// TestFindword 测试发现新词
func TestFindword(t *testing.T) {
	t.Skip()
	inputFile := "./data/findword/input.txt"
	fd, err := os.Open(inputFile)
	if err != nil {
		t.Error(err)
		return
	}
	defer fd.Close()
	words, err := Findword(fd, 0, 0, 0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v\n", words)
}
