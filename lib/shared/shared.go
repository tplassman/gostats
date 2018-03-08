package shared

type ShareCount struct {
	Index  int
	Count  int
	Source string
}

type GetShareCounter interface {
	GetShareCount(string) (int, error)
}
