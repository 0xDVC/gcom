package git

import (
	"os/exec"
	"strings"
	"time"
	"fmt"
	"strconv"
)

type Repo struct {
	Path string
}

func (r *Repo) GetCommits(revision string, since time.Time) ([]Commit, error) {
	format := "--pretty=format:%H%x00%an%x00%at%x00%s"
	args := []string{"-C", r.Path, "log", format}
	
	if revision != "" {
		args = append(args, revision)
	}
	
	if !since.IsZero() {
		args = append(args, "--since", since.Format("2006-01-02 15:04:05"))
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return r.parseCommits(string(output))
}

func (r *Repo) IsCommitPushed(hash string) bool {
	cmd := exec.Command("git", "-C", r.Path, "branch", "-r", "--contains", hash)
	output, err := cmd.Output()
	return err == nil && len(output) > 0
}

func (r *Repo) parseCommits(output string) ([]Commit, error) {
	var commits []Commit
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		commit, err := r.parseCommitLine(line)
		if err != nil {
			continue
		}
		
		commits = append(commits, commit)
	}
	
	return commits, nil
}

func (r *Repo) parseCommitLine(line string) (Commit, error) {
	parts := strings.Split(line, "\x00")
	
	if len(parts) != 4 {
		return Commit{}, fmt.Errorf("invalid commit line format: %s", line)
	}
	
	timestamp, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return Commit{}, fmt.Errorf("failed to parse timestamp: %v", err)
	}
	
	// Create and return the commit
	return Commit{
		Hash:      parts[0],
		Author:    parts[1],
		Date:      time.Unix(timestamp, 0),
		Message:   parts[3],
	}, nil
}

func (r *Repo) GetUnpushedCommits(since time.Time) ([]Commit, error) {
	return r.GetCommits("HEAD --not --remotes", since)
}

func (r *Repo) GetPushedCommits(since time.Time) ([]Commit, error) {
	return r.GetCommits("--remotes", since)
}

func (r *Repo) GetAllCommits(since time.Time) ([]Commit, error) {
	args := []string{
		"-C", r.Path,
		"log",
		"--since", since.Format(time.RFC3339),
		"--date=iso-strict",
		"--format=%H%x00%an%x00%at%x00%s",
		"--all",
	}
	
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git log failed in %s: %v\nOutput: %s", 
			r.Path, err, string(output))
	}
	
	commits, err := r.parseCommits(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to parse commits from %s: %v", r.Path, err)
	}
	
	return commits, nil
}

