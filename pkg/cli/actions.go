package cli

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Minnek-Digital-Studio/cominnek/config"
	"github.com/Minnek-Digital-Studio/cominnek/helper"
	"github.com/Minnek-Digital-Studio/cominnek/pkg/ask"
	pkg_action "github.com/Minnek-Digital-Studio/cominnek/pkg/cli/actions"
)

func askActions() {
	actions := []string{
		"Commit",
		"Create Branch",
		"Merge",
		"Pull Request",
		"Publish",
		"Push",
		"Release",
		"Stash",
		"Reset",
		"Exit",
	}

	if config.AppData.Action == "" {
		ask.One(&survey.Select{
			Message: "What do you want to do?",
			Options: actions,
		}, &config.AppData.Action, nil)
	}

	if config.AppData.Action == "Exit" {
		os.Exit(0)
	}

	if config.AppData.Action == "Commit" {
		pkg_action.Commit(true)
	}

	if config.AppData.Action == "Push" {
		pkg_action.Push()
	}

	if config.AppData.Action == "Release" {
		helper.PrintName()
		pkg_action.Release()
	}

	if config.AppData.Action == "Publish" {
		pkg_action.Publish()
	}

	if config.AppData.Action == "Pull Request" {
		pkg_action.PullRequest()
	}

	if config.AppData.Action == "Create Branch" {
		pkg_action.Branch()
	}

	if config.AppData.Action == "Stash" {
		pkg_action.Stash()
	}

	if config.AppData.Action == "Merge" {
		pkg_action.Merge()
	}

	if config.AppData.Action == "Reset" {
		pkg_action.Reset()
	}
}
