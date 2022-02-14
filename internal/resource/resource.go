package resource

type Spec struct {
	GPU                    *int64 `yaml:"gpu,omitempty"`
	GPUPerJob              *int64 `yaml:"gpuperjob,omitempty"`
	MaxMemoryPerJob        *int64 `yaml:"maxmemoryperjob,omitempty"`
	DefaultMemoryPerJob    *int64 `yaml:"defaultmemoryperjob,omitempty"`
	CPUPerJob              *int64 `yaml:"cpuperjob,omitempty"`
	StorageProxyCPURequest *int64 `yaml:"storageproxycpurequest,omitempty"`
	StorageProxyCPULimit   *int64 `yaml:"storageproxycpulimit,omitempty"`
	StorageProxyMemory     *int64 `yaml:"storageproxymemory,omitempty"`
	StorageSize            *int64 `yaml:"storagesize,omitempty"`
}

type Summary struct {
	Max  int64
	Used int64
}

type Quota struct {
	GPU     Summary
	CPU     Summary
	Memory  Summary
	Storage int64
}

type Request struct {
	GPU    int64
	CPU    int64
	Memory int64
}
