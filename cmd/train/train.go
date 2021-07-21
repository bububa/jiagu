package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/bububa/jiagu/perceptron"
	"github.com/bububa/jiagu/perceptron/model"
)

func Train(trainPath string, modelPath string, iters int) error {
	tagger := perceptron.New()
	var trainData []model.Sentence
	var sentence model.Sentence
	fd, err := os.Open(trainPath)
	if err != nil {
		return err
	}
	defer fd.Close()
	buf := bufio.NewReader(fd)
	for {
		line, err := buf.ReadString('\n')
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if len(sentence.Tags) == 0 {
				continue
			}
			trainData = append(trainData, sentence)
			sentence = model.Sentence{}
			continue
		}
		params := strings.Split(line, "	")
		if len(params) != 2 {
			continue
		}
		sentence.Words = append(sentence.Words, params[0])
		sentence.Tags = append(sentence.Tags, params[1])
	}
	tagger.Train(trainData, iters, false, true)
	return tagger.Save(modelPath)
}
