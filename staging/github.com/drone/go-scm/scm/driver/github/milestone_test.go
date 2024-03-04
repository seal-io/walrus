package github

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestMilestoneFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/milestones/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/milestone.json")

	client := NewDefault()
	got, res, err := client.Milestones.Find(context.Background(), "octocat/hello-world", 1)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Milestone)
	raw, _ := ioutil.ReadFile("testdata/milestone.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestMilestoneList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/milestones").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		MatchParam("state", "all").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/milestones.json")
	client := NewDefault()
	got, res, err := client.Milestones.List(context.Background(), "octocat/hello-world", scm.MilestoneListOptions{Page: 1, Size: 30, Open: true, Closed: true})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Milestone{}
	raw, _ := ioutil.ReadFile("testdata/milestones.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestMilestoneCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/milestones").
		File("testdata/milestone_create.json").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/milestone.json")

	client := NewDefault()
	dueDate, _ := time.Parse(scm.SearchTimeFormat, "2012-10-09T23:39:01Z")
	input := &scm.MilestoneInput{
		Title:       "v1.0",
		Description: "Tracking milestone for version 1.0",
		State:       "open",
		DueDate:     dueDate,
	}

	got, res, err := client.Milestones.Create(context.Background(), "octocat/hello-world", input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Milestone)
	raw, _ := ioutil.ReadFile("testdata/milestone.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestMilestoneUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Patch("/repos/octocat/hello-world/milestones/1").
		File("testdata/milestone_create.json").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/milestone.json")

	client := NewDefault()
	dueDate, _ := time.Parse(scm.SearchTimeFormat, "2012-10-09T23:39:01Z")
	input := &scm.MilestoneInput{
		Title:       "v1.0",
		Description: "Tracking milestone for version 1.0",
		State:       "open",
		DueDate:     dueDate,
	}

	got, res, err := client.Milestones.Update(context.Background(), "octocat/hello-world", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Milestone)
	raw, _ := ioutil.ReadFile("testdata/milestone.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestMilestoneDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Delete("/repos/octocat/hello-world/milestones/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.Milestones.Delete(context.Background(), "octocat/hello-world", 1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
