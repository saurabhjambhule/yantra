package git

import (
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const remoteRef string = "refs/remotes/origin/"

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

func GetDefaultBranch() string {
	repository := GetCurrentRepo()

	refs, err := repository.References()
	if err != nil {
		log.Fatal(err)
	}

	var dafaultBranch string
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		// The HEAD is omitted in a `git show-ref` so we ignore the symbolic
		// references, the HEAD
		if ref.Type() == plumbing.SymbolicReference {
			return nil
		}

		if strings.HasPrefix(ref.Name().String(), remoteRef) {
			dafaultBranch = ref.Name().String()
		}

		return nil
	})

	dafaultBranch = strings.TrimPrefix(dafaultBranch, remoteRef)

	return dafaultBranch
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
