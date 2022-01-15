# jiratool
Submit batch changes over one or multiple projects to Jira Cloud

# usage

| parameter | type   | mandatory | default | description                              |
|-----------|--------|-----------|---------|------------------------------------------|
| -h        | string | yes       |         | Jira Cloud Alias                         |         
| -u        | string | yes       |         | Jira Username                            |
| -a        | string | yes       |         | Jira API-Key                             |
| -p        | string | yes       |         | list of Jira projects (comma separated)  |
| -cv       | string | no        |         | create project version (release)         |
| -rv       | string | no        |         | release project version                  | 
| -rd       | string | no        | today   | release date of released project version |
