package models

type Configuration struct {
	Questions []Question `json:"questions"`
}

func (c Configuration) GetQuestions() []Question {
	questions := make([]Question, len(c.Questions))
	for index, question := range c.Questions {
		questions[index] = question
		questions[index].Id = index
	}

	return questions
}

type Question struct {
	Id           int           `json:"id"`
	Description  string        `json:"description"`
	Alternatives []Alternative `json:"alternatives"`
}

type Alternative struct {
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}

func (c Configuration) IsValidConfiguration() bool {
	for _, question := range c.Questions {
		if !question.isValidQuestion() {
			return false
		}
	}
	return true
}

func (q Question) isValidQuestion() bool {
	for _, a := range q.Alternatives {
		if !a.isValidAlternative() {
			return false
		}
	}
	return true
}

func (a Alternative) isValidAlternative() bool {
	return len(a.Emoji) > 0 && len(a.Name) > 0
}
