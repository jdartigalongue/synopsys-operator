/*
Copyright (C) 2019 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package synopsysctl

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// startCmd starts a Synopsys resource in the cluster
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Synopsys resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("must specify a sub-command")
	},
}

// startAlertCmd starts an Alert instance
var startAlertCmd = &cobra.Command{
	Use:           "alert NAME",
	Example:       "synopsysctl start alert <name>\nsynopsysctl start alert <name1> <name2>\nsynopsysctl start alert <name> -n <namespace>\nsynopsysctl start alert <name1> <name2> -n <namespace>",
	Short:         "Start an Alert instance",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return fmt.Errorf("this command takes one or more arguments")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		errors := []error{}
		for _, alertName := range args {
			alertNamespace, crdNamespace, _, err := getInstanceInfo(false, util.AlertCRDName, util.AlertName, namespace, alertName)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			log.Infof("starting Alert '%s' in namespace '%s'...", alertName, alertNamespace)

			// Get the Alert
			currAlert, err := util.GetAlert(alertClient, crdNamespace, alertName, metav1.GetOptions{})
			if err != nil {
				errors = append(errors, fmt.Errorf("error getting Alert '%s' in namespace '%s' due to %+v", alertName, alertNamespace, err))
				continue
			}

			// Make changes to Spec
			currAlert.Spec.DesiredState = ""
			// Update Alert
			_, err = util.UpdateAlert(alertClient, crdNamespace, currAlert)
			if err != nil {
				errors = append(errors, fmt.Errorf("error starting Alert '%s' in namespace '%s' due to %+v", alertName, alertNamespace, err))
				continue
			}

			log.Infof("successfully submitted start Alert '%s' in namespace '%s'", alertName, alertNamespace)
		}
		if len(errors) > 0 {
			return fmt.Errorf("%v", errors)
		}
		return nil
	},
}

// startBlackDuckCmd starts a Black Duck instance
var startBlackDuckCmd = &cobra.Command{
	Use:           "blackduck NAME",
	Example:       "synopsysctl start blackduck <name>\nsynopsysctl start blackduck <name1> <name2>\nsynopsysctl start blackduck <name> -n <namespace>\nsynopsysctl start blackduck <name1> <name2> -n <namespace>",
	Short:         "Start a Black Duck instance",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return fmt.Errorf("this command takes one or more arguments")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		errors := []error{}
		for _, blackDuckName := range args {
			blackDuckNamespace, crdNamespace, _, err := getInstanceInfo(false, util.BlackDuckCRDName, util.BlackDuckName, namespace, blackDuckName)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			log.Infof("starting Black Duck '%s' in namespace '%s'...", blackDuckName, blackDuckNamespace)

			// Get the Black Duck
			currBlackDuck, err := util.GetBlackduck(blackDuckClient, crdNamespace, blackDuckName, metav1.GetOptions{})
			if err != nil {
				errors = append(errors, fmt.Errorf("error getting Black Duck '%s' in namespace '%s' due to %+v", blackDuckName, blackDuckNamespace, err))
				continue
			}

			// Make changes to Spec
			currBlackDuck.Spec.DesiredState = ""
			// Update Blackduck
			_, err = util.UpdateBlackduck(blackDuckClient, currBlackDuck)
			if err != nil {
				errors = append(errors, fmt.Errorf("error starting Black Duck '%s' in namespace '%s' due to %+v", blackDuckName, blackDuckNamespace, err))
				continue
			}

			log.Infof("successfully submitted start Black Duck '%s' in namespace '%s'", blackDuckName, blackDuckNamespace)
		}
		if len(errors) > 0 {
			return fmt.Errorf("%v", errors)
		}
		return nil
	},
}

// startOpsSightCmd starts an OpsSight instance
var startOpsSightCmd = &cobra.Command{
	Use:           "opssight NAME",
	Example:       "synopsysctl start opssight <name>\nsynopsysctl start opssight <name1> <name2>",
	Short:         "Start an OpsSight instance",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return fmt.Errorf("this command takes one or more arguments")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		errors := []error{}
		for _, opsSightName := range args {
			opsSightNamespace, crdNamespace, _, err := getInstanceInfo(false, util.OpsSightCRDName, util.OpsSightName, namespace, opsSightName)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			log.Infof("starting OpsSight '%s' in namespace '%s'...", opsSightName, opsSightNamespace)

			// Get the OpsSight
			currOpsSight, err := util.GetOpsSight(opsSightClient, crdNamespace, opsSightName, metav1.GetOptions{})
			if err != nil {
				errors = append(errors, fmt.Errorf("error getting OpsSight '%s' in namespace '%s' due to %+v", opsSightName, opsSightNamespace, err))
				continue
			}

			// Make changes to Spec
			currOpsSight.Spec.DesiredState = ""
			// Update OpsSight
			_, err = util.UpdateOpsSight(opsSightClient, crdNamespace, currOpsSight)
			if err != nil {
				errors = append(errors, fmt.Errorf("error starting OpsSight '%s' in namespace '%s' due to %+v", opsSightName, opsSightNamespace, err))
				continue
			}

			log.Infof("successfully submitted start OpsSight '%s' in namespace '%s'", opsSightName, opsSightNamespace)
		}
		if len(errors) > 0 {
			return fmt.Errorf("%v", errors)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startAlertCmd.Flags().StringVarP(&namespace, "namespace", "n", namespace, "Namespace of the instance(s)")
	startCmd.AddCommand(startAlertCmd)

	startBlackDuckCmd.Flags().StringVarP(&namespace, "namespace", "n", namespace, "Namespace of the instance(s)")
	startCmd.AddCommand(startBlackDuckCmd)

	startOpsSightCmd.Flags().StringVarP(&namespace, "namespace", "n", namespace, "Namespace of the instance(s)")
	startCmd.AddCommand(startOpsSightCmd)
}
