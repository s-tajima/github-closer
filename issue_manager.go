package main

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"regexp"
	"strings"
	"time"
)

type IssueManager struct {
	Client       *github.Client
	Organization string
	Team         string
	Repository   string
	Duration     int
}

type Users []*github.User

func String(v string) *string { return &v }

func NewIssueManager(organization string, repository string, team string, token string, duration int) *IssueManager {

	tc := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

	im := &IssueManager{}
	im.Client = github.NewClient(tc)
	im.Organization = organization
	im.Repository = repository
	im.Team = team
	im.Duration = duration

	return im
}

func (im *IssueManager) FindIssues(spec string) ([]*github.Issue, error) {
	members, err := im.findUsersByTeamName(im.Team)
	if err != nil {
		return nil, err
	}
	queryString := im.buildQuery(spec)
	searchResult, _, err := im.Client.Search.Issues(queryString, &github.SearchOptions{})
	if err != nil {
		return nil, err
	}

	var targets []*github.Issue
Loop:
	for i, issue := range searchResult.Issues {
		for _, member := range members {
			if *issue.User.Login == *member.Login {
				targets = append(targets, &searchResult.Issues[i])
				continue Loop
			}
		}
	}

	return targets, nil
}

func (im *IssueManager) Close(issue *github.Issue) (*github.Issue, error) {
	issueRequest := &github.IssueRequest{}
	issueRequest.State = String("closed")

	i, _, err := im.Client.Issues.Edit(im.Organization, im.Repository, *issue.Number, issueRequest)

	return i, err
}

func (im *IssueManager) isUpdatedWithinDuration(issue *github.Issue) bool {
	dhInt := im.Duration * 24
	dur, _ := time.ParseDuration(fmt.Sprintf("%dh", dhInt))

	updatedAt := issue.UpdatedAt

	return time.Since(*updatedAt) < dur
}

func (im *IssueManager) Comment(issue *github.Issue, comment string) bool {
	ic := &github.IssueComment{Body: &comment}
	_, _, err := im.Client.Issues.CreateComment(im.Organization, im.Repository, *issue.Number, ic)

	return err != nil
}

func (im *IssueManager) findUsersByTeamName(name string) ([]*github.User, error) {
	teams, _, err := im.Client.Repositories.ListTeams(im.Organization, im.Repository, nil)
	if err != nil {
		return nil, err
	}

	for _, t := range teams {
		if *t.Name == name {
			users, _, err := im.Client.Organizations.ListTeamMembers(*t.ID, &github.OrganizationListTeamMembersOptions{})
			return users, err
		}
	}

	return nil, errors.New("team not found")
}

func (im *IssueManager) buildQuery(spec string) string {
	queries := regexp.MustCompile(" +").Split(spec, -1)
	queries = append(queries, "repo:"+im.Organization+"/"+im.Repository)

	return strings.Join(queries, " ")
}
