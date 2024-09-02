package chaincode

type User struct {
	UserName string `json:"UserName"`

	Role string `json:"Role"`
	Org  string `json:"Org"`

	ComputeResList []string `json:"ComputeResList"`
}
