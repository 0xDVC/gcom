package display

import (
	"fmt"
	"path/filepath"

	"github.com/0xdvc/gcom/internal/git"
)

func PrintCommits(commits []git.Commit) {
	for _, c := range commits {
		printCommit(c)
	}
}

func printCommit(c git.Commit) {
	pushStatus := "✓"
	if !c.Pushed {
		pushStatus = "✗"
	}
	
	fmt.Printf("\033[1m%s\033[0m %s\n", c.ShortHash(), pushStatus)
	fmt.Printf("Author: %s\n", c.Author)
	fmt.Printf("Date:   %s\n", c.Date.Format("Mon Jan 2 15:04:05 2006 -0700"))
	fmt.Printf("Repo:   %s\n", filepath.Base(c.RepoPath))
	fmt.Printf("\n    %s\n\n", c.Message)
} 