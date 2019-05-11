package rc

type ReplicationController struct {
	Namespace        string
	Replicas         int
	PullSecret       []string
	LivenessProbes   bool
	Containers map[string]Container
}

type Container struct {
	Image            string
	MinCPU           int
	MaxCPU           int
	MinMem           int
	MaxMem           int
}
