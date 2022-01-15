package internal

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestCreateRestClient(t *testing.T) {
	tests := []struct {
		testcase string
		testurl  *url.URL
		userinfo *url.Userinfo
		client   *http.Client
		url      string
		err      bool
	}{
		{
			"valid hostname with userinfo",
			&url.URL{Scheme: "https", Host: "rest.test.de"},
			url.UserPassword("username", "apikey"),
			http.DefaultClient,
			"https://username:apikey@rest.test.de/",
			false,
		},
		{
			"valid hostname without userinfo",
			&url.URL{Scheme: "https", Host: "rest.test.de"},
			nil,
			http.DefaultClient,
			"https://rest.test.de/",
			false,
		},
		{
			"missing hostname",
			nil,
			url.UserPassword("username", "apikey"),
			http.DefaultClient,
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			c, err := CreateRestClient(tt.userinfo, tt.testurl)
			switch {
			case err != nil && !tt.err:
				t.Errorf("got: %v - want: no Error", err)
			case err == nil && c.HttpClient != tt.client:
				t.Errorf("got: %v - want: %v", c.HttpClient, tt.client)
			case err == nil && c.BaseURL.ResolveReference(&url.URL{Path: "/"}).String() != tt.url:
				t.Errorf("got: %v - want: %v", c.BaseURL.ResolveReference(&url.URL{Path: "/"}), tt.url)
			}
		})
	}
}

func TestRestClient_GetProject(t *testing.T) {
	releaseDate := "2021-07-06"
	userReleaseDate := "6/Jul/2021"
	tests := []struct {
		testcase       string
		path           string
		response       []byte
		responseStatus int
		projectString  string
		project        Project
		err            bool
	}{
		{
			"get valid project without version",
			"/rest/api/3/project/DB",
			[]byte("{\"id\": \"10000\",\"key\": \"DB\",\"description\": \"This project was created as an test for REST.\",\"url\": \"https://www.example.com\",\"email\": \"from-jira@example.com\",\"assigneeType\": \"PROJECT_LEAD\",\"versions\": [],\"name\": \"Example\"}"),
			http.StatusOK,
			"DB",
			Project{
				Id:          "10000",
				Key:         "DB",
				Description: "This project was created as an test for REST.",
				Versions:    nil,
			},
			false,
		},
		{
			"get valid project with version",
			"/rest/api/3/project/DB",
			[]byte("{\"id\": \"10000\",\"key\": \"DB\",\"description\": \"This project was created as an test for REST.\",\"url\": \"https://www.example.com\",\"email\": \"from-jira@example.com\",\"assigneeType\": \"PROJECT_LEAD\",\"versions\": [{\"id\": \"10000\",\"description\": \"An excellent version\",\"name\": \"Test Version\",\"archived\": false,\"released\": true,\"releaseDate\": \"2021-07-06\",\"userReleaseDate\": \"6/Jul/2021\",\"projectId\": 10000}],\"name\": \"Example\"}"),
			http.StatusOK,
			"DB",
			Project{
				Id:          "10000",
				Key:         "DB",
				Description: "This project was created as an test for REST.",
				Versions: []Version{
					{
						Id:              "10000",
						Name:            "Test Version",
						Archived:        false,
						Released:        true,
						ReleaseDate:     &releaseDate,
						UserReleaseDate: &userReleaseDate,
						ProjectId:       10000,
					},
				},
			},
			false,
		},
		{
			"project not found",
			"/rest/api/3/project/DB",
			[]byte("{}"),
			http.StatusNotFound,
			"DB",
			Project{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.URL.String() != tt.path {
					t.Errorf("got: %v - want: %v", req.URL.String(), tt.path)
				}
				rw.WriteHeader(tt.responseStatus)
				rw.Write(tt.response)
			}))
			defer server.Close()

			u, _ := url.Parse(server.URL)
			c, _ := CreateRestClient(nil, u)
			prj, err := c.GetProject(tt.projectString)
			switch {
			case err != nil && !tt.err:
				t.Errorf("got: %v - want: no Error", err)
			case reflect.DeepEqual(prj, tt.project):
				t.Errorf("got: %v - want: %v", prj, tt.project)
			}
		})
	}
}

func TestRestClient_CreateVersion(t *testing.T) {
	tests := []struct {
		testcase       string
		path           string
		response       []byte
		responseStatus int
		version        string
		err            bool
	}{
		{
			"can create version",
			"/rest/api/3/version",
			[]byte("{\"id\": \"10000\",\"description\": \"An excellent version\",\"name\": \"Test Version\",\"archived\": false,\"released\": true,\"releaseDate\": \"2010-07-06\",\"userReleaseDate\": \"6/Jul/2010\",\"projectId\": 10000}"),
			http.StatusCreated,
			"2021-0X",
			false,
		},
		{
			"cannot create version",
			"/rest/api/3/version",
			[]byte("{}"),
			http.StatusBadRequest,
			"2021-0X",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.URL.String() != tt.path {
					t.Errorf("got: %v - want: %v", req.URL.String(), tt.path)
				}
				rw.WriteHeader(tt.responseStatus)
				rw.Write(tt.response)
			}))
			defer server.Close()

			u, _ := url.Parse(server.URL)
			c, _ := CreateRestClient(nil, u)
			err := c.CreateVersion(Version{
				Name:      tt.version,
				Archived:  false,
				Released:  false,
				ProjectId: 0,
			})
			if err != nil && !tt.err {
				t.Errorf("got: %v - want: no Error", err)
			}
		})
	}
}

func TestRestClient_UpdateVersion(t *testing.T) {
	tests := []struct {
		testcase       string
		path           string
		response       []byte
		responseStatus int
		versionId      string
		version        string
		err            bool
	}{
		{
			"can update version",
			"/rest/api/3/version/10000",
			[]byte("{\"id\": \"10000\",\"description\": \"An excellent version\",\"name\": \"Test Version\",\"archived\": false,\"released\": true,\"releaseDate\": \"2010-07-06\",\"userReleaseDate\": \"6/Jul/2010\",\"projectId\": 10000}"),
			http.StatusOK,
			"10000",
			"2021-0X",
			false,
		},
		{
			"cannot update version",
			"/rest/api/3/version/10000",
			[]byte("{}"),
			http.StatusBadRequest,
			"10000",
			"2021-0X",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.URL.String() != tt.path {
					t.Errorf("got: %v - want: %v", req.URL.String(), tt.path)
				}
				rw.WriteHeader(tt.responseStatus)
				rw.Write(tt.response)
			}))
			defer server.Close()

			u, _ := url.Parse(server.URL)
			c, _ := CreateRestClient(nil, u)
			releaseDate := "2021-07-06"
			err := c.UpdateVersion(Version{
				Id:          tt.versionId,
				Name:        tt.version,
				Archived:    false,
				Released:    true,
				ReleaseDate: &releaseDate,
				ProjectId:   10000,
			})
			if err != nil && !tt.err {
				t.Errorf("got: %v - want: no Error", err)
			}
		})
	}
}
