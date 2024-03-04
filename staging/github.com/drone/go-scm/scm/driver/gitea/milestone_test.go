package gitea

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

	mockServerVersion()

	gock.New("https://try.gitea.io").
		Get("/api/v1/repos/jcitizen/my-repo/milestones/1").
		Reply(200).
		Type("application/json").
		File("testdata/milestone.json")

	client, _ := New("https://try.gitea.io")
	got, _, err := client.Milestones.Find(context.Background(), "jcitizen/my-repo", 1)
	if err != nil {
		t.Error(err)
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
}

func TestMilestoneList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://try.gitea.io").
		Get("/api/v1/repos/jcitizen/my-repo/milestones").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/milestones.json")

	client, _ := New("https://try.gitea.io")
	got, res, err := client.Milestones.List(context.Background(), "jcitizen/my-repo", scm.MilestoneListOptions{})
	if err != nil {
		t.Error(err)
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

	t.Run("Page", testPage(res))
}

func TestMilestoneCreate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://try.gitea.io").
		Post("/api/v1/repos/jcitizen/my-repo/milestones").
		File("testdata/milestone_create.json").
		Reply(200).
		Type("application/json").
		File("testdata/milestone.json")

	client, _ := New("https://try.gitea.io")
	dueDate, _ := time.Parse(scm.SearchTimeFormat, "2012-10-09T23:39:01Z")
	input := &scm.MilestoneInput{
		Title:       "v1.0",
		Description: "Tracking milestone for version 1.0",
		State:       "open",
		DueDate:     dueDate,
	}
	got, _, err := client.Milestones.Create(context.Background(), "jcitizen/my-repo", input)
	if err != nil {
		t.Error(err)
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
}

func TestMilestoneUpdate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://try.gitea.io").
		Patch("/api/v1/repos/jcitizen/my-repo/milestones").
		File("testdata/milestone_create.json").
		Reply(200).
		Type("application/json").
		File("testdata/milestone.json")

	client, _ := New("https://try.gitea.io")
	dueDate, _ := time.Parse(scm.SearchTimeFormat, "2012-10-09T23:39:01Z")
	input := &scm.MilestoneInput{
		Title:       "v1.0",
		Description: "Tracking milestone for version 1.0",
		State:       "open",
		DueDate:     dueDate,
	}
	got, _, err := client.Milestones.Update(context.Background(), "jcitizen/my-repo", 1, input)
	if err != nil {
		t.Error(err)
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
}

func TestMilestoneDelete(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://try.gitea.io").
		Delete("/api/v1/repos/jcitizen/my-repo/milestones/1").
		Reply(200).
		Type("application/json")

	client, _ := New("https://try.gitea.io")
	_, err := client.Milestones.Delete(context.Background(), "jcitizen/my-repo", 1)
	if err != nil {
		t.Error(err)
	}
}

var mockPageHeaders = map[string]string{
	"Link": `<https://try.gitea.io/v1/resource?page=2>; rel="next",` +
		`<https://try.gitea.io/v1/resource?page=1>; rel="prev",` +
		`<https://try.gitea.io/v1/resource?page=1>; rel="first",` +
		`<https://try.gitea.io/v1/resource?page=5>; rel="last"`,
}

func mockServerVersion() {
	gock.New("https://try.gitea.io").
		Get("/api/v1/version").
		Reply(200).
		Type("application/json").
		File("testdata/version.json")
}
