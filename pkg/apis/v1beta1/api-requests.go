package apis

type APIEncryptReq struct {
	EncryptString string `json:"encrypt_string" binding:"required"`
}

type APIReq struct {
	AuthorizedKey string `json:"authorized_key" binding:"required"`
	Payload       string `json:"payload" binding:"required"`
}
