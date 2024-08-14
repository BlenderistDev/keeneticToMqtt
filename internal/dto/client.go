package dto

type Client struct {
	Mac    string `json:"mac"`
	Policy string `json:"policy"`
	Name   string `json:"name"`
	Permit bool   `json:"permit"`
}
