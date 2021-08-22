package deploy

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/saurabhjambhule/yantra/pkg/git"
)

// deployCmd represents the deploy command
var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy",
	Run: func(cmd *cobra.Command, args []string) {
		runECSTask()
	},
}

func getImageTag() string {
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
