package shared

type GetShareCounter interface {
	GetShareCount(int, string, chan<- GetShareCounter, chan<- error)
}
