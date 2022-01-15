package internal

import (
	"fmt"
	"strconv"
)

type Version struct {
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	Archived        bool    `json:"archived"`
	Released        bool    `json:"released"`
	ReleaseDate     *string `json:"releaseDate"`
	UserReleaseDate *string `json:"userReleaseDate"`
	ProjectId       int     `json:"projectId"`
}

func CreateVersion(prj *Project, verName string, c RestClient) error {
	prjId, err := strconv.Atoi(prj.Id)
	if err != nil {
		return fmt.Errorf("Projekt-Id %s ist ungültig", prj.Id)
	}
	ver := Version{
		Name:      verName,
		Archived:  false,
		Released:  false,
		ProjectId: prjId,
	}
	err = c.CreateVersion(ver)
	return err
}

func ReleaseVersion(prj *Project, relVer, relDate string, c RestClient) error {
	var ver *Version = nil
	for _, v := range prj.Versions {
		if v.Name == relVer {
			ver = &v
			break
		}
	}
	if ver == nil {
		return fmt.Errorf("Version %s ist in Projekt %s nicht vorhanden", relVer, prj.Key)
	}

	ver.ReleaseDate = &relDate
	ver.Released = true
	ver.UserReleaseDate = nil
	err := c.UpdateVersion(*ver)
	return err
}
