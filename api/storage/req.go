package storage

type LSReq struct {
	Path string `json:"path"`
}

type StatReq struct {
	Path string `json:"path"`
}

type GetReq struct {
	Path string `json:"path"`
}

type DownloadReq struct {
	Path string `json:"path"`
}
