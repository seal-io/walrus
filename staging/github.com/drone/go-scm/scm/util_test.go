// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import "testing"

func TestSplit(t *testing.T) {
	tests := []struct {
		value, owner, name string
	}{
		{"octocat/hello-world", "octocat", "hello-world"},
		{"octocat/hello/world", "octocat", "hello/world"},
		{"hello-world", "", "hello-world"},
		{value: ""}, // empty value returns nothing
	}
	for _, test := range tests {
		owner, name := Split(test.value)
		if got, want := owner, test.owner; got != want {
			t.Errorf("Got repository owner %s, want %s", got, want)
		}
		if got, want := name, test.name; got != want {
			t.Errorf("Got repository name %s, want %s", got, want)
		}
	}
}

func TestJoin(t *testing.T) {
	got, want := Join("octocat", "hello-world"), "octocat/hello-world"
	if got != want {
		t.Errorf("Got repository name %s, want %s", got, want)
	}
}

func TestTrimRef(t *testing.T) {
	tests := []struct {
		before, after string
	}{
		{
			before: "refs/tags/v1.0.0",
			after:  "v1.0.0",
		},
		{
			before: "refs/heads/master",
			after:  "master",
		},
		{
			before: "refs/heads/feature/x",
			after:  "feature/x",
		},
		{
			before: "master",
			after:  "master",
		},
	}
	for _, test := range tests {
		if got, want := TrimRef(test.before), test.after; got != want {
			t.Errorf("Got reference %s, want %s", got, want)
		}
	}
}

func TestExpandRef(t *testing.T) {
	tests := []struct {
		name, prefix, after string
	}{
		// tag references
		{
			after:  "refs/tags/v1.0.0",
			name:   "v1.0.0",
			prefix: "refs/tags",
		},
		{
			after:  "refs/tags/v1.0.0",
			name:   "v1.0.0",
			prefix: "refs/tags/",
		},
		// branch references
		{
			after:  "refs/heads/master",
			name:   "master",
			prefix: "refs/heads",
		},
		{
			after:  "refs/heads/master",
			name:   "master",
			prefix: "refs/heads/",
		},
		// is already a ref
		{
			after:  "refs/tags/v1.0.0",
			name:   "refs/tags/v1.0.0",
			prefix: "refs/heads/",
		},
	}
	for _, test := range tests {
		if got, want := ExpandRef(test.name, test.prefix), test.after; got != want {
			t.Errorf("Got reference %s, want %s", got, want)
		}
	}
}

func TestIsTag(t *testing.T) {
	tests := []struct {
		name string
		tag  bool
	}{
		// tag references
		{
			name: "refs/tags/v1.0.0",
			tag:  true,
		},
		{
			name: "refs/heads/master",
			tag:  false,
		},
	}
	for _, test := range tests {
		if got, want := IsTag(test.name), test.tag; got != want {
			t.Errorf("Got IsTag %v, want %v", got, want)
		}
	}
}

func TestIsPullRequest(t *testing.T) {
	tests := []struct {
		name string
		tag  bool
	}{
		{
			name: "refs/pull/12/head",
			tag:  true,
		},
		{
			name: "refs/pull/12/merge",
			tag:  true,
		},
		{
			name: "refs/pull-request/12/head",
			tag:  true,
		},
		{
			name: "refs/merge-requests/12/head",
			tag:  true,
		},
		// not pull requests
		{
			name: "refs/tags/v1.0.0",
			tag:  false,
		},
		{
			name: "refs/heads/master",
			tag:  false,
		},
	}
	for _, test := range tests {
		if got, want := IsPullRequest(test.name), test.tag; got != want {
			t.Errorf("Got IsPullRequest %v, want %v", got, want)
		}
	}
}

func TestExtractPullRequest(t *testing.T) {
	tests := []struct {
		name   string
		number int
	}{
		{
			name:   "refs/pull/12/head",
			number: 12,
		},
		{
			name:   "refs/pull/12/merge",
			number: 12,
		},
		{
			name:   "refs/pull-request/12/head",
			number: 12,
		},
		{
			name:   "refs/merge-requests/12/head",
			number: 12,
		},
		{
			name:   "refs/heads/master",
			number: 0,
		},
	}
	for _, test := range tests {
		if got, want := ExtractPullRequest(test.name), test.number; got != want {
			t.Errorf("Got pull request number %v, want %v", got, want)
		}
	}
}

func TestIsHash(t *testing.T) {
	tests := []struct {
		name string
		tag  bool
	}{
		{
			name: "aacad6eca956c3a340ae5cd5856aa9c4a3755408",
			tag:  true,
		},
		{
			name: "3da541559918a808c2402bba5012f6c60b27661c",
			tag:  true,
		},
		{
			name: "f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b",
			tag:  true,
		},
		// not a sha
		{
			name: "aacad6e",
			tag:  false,
		},
		{
			name: "master",
			tag:  false,
		},
		{
			name: "refs/heads/master",
			tag:  false,
		},
		{
			name: "issue/42",
			tag:  false,
		},
		{
			name: "feature/foo",
			tag:  false,
		},
	}
	for _, test := range tests {
		if got, want := IsHash(test.name), test.tag; got != want {
			t.Errorf("Got IsHash %v, want %v", got, want)
		}
	}
}
