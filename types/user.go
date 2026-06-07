package types

type UserGetEmailCodeReq struct {
	Email string `form:"email" json:"email" binding:"required"`
}

type UserVerifyEmailAndRegisterReq struct {
	Email      string `form:"email" json:"email" binding:"required"`
	VerifyCode string `form:"verify_code" json:"verify_code" binding:"required"`
	Password   string `form:"password" json:"password" binding:"required"`
}

type UserLoginReq struct {
	Account  string `form:"account" json:"account" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserVerifyEmailAndLoginReq struct {
	Email      string `form:"email" json:"email" binding:"required"`
	VerifyCode string `form:"verify_code" json:"verify_code" binding:"required"`
}
type UpdateUserAccountReq struct {
	Account string `form:"account" json:"account" binding:"required"`
}
type UpdateUserNameReq struct {
	UserName string `form:"username" json:"username" binding:"required"`
}
type UpdateUserGenderReq struct {
	Gender string `form:"gender" json:"gender" binding:"required"`
}
type UpdateUserAvatarReq struct {
}
type UpdateUserAvatarResp struct {
	AvatarUrl string
}
type User struct {
	Id        string
	Email     string
	UserName  string
	Account   string
	Avatar    string
	Gender    int
	RoleType  int
	CreatedAt int64
	UpdatedAT int64
}
