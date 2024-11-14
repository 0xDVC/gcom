package git

import (
    "time"
)

type Commit struct {
    Hash      string
    Author    string
    Date      time.Time
    Message   string
    RepoPath  string
    Pushed    bool
}

func (c *Commit) ShortHash() string {
    if len(c.Hash) < 8 {
        return c.Hash
    }
    return c.Hash[:8]
} 