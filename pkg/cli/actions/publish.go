package pkg_action

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Minnek-Digital-Studio/cominnek/config"
	git_controller "github.com/Minnek-Digital-Studio/cominnek/controllers/git"
	"github.com/Minnek-Digital-Studio/cominnek/controllers/loading"
	"github.com/Minnek-Digital-Studio/cominnek/pkg/ask"
	"github.com/Minnek-Digital-Studio/cominnek/pkg/git"
	"github.com/Minnek-Digital-Studio/cominnek/pkg/github"
)

func publishQuestions() {
	if config.AppData.Publish.Ticket == "" {
		loading.Start("Reading branches...")
		ticket := git_controller.GetTicketNumber()
		loading.Stop()

		if ticket != "" {
			config.AppData.Publish.Ticket = ticket
			return
		}

		ask.One(&survey.Input{
			Message: "Enter the ticket number:",
		}, &config.AppData.Publish.Ticket, survey.WithValidator(survey.Required))
	}
}

func Publish() {
	if git_controller.CheckChanges() {
		Commit(false)
	}

	pushQuestion()
	publishQuestions()

	if git_controller.CheckChanges() {
		executeCommit()
	}

	github.Publish(config.AppData.Publish.Ticket)

	if config.AppData.Push.Merge != "" {
		git.Merge(config.AppData.Push.Merge)
	}
}
