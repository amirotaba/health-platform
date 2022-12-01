package entity

type Identity struct {
	ID   int64      `json:"id,omitempty"`
	Role RoleEntity `json:"roles,omitempty"`
}

type SingInRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SingUpRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}
