package azure

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestContentFind(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/items").
		MatchParam("path", "README").
		MatchParam("versionDescriptor.version", "b1&b2").
		Reply(200).
		Type("application/json").
		File("testdata/content.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.Contents.Find(
		context.Background(),
		"REPOID",
		"README",
		"b1&b2",
	)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Content)
	raw, _ := ioutil.ReadFile("testdata/content.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestContentCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(201).
		Type("application/json").
		File("testdata/content_create.json")

	params := &scm.ContentParams{
		Message: "my commit message",
		Data:    []byte("bXkgbmV3IGZpbGUgY29udGVudHM="),
	}

	client := NewDefault("ORG", "PROJ")
	res, err := client.Contents.Create(
		context.Background(),
		"REPOID",
		"README",
		params,
	)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Status != 201 {
		t.Errorf("Unexpected Results")
	}
}

func TestContentUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/content_update.json")

	params := &scm.ContentParams{
		Message: "a new commit message",
		Data:    []byte("bXkgdXBkYXRlZCBmaWxlIGNvbnRlbnRz"),
		Sha:     "95b966ae1c166bd92f8ae7d1c313e738c731dfc3",
	}

	client := NewDefault("ORG", "PROJ")
	res, err := client.Contents.Update(
		context.Background(),
		"REPOID",
		"README",
		params,
	)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Status != 200 {
		t.Errorf("Unexpected Results")
	}
}

func TestContentDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/content_delete.json")

	params := &scm.ContentParams{
		Message: "a new commit message",
		BlobID:  "95b966ae1c166bd92f8ae7d1c313e738c731dfc3",
	}

	client := NewDefault("ORG", "PROJ")
	res, err := client.Contents.Delete(
		context.Background(),
		"REPOID",
		"README",
		params,
	)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Status != 200 {
		t.Errorf("Unexpected Results")
	}
}

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/content_list.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.Contents.List(
		context.Background(),
		"REPOID",
		"",
		"",
		scm.ListOptions{},
	)
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.ContentInfo{}
	raw, _ := ioutil.ReadFile("testdata/content_list.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func Test_generateURIFromRef(t *testing.T) {
	type args struct {
		ref string
	}
	tests := []struct {
		name    string
		args    args
		wantUri string
	}{
		{
			name:    "branch",
			args:    args{ref: "branch"},
			wantUri: "&versionDescriptor.versionType=branch&versionDescriptor.version=branch",
		},
		{
			name:    "commit",
			args:    args{ref: "6bbcbc818c804d35b88a12bbd2ed297e41c4d10d"},
			wantUri: "&versionDescriptor.versionType=commit&versionDescriptor.version=6bbcbc818c804d35b88a12bbd2ed297e41c4d10d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotUri := generateURIFromRef(tt.args.ref); gotUri != tt.wantUri {
				t.Errorf("generateURIFromRef() = %v, want %v", gotUri, tt.wantUri)
			}
		})
	}
}
