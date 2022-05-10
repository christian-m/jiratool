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

func InspectVersion(prj *Project, verName string, c RestClient) (string, error) {
	ver, err := getVersion(prj, verName)
	if err != nil {
		return "", err
	}
	var verData string
	if ver.Archived && ver.Released {
		verData = fmt.Sprintf("Version %s in Projekt %s ist archiviert (released am %s)", ver.Name, prj.Key, *ver.ReleaseDate)
	} else if ver.Archived {
		verData = fmt.Sprintf("Version %s in Projekt %s ist archiviert", ver.Name, prj.Key)
	} else if ver.Released {
		verData = fmt.Sprintf("Version %s in Projekt %s ist released am %s", ver.Name, prj.Key, *ver.ReleaseDate)
	} else {
		verData = fmt.Sprintf("Version %s in Projekt %s ist nicht released", ver.Name, prj.Key)
	}
	return verData, nil
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
	ver, err2 := getVersion(prj, relVer)
	if err2 != nil {
		return err2
	}

	ver.ReleaseDate = &relDate
	ver.Released = true
	ver.UserReleaseDate = nil
	err := c.UpdateVersion(*ver)
	return err
}

func getVersion(prj *Project, relVer string) (*Version, error) {
	var ver *Version = nil
	for _, v := range prj.Versions {
		if v.Name == relVer {
			ver = &v
			break
		}
	}
	if ver == nil {
		return nil, fmt.Errorf("Version %s ist in Projekt %s nicht vorhanden", relVer, prj.Key)
	}
	return ver, nil
}
