package models

type ProjectCreate struct {
	Domain     string
	EntryPoint string
	WebServer  int
	Database   int
	Cache      int
}

func NewProjectCreate(d string, ep string, ws int, db int, c int) *ProjectCreate {
	return &ProjectCreate{
		Domain:     d,
		EntryPoint: ep,
		WebServer:  ws,
		Database:   db,
		Cache:      c,
	}
}
