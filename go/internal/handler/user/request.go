package user

type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Timezone  string `json:"timezone"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Timezone  string `json:"timezone"`
}
