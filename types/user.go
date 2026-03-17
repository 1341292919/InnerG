package types

type UserGetEmailCodeReq struct {
	Email string `form:"email" binding:"required"`
}

type UserVerifyEmailAndRegisterReq struct {
	Email      string `form:"email" binding:"required"`
	VerifyCode string `form:"verify_code" binding:"required"`
	Password   string `form:"password" binding:"required"`
}

type UserLoginReq struct {
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserVerifyEmailAndLoginReq struct {
	Email      string `form:"email" binding:"required"`
	VerifyCode string `form:"verify_code" binding:"required"`
}
type UpdateUserAccountReq struct {
	Account string `form:"account" binding:"required"`
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
