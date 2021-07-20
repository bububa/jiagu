package jiagu

import (
	"github.com/bububa/jiagu/perceptron"
	"github.com/bububa/jiagu/perceptron/model"
)

var posModel *perceptron.Perceptron

// Pos 词性标注
func Pos(words []string) []model.Class {
	if posModel == nil {
		var err error
		if posModel, err = initPerceptron(POS_MODEL); err != nil {
			panic(err)
		}
	}
	return posModel.Predict(words)
}
