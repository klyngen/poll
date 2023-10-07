package models

type Poll struct {
	pollStatus map[int][]int32
}

type pollAnswerStatus struct {
	id    int
	votes []int32
}

func (p pollAnswerStatus) Id() int {
	return p.id
}

func (p pollAnswerStatus) Votes() []int32 {
	return p.votes
}

type PollAnswer struct {
	Id          int
	Alternative int
}

func NewPoll(configuration Configuration) *Poll {
	poll := Poll{}
	poll.pollStatus = make(map[int][]int32)
	for i, q := range configuration.Questions {
		poll.pollStatus[i] = make([]int32, len(q.Alternatives))
		for j := range poll.pollStatus[i] {
			poll.pollStatus[i][j] = 1
		}
	}
	return &poll
}

func (p Poll) AddAnswer(answer PollAnswer) {
	p.pollStatus[answer.Id][answer.Alternative]++
}

func (p Poll) GetAnswers() map[int][]int32 {
	return p.pollStatus
}
