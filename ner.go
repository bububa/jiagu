package jiagu

import (
	"github.com/bububa/jiagu/perceptron"
	"github.com/bububa/jiagu/perceptron/model"
)

var nerModel *perceptron.Perceptron

// NerModel get nerModel singleton
func NerModel() *perceptron.Perceptron {
	if nerModel == nil {
		var err error
		if nerModel, err = initPerceptron(NER_MODEL); err != nil {
			panic(err)
		}
	}
	return nerModel
}

// Ner 命名实体识别
func Ner(words []string) []model.Class {
	NerModel()
	return nerModel.Predict(words)
}
