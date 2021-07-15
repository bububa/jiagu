package jiagu

import (
	"github.com/bububa/jiagu/perceptron"
	"github.com/bububa/jiagu/perceptron/model"
)

var nerModel *perceptron.Perceptron

// Ner 命名实体识别
func Ner(words []string) []model.Class {
	return nerModel.Predict(words)
}
