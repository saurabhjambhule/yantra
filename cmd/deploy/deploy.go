package deploy

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/saurabhjambhule/yantra/pkg/git"
	"github.com/saurabhjambhule/yantra/pkg/aws"
	"github.com/saurabhjambhule/yantra/pkg/aws/ecr"
	"github.com/saurabhjambhule/yantra/internal/utils"
)

// deployCmd represents the deploy command
var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy",
	Run: func(cmd *cobra.Command, args []string) {
		imageTag := getImageTag()
		_ = imageTag
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getImageTag() string {
	// gitBranch := "tapish/TRKINTEG-553_m365_unactivated_count_fix"
	gitBranch := git.GetBranch()
	imageTag := gitBranch

	if gitBranch == git.GetDefaultBranch() {
		imageTag = "latest"
	} else {
		imageTag = strings.Replace(imageTag, "-", "_", -1)
		imageTag = strings.Replace(imageTag, "/", "_", -1)
		imageTag = "dev-" + imageTag
	}

	createdAt := checkImageFromECR("default", "us-east-1", imageTag, "dash-test")

	fmt.Println("The image: " + imageTag + " created " + createdAt + " before!")
	utils.UserConfirmation()

	return imageTag
}

func checkImageFromECR(awsProfile string, awsRegion string, imageTag string, repoName string) string {
	session := aws.StartSession(awsProfile, awsRegion)
	_, createdAt := ecr.DoesImageExist(session, imageTag, repoName)

	return createdAt
}
