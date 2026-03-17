package types

type UserGetEmailCodeReq struct {
	Email string `form:"email"`
}

type UserVerifyEmailAndRegisterReq struct {
	Email      string `form:"email" json:"email"`
	VerifyCode string `form:"verify_code" json:"verify_code"`
	Password   string `form:"password" json:"password"`
}

type UserLoginReq struct {
	Account  string `form:"account" json:"account"`
	Password string `form:"password" json:"password"`
}

type UserVerifyEmailAndLoginReq struct {
	Email      string `form:"email" json:"email"`
	VerifyCode string `form:"verify_code" json:"verify_code"`
}

type User struct {
	Id        string
	Email     string
	UserName  string
	Account   string
	Avatar    string
	CreatedAt int64
	UpdatedAT int64
}
