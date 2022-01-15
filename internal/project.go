package internal

type Project struct {
	Id          string    `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Versions    []Version `json:"versions"`
}
