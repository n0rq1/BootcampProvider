package client

type Engineer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Dev struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Engineers []Engineer `json:"engineers"`
}

type Ops struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Engineers []Engineer `json:"engineers"`
}

type DevOps struct{
	ID    string `json:"id"`
	Dev []Dev `json:"dev"`
	Ops []Ops `json:"ops"`
}