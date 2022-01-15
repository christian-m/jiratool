package internal

import "testing"

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
