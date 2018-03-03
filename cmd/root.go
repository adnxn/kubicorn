// Copyright © 2017 The Kubicorn Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/kris-nova/kubicorn/pkg/local"
	"github.com/kris-nova/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	bashCompletionFunc = `
__kubicorn_parse_list()
{
    local kubicorn_out
    if kubicorn_out=$(kubicorn list --no-headers 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${kubicorn_out[*]}" -- "$cur" ) )
    fi
}
__kubicorn_parse_profiles()
{
    local kubicorn_out
    if kubicorn_out=(amazon aws digitalocean do); then
        COMPREPLY=( $( compgen -W "${kubicorn_out[*]}" -- "$cur" ) )
    fi
}
__custom_func() {
    case ${last_command} in
        kubicorn_apply | kubicorn_create | kubicorn_delete | kubicorn_getconfig)
            __kubicorn_parse_list
            return
            ;;
        *)
            ;;
    esac
}
`
)

var cfgFile = local.Expand("kubicorn.cfg")

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kubicorn",
	Short: "Kubernetes cluster management, without any magic",
	Long: fmt.Sprintf(`
%s
`, Unicorn),
	Run: func(cmd *cobra.Command, args []string) {
		if logger.Fabulous {
			cmd.SetOutput(logger.FabulousWriter)
		}
		if os.Getenv("KUBICORN_TRUECOLOR") != "" {
			cmd.SetOutput(logger.FabulousTrueWriter)
		}
		cmd.Help()
	},
	BashCompletionFunction: bashCompletionFunc,
}

func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//flags here
	addPersistentFlagInt(RootCmd, &logger.Level, "verbose", "v", 3, "Log level")
	addPersistentFlagBool(RootCmd, &logger.Color, "color", "C", true, "Toggle colorized logs")
	addPersistentFlagBool(RootCmd, &logger.Fabulous, "fab", "f", false, "Toggle colorized logs")

	// add commands
	addCommands()
}

func addCommands() {
	RootCmd.AddCommand(AdoptCmd())
	RootCmd.AddCommand(ApplyCmd())
	RootCmd.AddCommand(CompletionCmd())
	RootCmd.AddCommand(CreateCmd())
	RootCmd.AddCommand(DeleteCmd())
	RootCmd.AddCommand(EditCmd())
	RootCmd.AddCommand(GetConfigCmd())
	RootCmd.AddCommand(ImageCmd())
	RootCmd.AddCommand(ListCmd())
	RootCmd.AddCommand(VersionCmd())

	// Add Prompt at the end to initialize all the other commands first.
	RootCmd.AddCommand(PromptCmd())
}

func flagApplyAnnotations(cmd *cobra.Command, flag, completion string) {
	if cmd.Flag(flag) != nil {
		if cmd.Flag(flag).Annotations == nil {
			cmd.Flag(flag).Annotations = map[string][]string{}
		}
		cmd.Flag(flag).Annotations[cobra.BashCompCustom] = append(
			cmd.Flag(flag).Annotations[cobra.BashCompCustom],
			completion,
		)
	}
}

func initConfig() {

	viper.SetConfigType("yaml")
	viper.SetConfigFile(cfgFile)

	viper.SetEnvPrefix("KUBICORN")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("WIP: unable to read config file"))
		os.Exit(1)
	}
}

func addPersistentFlagInt(cmd *cobra.Command, field *int, name, shorthand string, value int, usage string) {
	cmd.PersistentFlags().IntVarP(field, name, shorthand, value, usage)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func addPersistentFlagBool(cmd *cobra.Command, field *bool, name, shorthand string, value bool, usage string) {
	cmd.PersistentFlags().BoolVarP(field, name, shorthand, value, usage)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}
