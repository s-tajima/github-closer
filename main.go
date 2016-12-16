package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Options struct {
	Query    string `short:"q" long:"query" default:"type:issue is:open -label:Pending" description:"Query strings. For search Issues/Pull Requests."`
	Duration int    `short:"d" long:"duration" default:"30" description:"Duration. Issues would be closed if left over this duration. (days)"`
	Comment  string `short:"c" long:"comment" default:":alarm_clock: this Issue was left for a long time." description:"Comment. Would be posted before an Issue is closed."`
	DryRun   bool   `short:"n" long:"dry-run" description:"If true, show target Issues without closing."`
	RunOnce  bool   `short:"o" long:"run-once" description:"If true, close only one Issue."`
	Limit    int    `short:"l" long:"limit" default:"0" description:"A maximum number of closed Issues."`
	Debug    bool   `long:"debug"`
}

func main() {
	godotenv.Load()
	org := os.Getenv("GITHUB_ORGANIZATION")
	repo := os.Getenv("GITHUB_REPO")
	team := os.Getenv("GITHUB_TEAM")
	token := os.Getenv("GITHUB_ACCESS_TOKEN")

	var opts Options
	p := flags.NewParser(&opts, flags.Default)
	_, err := p.ParseArgs(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
	im := NewIssueManager(org, repo, team, token, opts.Duration)

	dryRunMessage := ""
	if opts.DryRun {
		dryRunMessage = " (dryrun)"
	}

	issues, err := im.FindIssues(opts.Query)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	closedNumber := 0
	for _, issue := range issues {
		if im.isUpdatedWithinDuration(issue) {
			if opts.Debug {
				log.Println(fmt.Sprintf("#%d %s %s was updated recently. skipped.", *issue.Number, *issue.HTMLURL, *issue.Title))
			}
			continue
		}

		if !opts.DryRun {
			if len(opts.Comment) > 0 {
				im.Comment(issue, opts.Comment)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			_, err = im.Close(issue)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		if opts.Debug {
			log.Println(fmt.Sprintf("#%d %s %s is closed.%s", *issue.Number, *issue.HTMLURL, *issue.Title, dryRunMessage))
		}

		closedNumber++
		if opts.Limit > 0 && opts.Limit <= closedNumber {
			break
		}

		if opts.RunOnce {
			break
		}
	}
}
