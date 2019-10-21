package ocr

type Detecion struct {
	DetecionMap
	DetecionList
	LatexNormal string `json:"latex_normal"`
}

type DetecionMap struct {
	Chart      int `json:"contains_chart"`
	Diagram    int `json:"contains_diagram"`
	Graph      int `json:"contains_graph"`
	Table      int `json:"contains_table"`
	IsInverted int `json:"is_inverted"`
	IsPrinted  int `json:"is_printed"`
	IsBlank    float64 `json:"is_blank"`
	IsNotMath  float64 `json:"is_not_math"`
}

type DetecionList struct {
	IsPrinted int `json:"is_printed"`
}