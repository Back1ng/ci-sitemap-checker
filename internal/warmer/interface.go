package warmer

type Warmer interface {
	Process(url string) error
	Add(url string)
	Refresh() error
}
