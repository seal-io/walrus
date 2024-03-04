// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"testing"
)

func TestClient_Base(t *testing.T) {
	client, err := New("https://dev.azure.com", "org", "proj")
	if err != nil {
		t.Error(err)
	}
	got, want := client.BaseURL.String(), "https://dev.azure.com/"
	if got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Default(t *testing.T) {
	client := NewDefault("org", "proj")
	if got, want := client.BaseURL.String(), "https://dev.azure.com/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_azure_special(t *testing.T) {
	client, _ := New("https://dev.azure.com", "org", "")
	if got, want := client.BaseURL.String(), "https://dev.azure.com/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
	client2, _ := New("https://dev.azure.com", "", "proj")
	if client2 != nil {
		t.Errorf("Want nil client, got %v", client2)
	}
}

func TestSanitizeBranchName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"master",
			args{
				"master",
			},
			"refs/heads/master",
		},
		{
			"refs/heads/master",
			args{
				"refs/heads/master",
			},
			"refs/heads/master",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeBranchName(tt.args.name); got != tt.want {
				t.Errorf("SanitizeBranchName() = %v, want %v", got, tt.want)
			}
		})
	}
}
