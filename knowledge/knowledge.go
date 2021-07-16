package knowledge

import (
	"io"

	"github.com/bububa/jiagu/perceptron"
	pmodel "github.com/bububa/jiagu/perceptron/model"
	"github.com/bububa/jiagu/utils"
)

// Knowledge 知识图谱关系提取
type Knowledge struct {
	model *perceptron.Perceptron
}

// New 新建图谱关系
func New(model *perceptron.Perceptron) *Knowledge {
	return &Knowledge{
		model: model,
	}
}

// NewFromReader 从model新建
func NewFromReader(r io.Reader) (*Knowledge, error) {
	aModel, err := perceptron.NewFromReader(r)
	if err != nil {
		return nil, err
	}
	return &Knowledge{
		model: aModel,
	}, nil
}

func (k *Knowledge) Entities(words []string) []Entity {
	labels := k.model.Predict(words)
	return lab2spo(words, labels)
}

func lab2spo(words []string, eppLabels []pmodel.Class) []Entity {
	var (
		subjects []Subject
		objects  []Subject
		entities []Entity
	)
	var index int
	for idx, word := range words {
		ep := eppLabels[idx].Label
		epPrefix := ep[0]
		var label string
		if len(ep) > 2 {
			label = ep[2:]
		}
		if epPrefix == 'B' && label == "实体" {
			subjects = append(subjects, Subject{
				Word:  word,
				Label: label,
				Index: index,
			})
		} else if (epPrefix == 'I' || epPrefix == 'E') && label == "实体" {
			l := len(subjects)
			if l == 0 {
				continue
			}
			subjects[l-1].Word += word
		}

		if epPrefix == 'B' && label != "实体" {
			objects = append(objects, Subject{
				Word:  word,
				Label: label,
				Index: index,
			})
		} else if (epPrefix == 'I' || epPrefix == 'E') && label != "实体" {
			l := len(objects)
			if l == 0 {
				continue
			}
			objects[l-1].Word += word
		}
		index++
	}
	subL := len(subjects)
	objL := len(objects)
	if subL == 0 || objL == 0 {
		return entities
	}
	if subL == 1 {
		for _, obj := range objects {
			predicate := utils.StringInRange(obj.Label, 0, -1)
			entities = append(entities, Entity{
				Subject:   subjects[0].Word,
				Attribute: predicate,
				Value:     obj.Word,
			})
		}
	}
	return entities
}
