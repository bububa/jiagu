package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"

	"github.com/bububa/jiagu/perceptron"
	"github.com/bububa/jiagu/perceptron/model"
)

func Eval(testPath string, modelPath string, showProgressBar bool) (float64, error) {
	tagger, err := perceptron.NewFromModelFile(modelPath)
	if err != nil {
		return 0, err
	}
	var (
		bar      *progressbar.ProgressBar
		correct  float64
		total    float64
		sentence model.Sentence
	)
	if showProgressBar {
		bar = progressbar.NewOptions64(-1,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(false),
			progressbar.OptionSetWidth(15),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionSetDescription("Testing model..."),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
			progressbar.OptionOnCompletion(func() {
				fmt.Fprint(os.Stderr, "\n")
			}),
		)
	}
	fd, err := os.Open(testPath)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	buf := bufio.NewReader(fd)
	for {
		line, err := buf.ReadString('\n')
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return 0, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			words := sentence.Words
			tags := sentence.Tags
			outputs := tagger.Predict(words)
			for idx, tag := range tags {
				if tag == outputs[idx].Label {
					correct += 1
				}
				total += 1
			}
			bar.Add(1)
			continue
		}
		params := strings.Split(line, "	")
		if len(params) != 2 {
			log.Printf("params: %+v, l: %d\n", params, len(params))
			continue
		}
		sentence.Words = append(sentence.Words, params[0])
		sentence.Tags = append(sentence.Tags, params[1])
	}
	return correct / total, nil
}
