package shared

type ShareCount interface {
	GetShareCount(int, string, chan<- ShareCount) error
}
