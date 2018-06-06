package main

import (
	"flag"
	"os"
	"time"

	"github.com/vickleford/impetus/impetus"
)

type repeatingFlag []string

func (i *repeatingFlag) String() string {
	return "bwahhhhhhh"
}

func (i *repeatingFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var ghReposToScan repeatingFlag
	var orgsToScan repeatingFlag
	scanDelay := flag.Int("scandelay", 1, "Scan the configured targets for idle PRs this many hours")
	idleHours := flag.Int("idlehours", 24, "The number of hours without activity to consider a PR idle")
	hipchatRoomID := flag.String("roomid", "", "Hipchat Room API ID")
	indefinite := flag.Bool("indefinite", false, "Run indefinitely with a configurable interval")
	flag.Var(&ghReposToScan, "repo", "Github project to scan in the format org/repo")
	flag.Var(&orgsToScan, "org", "Github organization to scan in the format org")

	flag.Parse()

	if *hipchatRoomID == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(ghReposToScan) == 0 && len(orgsToScan) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	scanner := impetus.NewGhePullScanner()
	scanner.IdleToleranceHours = time.Duration(*idleHours)
	controller := impetus.NewController(*hipchatRoomID, ghReposToScan, orgsToScan, scanner)
	for {
		controller.ScanAndReport()
		if !*indefinite {
			break
		}
		time.Sleep(time.Duration(*scanDelay) * time.Hour)
	}
}
