package scanner

import (
	"os"
	"path/filepath"
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
	repos := s.findRepos()
	since := s.calculateTimeRange()
	
	commitChan := make(chan git.Commit)
	errChan := make(chan error)
	var wg sync.WaitGroup
	
	// Scan repositories concurrently
	for _, repoPath := range repos {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			s.scanRepo(path, since, commitChan, errChan)
		}(repoPath)
	}
	
	// Close channels when done
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
	
	return commits, nil
}

func (s *Scanner) findRepos() []string {
	var repos []string
	pwd, err := os.Getwd()
	if err != nil {
		return repos
	}

	// Walk through all directories starting from current working directory
	filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil 
		}
		if info.IsDir() {
			if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
				repos = append(repos, path)
				return filepath.SkipDir 
			}
		}
		return nil
	})

	return repos
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
	
	// Get commits based on options
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
		errChan <- err
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