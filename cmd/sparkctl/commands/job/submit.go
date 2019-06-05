package job

import (
	"fmt"
	"os"
	"time"

	jobcontroller "spark-cluster/pkg/controller/job"

	"github.com/spf13/cobra"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const (
	Username        = "wdongyu"
	Email           = "wdongyu@outlook.com"
	waitingInterval = 10 * time.Second
)

func NewSubmitCommand() *cobra.Command {

	var (
		buildID int
	)

	var command = &cobra.Command{
		Use:   "submit",
		Short: "submit a spark job.",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			if err := submitJob(buildID); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	// command.Flags().IntVar(&buildID, "buildID", 1, "the ID(number) of the build.")

	return command
}

func submitJob(buildID int) error {
	if err := GitPush(); err != nil {
		fmt.Println(err)
		return err
	}

	controller, err := jobcontroller.New()
	if err != nil {
		return err
	}

	err = controller.Execute()
	if err != nil {
		return err
	}

	return nil
}

func GitPush() error {
	fmt.Printf("[sparkctl] Setting up repo and push origin\n")
	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// status, err := worktree.Status()
	// fmt.Println(status, err)

	commit, err := worktree.Commit("add .drone.yml", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  Username,
			Email: Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	obj, err := repo.CommitObject(commit)
	if err != nil {
		return err
	}
	fmt.Println(obj)

	if err = repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: Username,
			Password: "gogs",
		},
		Progress: os.Stdout,
	}); err != nil {
		return err
	}

	return nil
}
