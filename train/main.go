package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var (
		trainPath string
		testPath  string
		modelPath string
		iters     int
	)
	flag.StringVar(&trainPath, "train", "", "train data file")
	flag.StringVar(&testPath, "test", "", "text data file")
	flag.StringVar(&modelPath, "model", "", "model dir")
	flag.IntVar(&iters, "iters", 5, "iters")
	flag.Parse()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	if trainPath != "" {
		trainPath = filepath.Join(wd, trainPath)
		modelPath = filepath.Join(wd, modelPath)
		err := Train(trainPath, modelPath, iters)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if testPath != "" {
		testPath = filepath.Join(wd, testPath)
		modelPath = filepath.Join(wd, modelPath)
		precision, err := Eval(testPath, modelPath, true)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("\nPrecision: %f\n", precision)
	}
}
