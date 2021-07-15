package perceptron

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/shopspring/decimal"

	"github.com/bububa/jiagu/perceptron/model"
)

// AveragedPerceptron Averaged Perceptron for POS tagging
type AveragedPerceptron struct {
	weights   model.Weights // Each feature gets its own weight vector, so weights is a dict-of-dicts
	classes   map[string]struct{}
	totals    map[string]float64 // The accumulated values, for the averaging. These will be keyed by feature/clas tuples
	tstamps   map[string]float64 // The last time the feature was changed, for the averaging. Also keyed by feature/clas tuples
	instances float64            // Number of instances seen
}

// NewAveragedPerceptron 新建Perceptron
func NewAveragedPerceptron() *AveragedPerceptron {
	return &AveragedPerceptron{
		weights: model.NewWeights(),
		classes: make(map[string]struct{}),
		totals:  make(map[string]float64),
		tstamps: make(map[string]float64),
	}
}

// NewAveragedPerceptronFromModel 从model创建Perceptron
func NewAveragedPerceptronFromModel(aModel *model.PerceptronModel) *AveragedPerceptron {
	return &AveragedPerceptron{
		weights: aModel.Weights,
		classes: aModel.Classes,
		totals:  make(map[string]float64),
		tstamps: make(map[string]float64),
	}
}

// NewAveragedPerceptronFromReader 从io.Reader创建Perceptron
func NewAveragedPerceptronFromReader(r io.Reader) (*AveragedPerceptron, error) {
	var aModel model.PerceptronModel
	err := json.NewDecoder(r).Decode(&aModel)
	if err != nil {
		return nil, err
	}
	return &AveragedPerceptron{
		weights: aModel.Weights,
		classes: aModel.Classes,
		totals:  make(map[string]float64),
		tstamps: make(map[string]float64),
	}, nil
}

// Predict Dot-product the features and current weights and return the best label.
func (p *AveragedPerceptron) Predict(features []model.Feature) model.Class {
	scores := make(map[string]float64, len(p.classes))
	for _, feature := range features {
		if feature.IsZero() {
			continue
		}
		classes := p.weights.GetFeatureWeights(feature.Label)
		if classes == nil {
			continue
		}
		for label, weight := range classes {
			if _, found := p.classes[label]; !found {
				continue
			}
			scores[label] += weight * feature.Value
		}
	}
	totalScores := len(scores)
	if totalScores == 0 {
		return model.Class{}
	}
	scoreSlice := model.NewKVSlice(totalScores)
	for label, value := range scores {
		scoreSlice = append(scoreSlice, model.KV{Label: label, Value: value})
	}
	sort.Sort(sort.Reverse(scoreSlice))
	return scoreSlice[0]
}

// Update update the feature weights.
func (p *AveragedPerceptron) Update(truth string, guess string, features []model.Feature) {
	p.instances += 1
	if truth == guess {
		return
	}
	var classIters = []model.KV{
		{
			Label: truth,
			Value: 1.0,
		},
		{
			Label: guess,
			Value: -1.0,
		},
	}
	for _, feature := range features {
		for _, iter := range classIters {
			weight := p.weights.GetWeight(feature.Label, iter.Label)
			p.updateFeature(iter.Label, feature.Label, weight, iter.Value)
		}
	}
}

func (p *AveragedPerceptron) updateFeature(clas string, feat string, weight float64, value float64) {
	key := featureClassKey(feat, clas)
	p.totals[key] += (p.instances - p.tstamps[key]) * weight
	p.tstamps[key] = p.instances
	p.weights.SetWeight(feat, clas, weight+value)
}

// AverageWeights Average weights from all iterations.
func (p *AveragedPerceptron) AverageWeights() {
	for feat, classes := range p.weights {
		for clas, weight := range classes {
			key := featureClassKey(feat, clas)
			total := p.totals[key] + (p.instances-p.tstamps[key])*weight
			averaged, _ := decimal.NewFromFloat(total / p.instances).Round(3).Float64()
			if averaged > 0.0001 {
				p.weights.SetWeight(feat, clas, averaged)
			}
		}
	}
}

// AddClass add class
func (p *AveragedPerceptron) AddClass(clas string) {
	p.classes[clas] = struct{}{}
}

// Model
func (p AveragedPerceptron) Model() model.PerceptronModel {
	return model.PerceptronModel{
		Weights: p.weights,
		Classes: p.classes,
	}
}

func featureClassKey(feat string, clas string) string {
	return fmt.Sprintf("%s-%s", feat, clas)
}
