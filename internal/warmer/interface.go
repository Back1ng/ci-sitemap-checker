package warmer

type Warmer interface {
	Process(url string) *FailedCheck
	Refresh()
}
