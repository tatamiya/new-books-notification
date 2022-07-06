package subject

type SubjectDecoder struct {
	Taishou map[string]string `json:"taishou"`
	Keitai  map[string]string `json:"keitai"`
	Naiyou  map[string]string `json:"naiyou"`
}
