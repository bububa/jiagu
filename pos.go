package jiagu

import (
	"github.com/bububa/jiagu/perceptron"
	"github.com/bububa/jiagu/perceptron/model"
)

var posModel *perceptron.Perceptron

// PosModel get posModel singleton
func PosModel() *perceptron.Perceptron {
	if posModel == nil {
		var err error
		if posModel, err = initPerceptron(POS_MODEL); err != nil {
			panic(err)
		}
	}
	return posModel
}

// Pos 词性标注
func Pos(words []string) []model.Class {
	PosModel()
	return posModel.Predict(words)
}
