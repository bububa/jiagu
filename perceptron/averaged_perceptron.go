package perceptron

import (
	"encoding/gob"
	"io"
	"sort"

	"github.com/shopspring/decimal"

	"github.com/bububa/jiagu/perceptron/model"
)

// AveragedPerceptron Averaged Perceptron for POS tagging
type AveragedPerceptron struct {
	weights   *model.Weights // Each feature gets its own weight vector, so weights is a dict-of-dicts
	classes   map[string]struct{}
	totals    map[string]float64 // The accumulated values, for the averaging. These will be keyed by feature/clas tuples
	tstamps   map[string]float64 // The last time the feature was changed, for the averaging. Also keyed by feature/clas tuples
	instances float64            // Number of instances seen
}

// NewAveragedPerceptron 新建Perceptron
func NewAveragedPerceptron() *AveragedPerceptron {
	return &AveragedPerceptron{
		weights: model.NewWeights(0),
		classes: make(map[string]struct{}),
		totals:  make(map[string]float64),
		tstamps: make(map[string]float64),
	}
}

// NewAveragedPerceptronFromJSON 从JSON Model创建Perceptron
func NewAveragedPerceptronFromJSON(aModel model.PerceptronJSONModel) *AveragedPerceptron {
	return &AveragedPerceptron{
		weights: model.NewWeightsFromMap(aModel.Weights),
		classes: aModel.Classes,
		totals:  make(map[string]float64),
		tstamps: make(map[string]float64),
	}
}

// NewAveragedPerceptronFromReader 从io.Reader创建Perceptron
func NewAveragedPerceptronFromReader(r io.Reader) (*AveragedPerceptron, error) {
	var aModel model.PerceptronJSONModel
	err := gob.NewDecoder(r).Decode(&aModel)
	if err != nil {
		return nil, err
	}
	return &AveragedPerceptron{
		weights: model.NewWeightsFromMap(aModel.Weights),
		classes: aModel.Classes,
		totals:  make(map[string]float64),
		tstamps: make(map[string]float64),
	}, nil
}

func (p *AveragedPerceptron) Len() int {
	return len(p.classes)
}

// Predict Dot-product the features and current weights and return the best label.
func (p *AveragedPerceptron) Predict(features []model.Feature) model.Class {
	scores := make(map[string]float64, p.Len())
	for _, feature := range features {
		if feature.IsZero() {
			continue
		}
		classIter := p.weights.GetFeatureWeights(feature.Label)
		for kv := range classIter {
			if _, found := p.classes[kv.Label]; !found {
				continue
			}
			scores[kv.Label] += kv.Value * feature.Value
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
	key := model.FeatureClassKey(feat, clas)
	p.totals[key] += (p.instances - p.tstamps[key]) * weight
	p.tstamps[key] = p.instances
	p.weights.SetWeight(feat, clas, weight+value)
}

// AverageWeights Average weights from all iterations.
func (p *AveragedPerceptron) AverageWeights() {
	iter := p.weights.Features()
	for feat := range iter {
		classIter := p.weights.GetFeatureWeights(feat)
		for kv := range classIter {
			key := model.FeatureClassKey(feat, kv.Label)
			total := p.totals[key] + (p.instances-p.tstamps[key])*kv.Value
			averaged, _ := decimal.NewFromFloat(total / p.instances).Round(3).Float64()
			if averaged > 0.0001 {
				p.weights.SetWeight(feat, kv.Label, averaged)
			}
		}
	}
}

// AddClass add class
func (p *AveragedPerceptron) AddClass(clas string) {
	p.classes[clas] = struct{}{}
}

// Model
func (p AveragedPerceptron) Model() model.PerceptronJSONModel {
	return model.PerceptronJSONModel{
		Weights: p.weights.Map(),
		Classes: p.classes,
	}
}
