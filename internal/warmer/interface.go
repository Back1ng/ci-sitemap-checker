package warmer

type Warmer interface {
	Process(url string) *FailedCheck
	Add(url string)
	Refresh() []FailedCheck
}
