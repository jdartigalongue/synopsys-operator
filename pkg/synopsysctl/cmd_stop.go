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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// stopCmd stops a Synopsys resource in the cluster
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a Synopsys resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("must specify a sub-command")
	},
}

// stopAlertCmd stops an Alert instance
var stopAlertCmd = &cobra.Command{
	Use:           "alert NAME",
	Example:       "synopsysctl stop alert <name>\nsynopsysctl stop alert <name> -n <namespace>",
	Short:         "Stop an Alert instance",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.Help()
			return fmt.Errorf("this command takes 1 argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		alertName, alertNamespace, crdNamespace, _, err := getInstanceInfo(false, util.AlertCRDName, util.AlertName, namespace, args[0])
		if err != nil {
			return err
		}
		log.Infof("stopping Alert '%s' in namespace '%s'...", alertName, alertNamespace)

		// Get the Alert
		currAlert, err := util.GetAlert(alertClient, crdNamespace, alertName, v1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error stopping Alert '%s' in namespace '%s' due to %+v", alertName, alertNamespace, err)
		}

		// Make changes to Spec
		currAlert.Spec.DesiredState = "STOP"
		// Update Alert
		_, err = util.UpdateAlert(alertClient,
			crdNamespace, currAlert)
		if err != nil {
			return fmt.Errorf("error stopping Alert '%s' in namespace '%s' due to %+v", alertName, alertNamespace, err)
		}

		log.Infof("successfully submitted stop Alert '%s' in namespace '%s'", alertName, alertNamespace)
		return nil
	},
}

// stopBlackDuckCmd stops a Black Duck instance
var stopBlackDuckCmd = &cobra.Command{
	Use:           "blackduck NAME",
	Example:       "synopsysctl stop blackduck <name>\nsynopsysctl stop blackduck <name> -n <namespace>",
	Short:         "Stop a Black Duck instance",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.Help()
			return fmt.Errorf("this command takes 1 argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		blackDuckName, blackDuckNamespace, crdNamespace, _, err := getInstanceInfo(false, util.BlackDuckCRDName, util.BlackDuckName, namespace, args[0])
		if err != nil {
			return err
		}
		log.Infof("stopping Black Duck '%s' in namespace '%s'...", blackDuckName, blackDuckNamespace)

		// Get the Black Duck
		currBlackDuck, err := util.GetBlackduck(blackDuckClient, crdNamespace, blackDuckName, v1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error getting Black Duck '%s' in namespace '%s' due to %+v", blackDuckName, blackDuckNamespace, err)
		}

		// Make changes to Spec
		currBlackDuck.Spec.DesiredState = "STOP"
		// Update Black Duck
		_, err = util.UpdateBlackduck(blackDuckClient, currBlackDuck)
		if err != nil {
			return fmt.Errorf("error updating Black Duck '%s' in namespace '%s' due to %+v", blackDuckName, blackDuckNamespace, err)
		}

		log.Infof("successfully submitted stop Black Duck '%s' in namespace '%s'", blackDuckName, blackDuckNamespace)
		return nil
	},
}

// stopOpsSightCmd stops an OpsSight instance
var stopOpsSightCmd = &cobra.Command{
	Use:           "opssight NAME",
	Example:       "synopsysctl stop opssight <name>",
	Short:         "Stop an OpsSight instance",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.Help()
			return fmt.Errorf("this command takes 1 argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		opsSightName, opsSightNamespace, crdNamespace, _, err := getInstanceInfo(false, util.OpsSightCRDName, util.OpsSightName, namespace, args[0])
		if err != nil {
			return err
		}
		log.Infof("stopping OpsSight '%s' in namespace '%s'...", opsSightName, opsSightNamespace)

		// Get the OpsSight
		currOpsSight, err := util.GetOpsSight(opsSightClient, crdNamespace, opsSightName, v1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error getting OpsSight '%s' in namespace '%s' due to %+v", opsSightName, opsSightNamespace, err)
		}

		// Make changes to Spec
		currOpsSight.Spec.DesiredState = "STOP"
		// Update OpsSight
		_, err = util.UpdateOpsSight(opsSightClient,
			crdNamespace, currOpsSight)
		if err != nil {
			return fmt.Errorf("error updating OpsSight '%s' in namespace '%s' due to %+v", opsSightName, opsSightNamespace, err)
		}

		log.Infof("successfully submitted stop OpsSight '%s' in namespace '%s'", opsSightName, opsSightNamespace)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopAlertCmd.Flags().StringVarP(&namespace, "namespace", "n", namespace, "Namespace of the instance(s)")
	stopCmd.AddCommand(stopAlertCmd)

	stopBlackDuckCmd.Flags().StringVarP(&namespace, "namespace", "n", namespace, "Namespace of the instance(s)")
	stopCmd.AddCommand(stopBlackDuckCmd)

	stopCmd.AddCommand(stopOpsSightCmd)
}
