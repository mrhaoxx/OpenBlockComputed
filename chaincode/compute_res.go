package main

type ComputeRes struct {
	Id string `json:"Id"`

	State string `json:"State"`

	OwnerOrg string `json:"OwnerOrg"`
	UserOrg  string `json:"UserOrg"`

	UserOrgDueDate int `json:"UserOrgDueDate"`

	User        string `json:"User"`
	UserDueDate int    `json:"UserDueDate"`

	CPUSKU string `json:"CPUSKU"`
	CPUNum int    `json:"CPUNum"`

	GPUSKU string `json:"GPUSKU"`
	GPUNum int    `json:"GPUNum"`

	MemorySize int `json:"MemorySize"`

	ConnectionAbilities string `json:"ConnectionAbilities"`
}

func (c *ComputeRes) IsAvailable() bool {
	return c.User == ""
}
