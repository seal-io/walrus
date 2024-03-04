package gitlab

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

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/milestones/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/milestone.json")

	client := NewDefault()
	got, res, err := client.Milestones.Find(context.Background(), "diaspora/diaspora", 1)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Milestone)
	raw, err := ioutil.ReadFile("testdata/milestone.json.golden")
	if err != nil {
		t.Fatalf("ioutil.ReadFile: %v", err)
	}
	if err := json.Unmarshal(raw, want); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
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

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/milestones").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/milestones.json")

	client := NewDefault()
	got, res, err := client.Milestones.List(context.Background(), "diaspora/diaspora", scm.MilestoneListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Milestone{}
	raw, err := ioutil.ReadFile("testdata/milestones.json.golden")
	if err != nil {
		t.Fatalf("ioutil.ReadFile: %v", err)
	}
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
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

	gock.New("https://gitlab.com").
		Post("/api/v4/projects/diaspora/diaspora/milestones").
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
	got, res, err := client.Milestones.Create(context.Background(), "diaspora/diaspora", input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Milestone)
	raw, err := ioutil.ReadFile("testdata/milestone.json.golden")
	if err != nil {
		t.Fatalf("ioutil.ReadFile: %v", err)
	}
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
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

	gock.New("https://gitlab.com").
		Patch("/api/v4/projects/diaspora/diaspora/milestones/1").
		File("testdata/milestone_update.json").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/milestone.json")

	client := NewDefault()
	dueDate, _ := time.Parse(scm.SearchTimeFormat, "2012-10-09T23:39:01Z")
	input := &scm.MilestoneInput{
		Title:       "v1.0",
		Description: "Tracking milestone for version 1.0",
		State:       "close",
		DueDate:     dueDate,
	}
	got, res, err := client.Milestones.Update(context.Background(), "diaspora/diaspora", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Milestone)
	raw, err := ioutil.ReadFile("testdata/milestone.json.golden")
	if err != nil {
		t.Fatalf("ioutil.ReadFile: %v", err)
	}
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
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

	gock.New("https://gitlab.com").
		Delete("/api/v4/projects/diaspora/diaspora/milestones/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	_, err := client.Milestones.Delete(context.Background(), "diaspora/diaspora", 1)
	if err != nil {
		t.Error(err)
		return
	}
}
