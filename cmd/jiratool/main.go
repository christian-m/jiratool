package main

import (
	"bitbucket.org/christian_m/jiratool/internal"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	layoutISO = "2006-01-02"
)

var (
	flagUsername       = flag.String("u", "", "Jira Username")
	flagApiKey         = flag.String("a", "", "Jira API-Key")
	flagCloudAlias     = flag.String("h", "", "Jira Cloud Alias")
	flagProjects       = flag.String("p", "", "Jira Projekte (kommasepariert)")
	flagCreateVersion  = flag.String("cv", "", "Projektversion anlegen")
	flagReleaseVersion = flag.String("rv", "", "Projektversion Release")
	flagReleaseDate    = flag.String("rd", "", "Projektversion Release Datum")
)

func main() {
	flag.Parse()
	u, err := createUserInfo(*flagUsername, *flagApiKey)
	if err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *flagCloudAlias == "" {
		fmt.Println("Bitte den Jira Cloud Alias angeben:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	c, err := internal.CreateRestClient(u, &url.URL{Scheme: "https", Host: fmt.Sprintf("%s.atlassian.net", *flagCloudAlias)})
	if err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	prjKeys, err := resolveProjects(*flagProjects)
	if err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		os.Exit(1)
	}
	for _, pk := range prjKeys {
		prj, err := c.GetProject(pk)
		if err != nil {
			t, ok := err.(internal.RestError)
			if ok && t.Status() == http.StatusNotFound {
				log.Printf("Projekt %s in Jira nicht vorhanden", pk)
			} else {
				log.Printf("Projekt %s kann nicht gelesen werden (%s)", pk, err)
			}
			continue
		}
		switch {
		case *flagCreateVersion != "":
			ver := *flagCreateVersion
			err := internal.CreateVersion(prj, ver, c)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Version %s in Projekt %s angelegt", ver, prj.Key)
			}
		case *flagReleaseVersion != "" && *flagReleaseDate != "":
			ver := *flagReleaseVersion
			relDate := *flagReleaseDate
			_, err := time.Parse(layoutISO, relDate)
			if err != nil {
				log.Printf("Das Release Datum '%s' hat nicht das richtige Format (JJJJ-MM-TT)", relDate)
				break
			}
			err = internal.ReleaseVersion(prj, ver, relDate, c)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Version %s in Projekt %s released", ver, prj.Key)
			}
		default:
			log.Printf("In Projekt %s nichts geändert", prj.Key)
		}
	}
}

func createUserInfo(username, apiKey string) (*url.Userinfo, error) {
	if username == "" || apiKey == "" {
		return nil, fmt.Errorf("Bitte Jira-Usernamen und Passwort angeben:")
	}
	return url.UserPassword(username, apiKey), nil
}

func resolveProjects(projectKeys string) ([]string, error) {
	if projectKeys == "" {
		return nil, fmt.Errorf("Bitte mindestens ein Jira-Projekt angeben")
	}
	return strings.Split(projectKeys, ","), nil
}
