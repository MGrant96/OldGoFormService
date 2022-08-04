package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/google/uuid"
	"github.com/polyloop/formservice/data"
	"github.com/polyloop/formservice/graph/generated"
	"github.com/polyloop/formservice/graph/model"
	"github.com/polyloop/formservice/mailers"
)

func (r *mutationResolver) CreateForm(ctx context.Context, input model.NewForm) (*model.Form, error) {
	var participants []*model.Participant
	var questions []*model.Question

	for _, participant := range input.Participants {
		id := uuid.New().String()

		newParticipant := &model.Participant{
			ID:    id,
			Email: participant.Email,
		}
		participants = append(participants, newParticipant)
	}

	for _, question := range input.Questions {
		id := uuid.New().String()
		var options []*model.Option

		for _, option := range question.Options {
			newOption := &model.Option{
				Rank:        option.Rank,
				Description: option.Description,
			}
			options = append(options, newOption)
		}

		newQuestion := &model.Question{
			ID:        id,
			Title:     question.Title,
			InputType: question.InputType,
			Options:   options,
			Min:       question.Min,
			Max:       question.Max,
		}
		questions = append(questions, newQuestion)
	}

	formUUID := uuid.New().String()
	form := &model.Form{
		Title:        input.Title,
		ID:           formUUID,
		Participants: participants,
		Questions:    questions,
		TenderID:     input.TenderID,
		Description:  input.Description,
		Subject:      input.Subject,
		Message:      input.Message,
		// OrgID:        input.OrgID,
		AimID: input.AimID,
	}

	// r.forms = append(r.forms, form)
	data.CreateForm(form)
	for _, participant := range form.Participants {
		mailers.SendEmail(participant.ID, participant.Email, form.Message, form.Subject)
	}

	return form, nil
}

func (r *mutationResolver) DeleteForm(ctx context.Context, input model.DeleteForm) (*model.Form, error) {
	form := data.DeleteForm(input.ID)

	return form, nil
}

func (r *mutationResolver) AddFeedback(ctx context.Context, input []*model.NewFeedback) (*model.Form, error) {
	id := uuid.New().String()
	var feedbacks []*model.Feedback
	for _, feedback := range input {
		question := data.GetQuestion(feedback.QuestionID, feedback.FormID)
		newFeedback := &model.Feedback{
			ID:       id,
			Response: feedback.Response,
			Stage:    feedback.Stage,
			QuestionTitle: question.Title,
			QuestionType: question.InputType,
		}
		feedbacks = append(feedbacks, newFeedback)
		participant := data.GetParticipant(feedback.ParticipantID, feedback.FormID)

		updatedParticipant := &model.Participant{
			ID:        feedback.ParticipantID,
			Email:     participant.Email,
			Latitude:  &feedback.Latitude,
			Longitude: &feedback.Longitude,
		}

		data.UpdateParticipant(feedback.FormID, updatedParticipant)
	}

	form := data.AddFeedback(input[0].FormID, input[0].ParticipantID, feedbacks, input[0].QuestionID)
	return form, nil
}

func (r *mutationResolver) UpdateQuestion(ctx context.Context, input model.UpdateQuestion) (*model.Form, error) {
	var options []*model.Option

	for _, option := range input.Options {
		newOption := &model.Option{
			Rank:        option.Rank,
			Description: option.Description,
		}
		options = append(options, newOption)
	}

	newQuestion := &model.Question{
		ID:        input.ID,
		Title:     input.Title,
		InputType: input.InputType,
		Options:   options,
		Min:       input.Min,
		Max:       input.Max,
	}
	form := data.UpdateQuestion(input.FormID, newQuestion)
	return form, nil
}

func (r *mutationResolver) ResendEmail(ctx context.Context, input model.NewEmail) (string, error) {
	for _, participant := range input.Participants {
		mailers.SendEmail(participant.ParticipantID, participant.Email, input.Message, input.Subject)
	}
	return "Resent Email", nil
}

func (r *queryResolver) Forms(ctx context.Context) ([]*model.Form, error) {
	forms := data.GetForms()
	return forms, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
