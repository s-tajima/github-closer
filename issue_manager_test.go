package main

import (
	"testing"
)

func TestBuildQuery(t *testing.T) {
	im := &IssueManager{}
	im.Organization = "sample-organization"
	im.Repository = "sample-repository"

	q := im.buildQuery("is:open")
	expect := "is:open repo:sample-organization/sample-repository"

	if q != expect {
		t.Error("buildQuery returns wrong result. (" + q + ")")
	}
}
