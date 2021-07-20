package stopwords

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/bububa/jiagu/utils"
)

// Stopwords 停用词
type Stopwords struct {
	utils.StringSet
}

// New 新建停用词
func New() *Stopwords {
	return &Stopwords{utils.InitStringSet()}
}

// Load 从io.Reader 加载stopwords
func (s *Stopwords) Load(r io.Reader) error {
	buf := bufio.NewReader(r)
	var kws []string
	for {
		kw, err := buf.ReadString('\n')
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		kw = strings.TrimSpace(kw)
		kws = append(kws, kw)
	}
	s.Add(kws)
	return nil
}

// LoadFile 加载stopwords文件
func (s *Stopwords) LoadFile(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	buf := bufio.NewReader(fd)
	var kws []string
	for {
		kw, err := buf.ReadString('\n')
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		kw = strings.TrimSpace(kw)
		kws = append(kws, kw)
	}
	s.Add(kws)
	return nil
}
