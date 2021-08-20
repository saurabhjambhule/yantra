package deploy

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/saurabhjambhule/yantra/pkg/aws"
	"github.com/saurabhjambhule/yantra/pkg/git"
	"github.com/saurabhjambhule/yantra/internal/utils"
)

// deployCmd represents the deploy command
var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy",
	Run: func(cmd *cobra.Command, args []string) {
		runECSTask()
	},
}

func ecsDeploy()  {
	imageTag := getImageTag()

	awsProfile := "default"
	awsRegion := "us-east-1"
	session := aws.StartSession(awsProfile, awsRegion)
	createdAt := checkImageFromECR(session, imageTag, "dash-test")

	fmt.Println("The image: " + imageTag + " created " + createdAt + " before!")
	utils.UserConfirmation()
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

	return imageTag
}
