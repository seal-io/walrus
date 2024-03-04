// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import (
	"regexp"
	"strconv"
	"strings"
)

// regular expression to extract the pull request number
// from the git ref (e.g. refs/pulls/{d}/head)
var re = regexp.MustCompile("\\d+")

// regular expressions to test whether or not a string is
// a sha1 or sha256 commit hash.
var (
	sha1   = regexp.MustCompile("^([a-f0-9]{40})$")
	sha256 = regexp.MustCompile("^([a-f0-9]{64})$")
)

// Split splits the full repository name into segments.
func Split(s string) (owner, name string) {
	parts := strings.SplitN(s, "/", 2)
	switch len(parts) {
	case 1:
		name = parts[0]
	case 2:
		owner = parts[0]
		name = parts[1]
	}
	return
}

// Join joins the repository owner and name segments to
// create a fully qualified repository name.
func Join(owner, name string) string {
	return owner + "/" + name
}

// TrimRef returns ref without the path prefix.
func TrimRef(ref string) string {
	ref = strings.TrimPrefix(ref, "refs/heads/")
	ref = strings.TrimPrefix(ref, "refs/tags/")
	return ref
}

// ExpandRef returns name expanded to the fully qualified
// reference path (e.g refs/heads/master).
func ExpandRef(name, prefix string) string {
	prefix = strings.TrimSuffix(prefix, "/")
	if strings.HasPrefix(name, "refs/") {
		return name
	}
	return prefix + "/" + name
}

// ExtractPullRequest returns name extraced pull request
// number from the reference path.
func ExtractPullRequest(ref string) int {
	s := re.FindString(ref)
	d, _ := strconv.Atoi(s)
	return d
}

// IsBranch returns true if the reference path points to
// a branch.
func IsBranch(ref string) bool {
	return strings.HasPrefix(ref, "refs/heads/")
}

// IsTag returns true if the reference path points to
// a tag object.
func IsTag(ref string) bool {
	return strings.HasPrefix(ref, "refs/tags/")
}

// IsPullRequest returns true if the reference path points
// to a pull request object.
func IsPullRequest(ref string) bool {
	return strings.HasPrefix(ref, "refs/pull/") ||
		strings.HasPrefix(ref, "refs/pull-request/") ||
		strings.HasPrefix(ref, "refs/merge-requests/")
}

// IsHash returns true if the string is a commit hash.
func IsHash(s string) bool {
	return sha1.MatchString(s) || sha256.MatchString(s)
}
