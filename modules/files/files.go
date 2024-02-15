package files

type DeleteFilesReq struct {
	Destinations []string `json:"destinations" form:"destinations"`
}
