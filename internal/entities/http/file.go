package http

const (
	DefaultType    = "core"
	DefaultVersion = "1.0.0"
	DefaultHash    = "null"
)

type FileRequest struct {
	Type    string `json:"type,omitempty"`
	Version string `json:"version,omitempty"`
	Hash    string `json:"hash,omitempty"`
}

func (freq *FileRequest) Validate() {
	if freq.Type == "" {
		freq.Type = DefaultType
	}

	if freq.Version == "" {
		freq.Version = DefaultVersion
	}

	if freq.Hash == "" {
		freq.Hash = DefaultHash
	}
}

func (freq *FileRequest) ToPath() string {
	return freq.Type + "/" + freq.Version + ".json"
}

type FileResponse struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
	Content string `json:"content"`
}
