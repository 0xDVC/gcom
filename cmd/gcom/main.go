package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/0xdvc/gcom/internal/scanner"
    "github.com/0xdvc/gcom/internal/display"
)

func main() {
    // CLI flags with shorthand and long versions
    var opts scanner.Options
    
    flag.BoolVar(&opts.Today, "t", false, "Show today's commits")
    flag.BoolVar(&opts.Today, "today", false, "Show today's commits")
    
    flag.BoolVar(&opts.Unpushed, "u", false, "Show only unpushed commits")
    flag.BoolVar(&opts.Unpushed, "unpushed", false, "Show only unpushed commits")
    
    flag.BoolVar(&opts.Pushed, "p", false, "Show only pushed commits")
    flag.BoolVar(&opts.Pushed, "pushed", false, "Show only pushed commits")
    
    flag.IntVar(&opts.Hours, "h", 0, "Show commits from last n hours")
    flag.IntVar(&opts.Hours, "hours", 0, "Show commits from last n hours")
    
    flag.IntVar(&opts.Days, "d", 0, "Show commits from last n days")
    flag.IntVar(&opts.Days, "days", 0, "Show commits from last n days")
    
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of gcom:\n")
        fmt.Fprintf(os.Stderr, "\ngcom [flags...]\n\n")
        fmt.Fprintf(os.Stderr, "Flags can be combined. Examples:\n")
        fmt.Fprintf(os.Stderr, "  gcom -t -u         Show today's unpushed commits\n")
        fmt.Fprintf(os.Stderr, "  gcom -d 7 -p       Show pushed commits from last 7 days\n")
        fmt.Fprintf(os.Stderr, "  gcom -h 12 -u -p   Show all commits from last 12 hours\n\n")
        fmt.Fprintf(os.Stderr, "Available flags:\n")
        flag.PrintDefaults()
    }
    
    flag.Parse()

    s := scanner.New(opts)

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