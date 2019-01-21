package api

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"fmt"
	"net/http"
	"net/url"
	"sort"
)

//import "github.com/Jeffail/gabs"

const API_HOST = "api.ok.ru"
const PATH = "fb.do"

type Api struct {
	AppId string
	SessionData SessionData
}

type SessionData struct {
	session_key    string
	session_secret string
}



func NewOkAuthProvider(app_id string, session_key string, session_secret_key string) AuthProvider {
	return &Api{AppId: app_id, SessionData: SessionData{session_key: session_key,session_secret:session_secret_key}}
}

func NewSessionData(key string, secret string) SessionData {
	return SessionData{key, secret}
}

type ErrorResponse struct {
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

type Friends struct {
	Uids []string `json:"uids,omitempty"`
}

//func NewApi(AppId ) Api {
//    return Api{
//
//    }
//
//}


func (api *Api) makeSig(session SessionData, params map[string]string) string {
	var buffer bytes.Buffer
	params["session_key"] = session.session_key + session.session_secret
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		appendix := fmt.Sprintf("%s=%s", key, params[key])
		buffer.WriteString(appendix)
	}
	h := md5.New()
	h.Write(buffer.Bytes())
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (api *Api) apiRequest(session SessionData, params map[string]string, obj interface{}) error {
	u := url.URL{}

	u.Scheme = "http"
	u.Path = PATH
	u.Host = API_HOST
	query := u.Query()

	params["sig"] = api.makeSig(session, params)
	params["session_key"] = session.session_key
	fmt.Println("SKEY", params["session_key"])

	for key, value := range params {
		query.Add(key, value)
	}
	u.RawQuery = query.Encode()

	httpClient := &http.Client{}
	log.Println("u:: ", u.String())
	res, err := httpClient.Get(u.String())
	if err != nil {
		return err
	}

	//bodyBytes, _ := ioutil.ReadAll(res.Body)
	//bodyString := string(bodyBytes)
	//log.Println(bodyString)

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(obj); err != nil {
		errResponse := &ErrorResponse{}
		if err := decoder.Decode(errResponse); err == nil {
			return errors.New(fmt.Sprintf("%d,%s", errResponse.ErrorCode, errResponse.ErrorMsg))
		} else {
			return err
		}
	}
	log.Println("obj: ", obj)
	return err
}

func (api *Api) Auth() (User, error) {
	user := User{}
	params := make(map[string]string)
	params["application_key"] = api.AppId
	params["format"] = "json"
	params["method"] = "users.getCurrentUser"
	params["fields"] = "uid,first_name,last_name,pic_base"
	err := api.apiRequest(api.SessionData, params, &user)
	return user, err
}

func (api *Api) FriendsData() ([]User, error) {
	friends := Friends{}
	if err := api.apiRequest(api.SessionData, map[string]string{
		"application_key": api.AppId,
		"format":          "json",
		"method":          "friends.getAppUsers",
	}, &friends); err != nil {
		return nil, err
	}
	log.Printf("uids", friends.Uids)
	users := make([]User, len(friends.Uids), len(friends.Uids))
	err := api.apiRequest(api.SessionData, map[string]string{
		"application_key": api.AppId,
		"uids":            strings.Join(friends.Uids, ","),
		"format":          "json",
		"fields":          "uid,first_name,last_name,pic_base",
		"method":          "users.getInfo",
	}, &friends)

	return users, err
}

func (api *Api) Friends() ([]string, error) {
	friends := Friends{}
	err := api.apiRequest(api.SessionData, map[string]string{
		"application_key": api.AppId,
		"format":          "json",
		"method":          "friends.getAppUsers",
	}, &friends)
	return friends.Uids, err
}
