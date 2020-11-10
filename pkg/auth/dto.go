package auth

type RegisterDto struct {
	Name            string `validate:"required,alpha"`
	Email           string `validate:"required,email"`
	Password        string `validate:"required,min=5"`
	ConfirmPassword string `validate:"required,eqfield=Password"`
}
