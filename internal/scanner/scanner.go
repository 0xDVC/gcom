package scanner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/0xdvc/gcom/internal/git"
)

type Options struct {
	Today    bool
	Hours    int
	Days     int
	Unpushed bool
	Pushed   bool
}

type Scanner struct {
	opts Options
}

func New(opts Options) *Scanner {
	return &Scanner{opts: opts}
}

func (s *Scanner) Scan() ([]git.Commit, error) {
	commitChan := make(chan git.Commit, 1000)
	errChan := make(chan error, 100)
	var wg sync.WaitGroup
	
	repos, err := s.findRepos()
	if err != nil {
		return nil, err
	}
	
	fmt.Printf("Found %d repositories\n", len(repos))
	since := s.calculateTimeRange()
	
	for _, repo := range repos {
		wg.Add(1)
		go func(repoPath string) {
			defer wg.Done()
			s.scanRepo(repoPath, since, commitChan, errChan)
		}(repo)
	}
	
	go func() {
		wg.Wait()
		close(commitChan)
		close(errChan)
	}()
	
	var commits []git.Commit
	for commit := range commitChan {
		if s.shouldIncludeCommit(commit) {
			commits = append(commits, commit)
		}
	}
	
	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	default:
	}
	
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})
	
	return commits, nil
}

func (s *Scanner) findRepos() ([]string, error) {
	var repos []string
	pwd, err := os.Getwd()
	if err != nil {
		return repos, fmt.Errorf("failed to get working directory: %v", err)
	}

	fmt.Printf("Scanning for git repositories in: %s\n", pwd)

	if isGitRepo(pwd) {
		repos = append(repos, pwd)
	}

	entries, err := os.ReadDir(pwd)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(pwd, entry.Name())
			if isGitRepo(path) {
				repos = append(repos, path)
			}
		}
	}

	if len(repos) == 0 {
		return nil, fmt.Errorf("no git repositories found in %s", pwd)
	}

	return repos, nil
}

func isGitRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		cmd := exec.Command("git", "-C", path, "rev-parse", "--git-dir")
		if err := cmd.Run(); err == nil {
			fmt.Printf("Found git repository: %s\n", path)
			return true
		}
	}
	return false
}

func (s *Scanner) calculateTimeRange() time.Time {
	now := time.Now()
	
	if s.opts.Today {
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}
	
	if s.opts.Hours > 0 {
		return now.Add(-time.Duration(s.opts.Hours) * time.Hour)
	}
	
	if s.opts.Days > 0 {
		return now.AddDate(0, 0, -s.opts.Days)
	}
	
	return now.Add(-24 * time.Hour)
}

func (s *Scanner) scanRepo(path string, since time.Time, commitChan chan<- git.Commit, errChan chan<- error) {
	repo := &git.Repo{Path: path}
	
	var commits []git.Commit
	var err error
	
	if s.opts.Unpushed {
		commits, err = repo.GetUnpushedCommits(since)
	} else if s.opts.Pushed {
		commits, err = repo.GetPushedCommits(since)
	} else {
		commits, err = repo.GetAllCommits(since)
	}
	
	if err != nil {
		errChan <- fmt.Errorf("error scanning repo %s: %w", path, err)
		return
	}
	
	for _, commit := range commits {
		commit.RepoPath = path
		commitChan <- commit
	}
}

func (s *Scanner) shouldIncludeCommit(commit git.Commit) bool {
	if !s.opts.Pushed && !s.opts.Unpushed {
		return true
	}
	
	// Include only pushed commits
	if s.opts.Pushed && commit.Pushed {
		return true
	}
	
	if s.opts.Unpushed && !commit.Pushed {
		return true
	}
	
	return false
} 