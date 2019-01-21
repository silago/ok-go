package api

type AuthProvider interface {
	Friends() ( []string, error )
	FriendsData() ([]User, error)
	Auth() (User, error)
}

type User struct {
	Uid       string `json:"uid,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	PicBase   string `json:"pic_base,omitempty"`
}

//func (api *Api) Friends(data SessionData) ([]string, error) {
//func (api *Api) FriendsData(data SessionData) ([]User, error) {
//unc (api *Api) Auth(data SessionData) (User, error) {
