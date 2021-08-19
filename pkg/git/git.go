package git

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetCurrentRepo() *git.Repository {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	repository, err := git.PlainOpen(dir)
	if err != nil {
		log.Fatal(err)
	}

	return repository
}

func GetBranch() string {
	repository := GetCurrentRepo()

	branchRefs, err := repository.Branches()
	if err != nil {
		log.Fatal(err)
	}

	headRef, err := repository.Head()
	if err != nil {
		log.Fatal(err)
	}

	var currentBranch string
	err = branchRefs.ForEach(func(branchRef *plumbing.Reference) error {
		if branchRef.Hash() == headRef.Hash() {
			currentBranch = branchRef.Name().Short()

			return nil
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return currentBranch
}

func GetCommit() string {
	repository := GetCurrentRepo()

	headRef, err := repository.Head()
	if err != nil {
		log.Fatal(err)
	}
	headSha := headRef.Hash().String()

	return headSha
}
