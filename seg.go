package jiagu

import (
	"fmt"
	"io"

	"github.com/bububa/jiagu/segment"
)

var seg *segment.Segment

// Segment get seg singleton
func Segment() *segment.Segment {
	if seg == nil {
		vocabR, err := dictFS.Open(fmt.Sprintf("dict/%s", VOCAB_DICT))
		if err != nil {
			panic(err)
		}
		defer vocabR.Close()
		modelR, err := modelFS.Open(fmt.Sprintf("model/%s", CWS_MODEL))
		if err != nil {
			panic(err)
		}
		defer modelR.Close()
		seg, err = segment.NewFromReader(vocabR, modelR)
		if err != nil {
			panic(err)
		}
	}
	return seg
}

// Seg 分词
func Seg(sentence string) []string {
	model := Segment()
	return model.Seg(sentence, segment.Default_SegMode)
}

// LoadUserDict 加载用户词典
func LoadUserDict(r io.Reader) {
	model := Segment()
	model.LoadUserDict(r)
}

// AddVocabs 添加用户词典word
func AddVocabs(words []string) {
	model := Segment()
	for _, w := range words {
		model.AddVocab(w, 1)
	}
}

// DelVocabs 删除用户词典word
func DelVocabs(words []string) {
	model := Segment()
	for _, w := range words {
		model.DelVocab(w, 1)
	}
}
