package sizes

type SizeInterface interface {
	GetSize(name string) map[string]*Size
}