package main

import (
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/bububa/jiagu/classify/bayes"
	"github.com/bububa/jiagu/perceptron/model"
)

func main() {
	var (
		inputFile      string
		outputFile     string
		sentimentModel bool
	)
	flag.StringVar(&inputFile, "i", "", "json model file")
	flag.StringVar(&outputFile, "o", "", "output model file")
	flag.BoolVar(&sentimentModel, "sentiment", false, "convert sentiment model")
	flag.Parse()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	inputFile = filepath.Join(wd, inputFile)
	outputFile = filepath.Join(wd, outputFile)
	log.Printf("converting: %s -> %s \n", inputFile, outputFile)

	fd, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer fd.Close()

	oFd, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer oFd.Close()
	gw := gzip.NewWriter(oFd)
	defer gw.Close()
	/*
		zipWriter := zip.NewWriter(oFd)
		defer zipWriter.Close()

		fInfo, err := fd.Stat()
		if err != nil {
			log.Fatalln(err)
		}
		zipFile, err := zipWriter.Create(fInfo.Name())
		if err != nil {
			log.Fatalln(err)
		}
	*/
	if sentimentModel {
		var jsonModel bayes.Model
		err = json.NewDecoder(fd).Decode(&jsonModel)
		if err != nil {
			log.Fatalln(err)
		}
		err = gob.NewEncoder(gw).Encode(jsonModel)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		var jsonModel model.PerceptronJSONModel
		err = json.NewDecoder(fd).Decode(&jsonModel)
		if err != nil {
			log.Fatalln(err)
		}
		err = gob.NewEncoder(gw).Encode(jsonModel)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Printf("converted: %s -> %s \n", inputFile, outputFile)
}
