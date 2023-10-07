package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var emojis = "💩🐂🐈🐕🌭👿🥚🦇🦆🔔☁🦉🤔"

func getEmojiByIndex(index int) byte {
	emojiLength := len(emojis)

	wrappedIndex := index % emojiLength

	return emojis[wrappedIndex]
}

func createConfiguration(size int) Configuration {
	configuration := Configuration{}

	for i := 0; i < size; i++ {
		question := Question{
			Description: string(getEmojiByIndex(i)),
		}

		for j := 0; j < 10; j++ {
			question.Alternatives = append(question.Alternatives, Alternative{
				Name:  string(getEmojiByIndex(j)),
				Emoji: string(getEmojiByIndex(j)),
			})
		}
		configuration.Questions = append(configuration.Questions, question)
	}

	return configuration
}

func TestNewPoll(t *testing.T) {
	size := 10
	config := createConfiguration(size)
	assert.True(t, config.IsValidConfiguration(), "Configuration was not valid")

	poll := NewPoll(config)

	assert.Equal(t, size, len(poll.pollStatus))
}

func Test_poll_AddAnswer(t *testing.T) {
	size := 10
	config := createConfiguration(size)

	poll := NewPoll(config)

	for i := 0; i < 10; i++ {
		poll.AddAnswer(PollAnswer{
			Id:          i,
			Alternative: i,
		})
	}

	assert.Equal(t, int32(1), poll.pollStatus[1][1], "There should be only one vote for the value at this index")
	assert.Equal(t, int32(0), poll.pollStatus[2][1], "There should be no votes")
}

func BenchmarkAddAnswer(b *testing.B) {
	size := 10
	config := createConfiguration(size)
	poll := NewPoll(config)

	for i := 0; i < b.N; i++ {
		id := i % size
		poll.AddAnswer(PollAnswer{
			Id:          id,
			Alternative: id,
		})

		poll.AddAnswer(PollAnswer{
			Id:          id,
			Alternative: id,
		})

		poll.AddAnswer(PollAnswer{
			Id:          id,
			Alternative: id,
		})

		poll.AddAnswer(PollAnswer{
			Id:          id,
			Alternative: id,
		})
	}
}
