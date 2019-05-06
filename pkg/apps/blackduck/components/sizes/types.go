package sizes

type ContainerSize struct {
	MinCPU   int
	MaxCPU   int
	MinMem   int
	MaxMem   int
}

type Size struct {
	Replica int
	Containers map[string]ContainerSize
}