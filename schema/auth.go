package schema

type ReqRegister struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}
type ResRegister struct {
	Message string `json:"message"`
}

type ReqLogin struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

type ResLogin struct {
	AccessToken string `json:"access_token"`
}

type ResLogout struct {
	Message string `json:"message"`
}
