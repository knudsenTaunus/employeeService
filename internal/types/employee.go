package types

// Employee is the struct used for response
type Employee struct {
	ID        int `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Salary    int `json:"salary"`
}
