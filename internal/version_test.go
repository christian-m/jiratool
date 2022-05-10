package internal

import "testing"

func TestInspectVersion(t *testing.T) {
	testReleaseDate := "2021-04-01"
	tests := []struct {
		testcase    string
		project     Project
		versionName string
		versionData string
		err         bool
	}{
		{
			"inspect valid not released version",
			Project{
				Key: "PRJ",
				Versions: []Version{
					{
						Id:        "10001",
						Name:      "2021-02",
						ProjectId: 10000,
					},
				},
			},
			"2021-02",
			"Version 2021-02 in Projekt PRJ ist nicht released",
			false,
		},
		{
			"inspect valid released version",
			Project{
				Key: "PRJ",
				Versions: []Version{
					{
						Id:          "10001",
						Name:        "2021-02",
						Released:    true,
						ReleaseDate: &testReleaseDate,
						ProjectId:   10000,
					},
				},
			},
			"2021-02",
			"Version 2021-02 in Projekt PRJ ist released am 2021-04-01",
			false,
		},
		{
			"inspect valid archived version",
			Project{
				Key: "PRJ",
				Versions: []Version{
					{
						Id:        "10001",
						Name:      "2021-02",
						Archived:  true,
						ProjectId: 10000,
					},
				},
			},
			"2021-02",
			"Version 2021-02 in Projekt PRJ ist archiviert",
			false,
		},
		{
			"inspect valid released and archived version",
			Project{
				Key: "PRJ",
				Versions: []Version{
					{
						Id:          "10001",
						Name:        "2021-02",
						Released:    true,
						ReleaseDate: &testReleaseDate,
						Archived:    true,
						ProjectId:   10000,
					},
				},
			},
			"2021-02",
			"Version 2021-02 in Projekt PRJ ist archiviert (released am 2021-04-01)",
			false,
		},
		{
			"inspect version in project not present",
			Project{
				Key: "PRJ",
				Versions: []Version{
					{
						Id:        "10000",
						Name:      "2021-01",
						ProjectId: 10000,
					},
					{
						Id:        "10001",
						Name:      "2021-02",
						ProjectId: 10000,
					},
				},
			},
			"2021-03",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			verData, err := InspectVersion(&tt.project, tt.versionName, &TestRestClient{})
			if verData != tt.versionData {
				t.Errorf("got: %s - want: %s", verData, tt.versionData)
			}
			if err != nil && !tt.err {
				t.Errorf("got: %v - want: no error", err)
			}
		})
	}
}

func TestCreateVersion(t *testing.T) {
	tests := []struct {
		testcase    string
		project     Project
		versionName string
		err         bool
	}{
		{
			"release valid project",
			Project{
				Id: "10000",
			},
			"2021-02",
			false,
		},
		{
			"release invalid project-id",
			Project{
				Id: "CHAR_ID",
			},
			"2021-02",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			err := CreateVersion(&tt.project, tt.versionName, &TestRestClient{})
			if err != nil && !tt.err {
				t.Errorf("got: %v - want: no error", err)
			}
		})
	}
}

func TestReleaseVersion(t *testing.T) {
	tests := []struct {
		testcase       string
		project        Project
		releaseVersion string
		releaseDate    string
		err            bool
	}{
		{
			"release valid project",
			Project{
				Versions: []Version{
					{
						Id:        "10000",
						Name:      "2021-01",
						ProjectId: 10000,
					},
					{
						Id:        "10001",
						Name:      "2021-02",
						ProjectId: 10000,
					},
				},
			},
			"2021-02",
			"2021-04-01",
			false,
		},
		{
			"release version in project not present",
			Project{
				Versions: []Version{
					{
						Id:        "10000",
						Name:      "2021-01",
						ProjectId: 10000,
					},
					{
						Id:        "10001",
						Name:      "2021-02",
						ProjectId: 10000,
					},
				},
			},
			"2021-03",
			"2021-04-01",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			c := &TestRestClient{}
			err := ReleaseVersion(&tt.project, tt.releaseVersion, tt.releaseDate, c)
			if err != nil && !tt.err {
				t.Errorf("got: %v - want: no error", err)
			}
		})
	}
}

type TestRestClient struct{}

func (c *TestRestClient) GetProject(prjKey string) (*Project, error) {
	return nil, nil
}

func (c *TestRestClient) CreateVersion(version Version) error {
	return nil
}

func (c *TestRestClient) UpdateVersion(version Version) error {
	return nil
}

func TestGetVersion(t *testing.T) {
	tests := []struct {
		testcase    string
		project     Project
		versionName string
		err         bool
	}{
		{
			"version in project not present",
			Project{
				Key: "PRJ",
				Versions: []Version{
					{
						Id:        "10000",
						Name:      "2021-01",
						ProjectId: 10000,
					},
				},
			},
			"2021-03",
			true,
		},
		{
			"version in project present",
			Project{
				Key: "PRJ",
				Versions: []Version{
					{
						Id:        "10000",
						Name:      "2021-01",
						ProjectId: 10000,
					},
				},
			},
			"2021-03",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			_, err := getVersion(&tt.project, tt.versionName)
			if err != nil && !tt.err {
				t.Errorf("got: %v - want: no error", err)
			}
		})
	}
}
