package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/0xdvc/gcom/internal/scanner"
    "github.com/0xdvc/gcom/internal/display"
)

func main() {
    // CLI flags
    unpushed := flag.Bool("u", false, "Show only unpushed commits")
    pushed := flag.Bool("p", false, "Show only pushed commits")
    today := flag.Bool("t", false, "Show today's commits")
    hours := flag.Int("h", 0, "Show commits from last n hours")
    days := flag.Int("d", 0, "Show commits from last n days")
    
    flag.Parse()

    // Initialize scanner
    s := scanner.New(scanner.Options{
        Today:    *today,
        Hours:    *hours,
        Days:     *days,
        Unpushed: *unpushed,
        Pushed:   *pushed,
    })

    // Scan for commits
    commits, err := s.Scan()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error scanning commits: %v\n", err)
        os.Exit(1)
    }

    if err := display.PrintCommits(commits); err != nil {
        fmt.Fprintf(os.Stderr, "Error displaying commits: %v\n", err)
        os.Exit(1)
    }
} 