package display

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/0xdvc/gcom/internal/git"
)

// ANSI color codes
const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	Dim        = "\033[2m"
	
	Black      = "\033[30m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Magenta    = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	
	BgBlack    = "\033[40m"
	BgRed      = "\033[41m"
	BgGreen    = "\033[42m"
	BgYellow   = "\033[43m"
	BgBlue     = "\033[44m"
	BgMagenta  = "\033[45m"
	BgCyan     = "\033[46m"
	BgWhite    = "\033[47m"
)

func PrintCommits(commits []git.Commit) error {
	var buf bytes.Buffer
	
	header := fmt.Sprintf("🔍 Found %d commits", len(commits))
	fmt.Fprintf(&buf, "\n%s%s%s\n", Bold, header, Reset)
	fmt.Fprintf(&buf, "%s%s%s\n\n", Dim, strings.Repeat("─", len(header)), Reset)
	
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	
	for _, c := range commits {
		commitDate := c.Date.Truncate(24 * time.Hour)
		switch {
		case commitDate.Equal(today):
			fmt.Fprintf(&buf, "%s%s📅 Today%s\n%s%s%s\n", 
				Bold, Green, Reset, 
				Dim, strings.Repeat("─", 40), Reset)
		case commitDate.Equal(yesterday):
			fmt.Fprintf(&buf, "%s%s📅 Yesterday%s\n%s%s%s\n", 
				Bold, Yellow, Reset, 
				Dim, strings.Repeat("─", 40), Reset)
		default:
			fmt.Fprintf(&buf, "%s%s📅 %s%s\n%s%s%s\n", 
				Bold, Blue, 
				c.Date.Format("Monday, January 2"), Reset,
				Dim, strings.Repeat("─", 40), Reset)
		}
		
		formatCommit(&buf, c)
	}
	
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "less"
	}
	
	if strings.Contains(pager, "less") {
		pager = "less -R -F -X"
	}
	
	cmd := exec.Command("sh", "-c", pager)
	cmd.Stdin = &buf
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func formatCommit(buf *bytes.Buffer, c git.Commit) {
	pushStatus := fmt.Sprintf("%s✓%s", Green, Reset)
	if !c.Pushed {
		pushStatus = fmt.Sprintf("%s✗%s", Red, Reset)
	}
	
	hash := fmt.Sprintf("%s%s %s %s", 
		BgBlue, Black, c.ShortHash(), Reset)
	
	repoName := filepath.Base(c.RepoPath)
	repo := fmt.Sprintf("%s%s%s", Magenta, repoName, Reset)
	
	author := fmt.Sprintf("%s%s%s", Cyan, c.Author, Reset)
	
	timeStr := fmt.Sprintf("%s%s%s", 
		Dim, 
		c.Date.Format("15:04:05"), 
		Reset)
	
	fmt.Fprintf(buf, "\n%s %s [%s] %s@%s\n",
		hash,
		pushStatus,
		repo,
		author,
		timeStr)
	
	messageLines := strings.Split(c.Message, "\n")
	for i, line := range messageLines {
		if i == 0 {
			fmt.Fprintf(buf, "    %s%s%s\n", Bold, line, Reset)
		} else if line != "" {
			fmt.Fprintf(buf, "    %s%s%s\n", Dim, line, Reset)
		}
	}
	
	fmt.Fprintf(buf, "%s%s%s\n", 
		Dim, 
		strings.Repeat("· ", 40), 
		Reset)
} 