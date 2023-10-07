package presentation

import (
	"context"
	"log"

	"github.com/klyngen/votomatic-3000/packages/backend/models"
	"github.com/klyngen/votomatic-3000/packages/backend/protoclient"
)

type grpcServer struct {
	configuration *protoclient.ConfigurationResponse
	poll          models.Poll
	listeners     []chan *protoclient.VoteUpdate
	protoclient.UnimplementedVoteServiceServer
}

// GetConfiguration implements protoclient.VoteServiceServer.
func (g *grpcServer) GetConfiguration(context.Context, *protoclient.EmptyRequest) (response *protoclient.ConfigurationResponse, err error) {
	response = g.configuration
	return
}

// GetVoteStatus implements protoclient.VoteServiceServer.
func (g *grpcServer) GetVoteStatus(context.Context, *protoclient.EmptyRequest) (*protoclient.VoteStatus, error) {
	questions := make([]*protoclient.Question, len(g.configuration.Questions))
	pollStatus := g.poll.GetAnswers()

	for i, ps := range pollStatus {
		questions[i] = &protoclient.Question{
			Id:    int32(i),
			Votes: ps,
		}
	}

	voteStatusReponse := &protoclient.VoteStatus{
		Questions: questions,
	}

	return voteStatusReponse, nil
}

func (g *grpcServer) addListener(listener chan *protoclient.VoteUpdate) {
	g.listeners = append(g.listeners, listener)
}

// GetVoteUpdates implements protoclient.VoteServiceServer.
func (g *grpcServer) GetVoteUpdates(r *protoclient.EmptyRequest, listener protoclient.VoteService_GetVoteUpdatesServer) error {
	updateChan := make(chan *protoclient.VoteUpdate)
	log.Println("Adding new vote update listener")
	g.addListener(updateChan)

	for {
		select {
		case update := <-updateChan:
			listener.Send(update)
			break

		case <-listener.Context().Done():
			log.Println("Client closed grpc stream connection")
			return nil
		}
	}

}

func (g *grpcServer) updateListeners(index int, alternative int) {
	update := &protoclient.VoteUpdate{
		QuestionId: int32(index),
		VoteIndex:  int32(alternative),
	}
	for _, vs := range g.listeners {
		vs <- update
	}
}

// PushVote implements protoclient.VoteServiceServer.
func (g *grpcServer) PushVote(ctx context.Context, request *protoclient.PushVoteRequest) (*protoclient.EmptyRequest, error) {
	g.poll.AddAnswer(models.PollAnswer{
		Id:          int(request.QuestionId),
		Alternative: int(request.VoteIndex),
	})

	go g.updateListeners(int(request.QuestionId), int(request.VoteIndex))

	return &protoclient.EmptyRequest{}, nil
}

func NewGrpcServer(configuration models.Configuration, poll models.Poll) protoclient.VoteServiceServer {
	configurationResponse := createPollConfiguration(configuration)
	return &grpcServer{
		configuration: configurationResponse,
		poll:          poll,
		listeners:     make([]chan *protoclient.VoteUpdate, 0),
	}
}

func createPollConfiguration(configuration models.Configuration) *protoclient.ConfigurationResponse {
	questions := make([]*protoclient.ConfigurationQuestion, len(configuration.Questions))

	for i, question := range configuration.Questions {
		questions[i] = &protoclient.ConfigurationQuestion{
			Description: question.Description,
		}
		questions[i].Alternatives = make([]*protoclient.ConfigurationQuestionAlternative, len(question.Alternatives))

		for j, alternative := range question.Alternatives {
			questions[i].Alternatives[j] = &protoclient.ConfigurationQuestionAlternative{
				Name:  alternative.Name,
				Emoji: alternative.Emoji,
			}
		}
	}

	configResponse := protoclient.ConfigurationResponse{
		Questions: questions,
	}

	return &configResponse
}
