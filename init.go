package jiagu

import (
	"embed"
	_ "embed"
	"fmt"

	"github.com/bububa/jiagu/perceptron"
)

const (
	// POS_MODEL 词性标注model文件
	POS_MODEL = "pos.model"
	// NER_MODEL 命名实体识别model文件
	NER_MODEL = "ner.model"
	// KG_MODEL 知识图谱model文件
	KG_MODEL = "kg.model"
	// CWS_MODEL 分词model文件
	CWS_MODEL = "cws.model"
	// VOCAB_DICT 分词字典
	VOCAB_DICT = "jiagu.dict"
	// STOPWORDS stopwords字典
	STOPWORDS_DICT = "stopwords.txt"
)

//go:embed model/*
var modelFS embed.FS

//go:embed dict/*
var dictFS embed.FS

func initPerceptron(modelFile string) (*perceptron.Perceptron, error) {
	fd, err := modelFS.Open(fmt.Sprintf("model/%s", modelFile))
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return perceptron.NewFromReader(fd)
}
