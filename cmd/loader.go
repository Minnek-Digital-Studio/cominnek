// This file is to add all commands generated by plugins

package cmd

import (
	"fmt"

	"github.com/Minnek-Digital-Studio/cominnek/controllers/app"
	"github.com/spf13/cobra"
)

var flags []app.IFlag
var subFlags [][]app.IFlag
var indexer []struct {
	Index int
	Name  string
}

func mainLoader(plugin app.IPlugin) cobra.Command {
	cmd := cobra.Command{
		Use:     plugin.Name,
		Short:   plugin.Help,
		Version: plugin.Version,
		Run: func(cmd *cobra.Command, args []string) {
			println("Hello from plugin", plugin.Name)

			for _, flag := range flags {
				fmt.Println(GetData(flag.Long, *cmd))
			}
		},
	}

	cmd.SetVersionTemplate(plugin.Name + " {{.Version}}\n")
	return cmd
}

func subLoader(subCommad app.IScripts) cobra.Command {
	cmd := cobra.Command{
		Use:   subCommad.Name,
		Short: subCommad.Help,
		Run: func(cmd *cobra.Command, args []string) {
			println("Hello from plugin subcommand", subCommad.Name)

			for _, index := range indexer {
				if index.Name != subCommad.Name {
					continue
				}

				flg := subFlags[index.Index]

				for _, flag := range flg {
					fmt.Println(GetData(flag.Long, *cmd))
				}
			}
		},
	}

	setSubCommand(&cmd, subCommad.Scripts)

	return cmd
}

func setSubCommand(cmd *cobra.Command, sub []app.IScripts) {
	if len(sub) > 0 {
		for _, sub := range sub {
			subLoad := subLoader(sub)
			pos := len(subFlags)
			subFlags = append(subFlags, sub.Flags)
			indexer = append(indexer, struct {
				Index int
				Name  string
			}{
				Index: pos,
				Name:  sub.Name,
			})

			for _, flag := range sub.Flags {
				FlagLoader(flag, &subLoad)
			}

			cmd.AddCommand(&subLoad)
		}
	}
}

func FlagLoader(flag app.IFlag, cmd *cobra.Command) {
	if flag.Type == "bool" {
		defaultF := flag.Default.(bool)
		cmd.PersistentFlags().BoolP(flag.Long, flag.Short, defaultF, flag.Help)
	}

	if flag.Type == "string" {
		defaultF := flag.Default.(string)
		cmd.PersistentFlags().StringP(flag.Long, flag.Short, defaultF, flag.Help)
	}
}

func GetData(name string, cmd cobra.Command) interface{} {
	return cmd.PersistentFlags().Lookup(name).Value
}

func init() {
	app.ConfigReader()
	plugins := app.ConfigGlobal.Plugins

	for _, plugin := range plugins {
		mainLoad := mainLoader(plugin)

		setSubCommand(&mainLoad, plugin.Scripts)

		flags = plugin.Flags

		for _, flag := range plugin.Flags {
			FlagLoader(flag, &mainLoad)
		}

		rootCmd.AddCommand(&mainLoad)
	}

}