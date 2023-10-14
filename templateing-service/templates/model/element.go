package model

type Element struct {
	Type      string `json:"type"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	ValueFrom string `json:"value_from"`
	Font      string `json:"font"`
	Size      int    `json:"size"`
}
