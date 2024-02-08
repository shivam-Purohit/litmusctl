/*
Copyright © 2021 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package get

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/litmuschaos/litmusctl/pkg/apis/environment"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

var ChaosEnvironmentCmd = &cobra.Command{
	Use:   "chaos-environment",
	Short: "Get Chaos Environment within the project",
	Long:  `Display the Chaos Environment within the project with the targeted id `,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectID, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if projectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&projectID)

			if projectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		environmentID, err := cmd.Flags().GetString("environment-id")
		utils.PrintError(err)

		if environmentID == "" {
			utils.White_B.Print("\nEnter the Environment ID: ")
			fmt.Scanln(&environmentID)

			if environmentID == "" {
				utils.Red.Println("⛔ Environment ID can't be empty!!")
				os.Exit(1)
			}
		}

		environmentGet, err := environment.GetChaosEnvironment(projectID, environmentID, credentials)
		if err != nil {
			if strings.Contains(err.Error(), "permission_denied") {
				utils.Red.Println("❌ You don't have enough permissions to access this resource.")
				os.Exit(1)
			} else {
				utils.PrintError(err)
				os.Exit(1)
			}
		}
		environmentGetData := environmentGet.Data.EnvironmentDetails

		writer := tabwriter.NewWriter(os.Stdout, 30, 8, 0, '\t', tabwriter.AlignRight)
		writer.Flush()
		intUpdateTime, err := strconv.ParseInt(environmentGetData.UpdatedAt, 10, 64)
		if err != nil {
			utils.Red.Println("Error converting UpdatedAt to int64:", err)
		}
		updatedTime := time.Unix(intUpdateTime, 0).String()
		intCreatedTime, err := strconv.ParseInt(environmentGetData.CreatedAt, 10, 64)
		if err != nil {
			utils.Red.Println("Error converting CreatedAt to int64:", err)
		}
		createdTime := time.Unix(intCreatedTime, 0).String()
		writer.Flush()
		utils.White_B.Fprintln(writer, "CHAOS ENVIRONMENT DETAILS")
		utils.White.Fprintln(writer, "CHAOS ENVIRONMENT ID\t", environmentGetData.EnvironmentID)
		utils.White.Fprintln(writer, "CHAOS ENVIRONMENT NAME\t", environmentGetData.Name)
		utils.White.Fprintln(writer, "CHAOS ENVIRONMENT Type\t", environmentGetData.Type)
		utils.White.Fprintln(writer, "CREATED AT\t", createdTime)
		utils.White.Fprintln(writer, "CREATED BY\t", environmentGetData.CreatedBy.Username)
		utils.White.Fprintln(writer, "UPDATED AT\t", updatedTime)
		utils.White.Fprintln(writer, "UPDATED BY\t", environmentGetData.UpdatedBy.Username)
		utils.White.Fprintln(writer, "CHAOS INFRA IDs\t", strings.Join(environmentGetData.InfraIDs, ", "))
		utils.White.Fprintln(writer, "TAGS\t", strings.Join(environmentGetData.Tags, ", "))
		writer.Flush()

	},
}

func init() {
	GetCmd.AddCommand(ChaosEnvironmentCmd)
	ChaosEnvironmentCmd.Flags().String("project-id", "", "Set the project-id to get Chaos Environment from a particular project.")
	ChaosEnvironmentCmd.Flags().String("environment-id", "", "Set the environment-id to get Chaos Environment")
}
