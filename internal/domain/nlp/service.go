package nlp

import (
	"fmt"
	"os"
	"strings"

	"github.com/jbrukh/bayesian"
)

type Class = bayesian.Class

const (
	Good Class = "Good"
	Bad  Class = "Bad"
)

type NLPService struct {
	modelPath string
}

func NewNLPService(modelPath string) *NLPService {
	return &NLPService{modelPath: modelPath}
}

func (s *NLPService) TrainModel(goodStuff, badStuff []string) error {
	classifier := bayesian.NewClassifierTfIdf(Good, Bad)

	for _, goodWords := range goodStuff {
		goodData := s.parseWords(goodWords)
		classifier.Learn(goodData, Good)
	}

	for _, badWords := range badStuff {
		badData := s.parseWords(badWords)
		classifier.Learn(badData, Bad)
	}

	classifier.ConvertTermsFreqToTfIdf()

	if err := classifier.WriteToFile(s.modelPath); err != nil {
		return fmt.Errorf("failed to write model to %s: %w", s.modelPath, err)
	}

	return nil
}

func (s *NLPService) ClassifyContent(text string) (bayesian.Class, error) {
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("text to classify cannot be empty")
	}

	if !s.modelExists() {
		return "", fmt.Errorf("model file not found at '%s'", s.modelPath)
	}

	classifier, err := bayesian.NewClassifierFromFile(s.modelPath)
	if err != nil {
		return "", fmt.Errorf("failed to load classifier from model file '%s': %w", s.modelPath, err)
	}

	words := s.parseWords(text)

	_, score, _ := classifier.LogScores(words)
	if score < 0 || score >= len(classifier.Classes) {
		return "", fmt.Errorf("invalid classification score: %d", score)
	}

	return classifier.Classes[score], nil
}

func (s *NLPService) modelExists() bool {
	if _, err := os.Stat(s.modelPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (s *NLPService) parseWords(phrase string) []string {
	if strings.TrimSpace(phrase) == "" {
		return []string{}
	}

	fields := strings.Fields(phrase)
	words := make([]string, 0, len(fields))

	for _, word := range fields {
		cleaned := strings.TrimSpace(strings.ToLower(word))
		if cleaned != "" {
			words = append(words, cleaned)
		}
	}

	return words
}
