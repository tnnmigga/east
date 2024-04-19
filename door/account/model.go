package account

type Account struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type WebResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
