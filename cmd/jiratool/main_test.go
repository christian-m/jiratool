package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestCreateUserInfo(t *testing.T) {
	tests := []struct {
		testcase string
		username string
		apiKey   string
		userinfo *url.Userinfo
		err      bool
	}{
		{
			"valid credentials",
			"username",
			"apikey",
			url.UserPassword("username", "apikey"),
			false,
		},
		{
			"missing username",
			"",
			"apikey",
			nil,
			true,
		},
		{
			"missing apiKey",
			"username",
			"",
			nil,
			true,
		},
		{
			"no credentials provided",
			"",
			"",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			userinfo, err := createUserInfo(tt.username, tt.apiKey)
			switch {
			case err != nil && !tt.err:
				t.Errorf("got: %v - want: no Error", err)
			case !reflect.DeepEqual(userinfo, tt.userinfo):
				t.Errorf("Got: %v - Want: %v", userinfo, tt.userinfo)
			}
		})
	}
}

func TestResolveProjects(t *testing.T) {
	tests := []struct {
		testcase      string
		projectString string
		projects      []string
		err           bool
	}{
		{
			"one valid project",
			"DB",
			[]string{"DB"},
			false,
		},
		{
			"some valid projects",
			"DB,MN,REL",
			[]string{"DB", "MN", "REL"},
			false,
		},
		{
			"missing projects",
			"",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			projects, err := resolveProjects(tt.projectString)
			switch {
			case err != nil && !tt.err:
				t.Errorf("got: %v - want: no Error", err)
			case !reflect.DeepEqual(projects, tt.projects):
				t.Errorf("Got: %v - Want: %v", projects, tt.projects)
			}
		})
	}
}
