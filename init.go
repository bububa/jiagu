package jiagu

import (
	"embed"
	_ "embed"
	"fmt"

	"github.com/bububa/jiagu/knowledge"
	"github.com/bububa/jiagu/perceptron"
	"github.com/bububa/jiagu/segment"
	"github.com/bububa/jiagu/textrank"
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

func init() {
	var err error
	if posModel, err = initPerceptron(POS_MODEL); err != nil {
		panic(err)
	}
	if nerModel, err = initPerceptron(NER_MODEL); err != nil {
		panic(err)
	}
	initSegNroute()
	initKnowledge()
	initKeywords()
	initSummarize()
}

func initPerceptron(modelFile string) (*perceptron.Perceptron, error) {
	fd, err := modelFS.Open(fmt.Sprintf("model/%s", modelFile))
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return perceptron.NewFromReader(fd)
}

func initSegNroute() {
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
	segNroute, err = segment.NewFromReader(vocabR, modelR)
	if err != nil {
		panic(err)
	}
}

func initKnowledge() {
	modelR, err := modelFS.Open(fmt.Sprintf("model/%s", KG_MODEL))
	if err != nil {
		panic(err)
	}
	defer modelR.Close()
	knowledgeModel, err = knowledge.NewFromReader(modelR)
	if err != nil {
		panic(err)
	}
}

func initKeywords() {
	stopwordsR, err := dictFS.Open(fmt.Sprintf("dict/%s", STOPWORDS_DICT))
	if err != nil {
		panic(err)
	}
	defer stopwordsR.Close()
	keywordsModel = textrank.NewKeywords(segNroute)
	keywordsModel.LoadStopwords(stopwordsR)
}

func initSummarize() {
	stopwordsR, err := dictFS.Open(fmt.Sprintf("dict/%s", STOPWORDS_DICT))
	if err != nil {
		panic(err)
	}
	defer stopwordsR.Close()
	summarizeModel = textrank.NewSummarize(segNroute)
	summarizeModel.LoadStopwords(stopwordsR)
}
