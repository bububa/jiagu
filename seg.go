package jiagu

import (
	"io"

	"github.com/bububa/jiagu/segment"
)

var segNroute *segment.Segment

// Seg 分词
func Seg(sentence string, mode segment.SegMode) []string {
	return segNroute.Seg(sentence, mode)
}

// LoadUserDict 加载用户词典
func LoadUserDict(r io.Reader) {
	segNroute.LoadUserDict(r)
}

// AddVocabs 添加用户词典word
func AddVocabs(words []string) {
	for _, w := range words {
		segNroute.AddVocab(w, 1)
	}
}

// DelVocabs 删除用户词典word
func DelVocabs(words []string) {
	for _, w := range words {
		segNroute.DelVocab(w, 1)
	}
}
