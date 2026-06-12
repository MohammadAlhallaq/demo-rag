package types

type Doc struct {
	ID     string    `json:"id"`
	Source string    `json:"source"`
	Text   string    `json:"text"`
	Vec    []float64 `json:"vec"`
}

type Scored struct {
	Doc   Doc
	Score float64
}

type Index struct {
	Docs []Doc `json:"docs"`
}
