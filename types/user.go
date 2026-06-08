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
type UpdateUserSignatureReq struct {
	Signature string `form:"signature" json:"signature"`
}
type UpdateUserProfilePublicReq struct {
	ProfilePublic *int8 `form:"profile_public" json:"profile_public" binding:"required,oneof=0 1"`
}
type GetUserInfoByIDReq struct {
	UserID int64 `form:"user_id" json:"user_id" binding:"required"`
}
type UpdateUserAvatarReq struct {
}
type UpdateUserAvatarResp struct {
	AvatarUrl string
}
type User struct {
	Id            string
	Email         string
	UserName      string
	Account       string
	Avatar        string
	Signature     string
	ProfilePublic int8
	Gender        int
	RoleType      int
	CreatedAt     int64
	UpdatedAT     int64
}

type UserProfile struct {
	Id        string
	UserName  string
	Account   string
	Avatar    string
	Signature string
	Protected bool
	Gender    int
	RoleType  int
	CreatedAt int64
	UpdatedAT int64
}
