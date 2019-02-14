// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("edit called")
	},
}

var editBlackduckCmd = &cobra.Command{
	Use:   "blackduck",
	Short: "Edit an instance of Blackduck",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Editing Blackduck")
	},
}

var editOpsSightCmd = &cobra.Command{
	Use:   "opssight",
	Short: "Edit an instance of OpsSight",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Editing OpsSight")
	},
}

var editAlertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Edit an instance of Alert",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Editing Alert")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.AddCommand(editBlackduckCmd)
	editCmd.AddCommand(editOpsSightCmd)
	editCmd.AddCommand(editAlertCmd)
}
