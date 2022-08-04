package data

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/polyloop/formservice/graph/model"
)

func getTable() (dynamo.Table) {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String("eu-west-1")})
	table := db.Table("Form")

	return table
}

func CreateForm(form *model.Form) {
	table := getTable()
	preForm := "PRE Form"
	form.Stage = &preForm
	table.Put(form).Run()
}

func GetForms() ([]*model.Form) {
	var results []*model.Form
	table := getTable()
	table.Scan().All(&results)

	return results
}

func GetForm(formID string) (*model.Form, dynamo.Table) {
	var result *model.Form
	table := getTable()
	table.Get("ID", formID).One(&result)

	return result, table
}

func GetQuestion(questionID string, formID string) (*model.Question) {
	var result *model.Question
	form, _ := GetForm(formID)

	for _, question := range form.Questions {
		if (questionID == question.ID) {
			result = question
		}
	}

	return result
}

func GetParticipant(participantID string, formID string) (*model.Participant) {
	var result *model.Participant
	form, _ := GetForm(formID)

	for _, participant := range form.Participants {
		if (participantID == participant.ID) {
			result = participant
			break
		}
	}

	return result
}

func AddFeedback(formID string, participantID string, feedback []*model.Feedback, questionID string) (*model.Form) {
	form, table := GetForm(formID)
	for _, participant := range form.Participants {
		if (participantID == participant.ID) {
			participant.Feedback = append(participant.Feedback, feedback...)
			break
		}
	}

	for _, question := range form.Questions {
		if (questionID == question.ID) {
			question.Feedback = append(question.Feedback, feedback...)
			break
		}
	}
	postForm := "POST Form"
	form.Stage = &postForm
	table.Put(form).Run()
	return form
}

func UpdateQuestion(formID string, updatedQuestion *model.Question) (*model.Form) {
	form, table := GetForm(formID)
	for i, question := range form.Questions {
		if (updatedQuestion.ID == question.ID) {
			form.Questions[i] = updatedQuestion
			break
		}
	}
	table.Put(form).Run()
	return form
}

func UpdateParticipant(formID string, updatedParticipant *model.Participant) {
	form, table := GetForm(formID)
	for i, participant := range form.Participants {
		if (updatedParticipant.ID == participant.ID) {
			form.Participants[i] = updatedParticipant
			break
		}
	}
	table.Put(form).Run()
}

func DeleteForm(formID string) (*model.Form) {
	form, table := GetForm(formID)
	err := table.Delete("ID", formID).Run()
	fmt.Println("Error: ", err)

	return form
}
