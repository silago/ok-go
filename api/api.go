package api

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"strings"
	//"errors"
	//	"log"
	//"io/ioutil"
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
}

type SessionData struct {
	session_key    string
	session_secret string
}

func NewSessionData(key string, secret string) SessionData {
	return SessionData{key, secret}
}

type ErrorResponse struct {
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

type User struct {
	Uid       string `json:"uid,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	PicBase   string `json:"pic_base,omitempty"`
	//ErrorCode int    `json:"error_code,omitempty"`
	//ErrorMsg  string `json:"error_msg,omitempty"`
}

type Friends struct {
	Uids []string `json:"uids,omitempty"`
	//ErrorCode int      `json:"error_code,omitempty"`
	//ErrorMsg  string   `json:"error_msg,omitempty"`
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
	fmt.Println("u:: ", u.String())
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
		if err := decoder.Decode(errResponse); err != nil {
			return errors.New(fmt.Sprintf("%d,%s", errResponse.ErrorCode, errResponse.ErrorMsg))
		}
	}
	return err
}

func (api *Api) Auth(data SessionData) (User, error) {
	user := User{}
	params := make(map[string]string)
	params["application_key"] = api.AppId
	params["format"] = "json"
	params["method"] = "users.getCurrentUser"
	err := api.apiRequest(data, params, &user)
	return user, err
}

func (api *Api) Friends(data SessionData) ([]User, error) {
	friends := Friends{}
	if err := api.apiRequest(data, map[string]string{
		"applicationKey": api.AppId,
		"format":         "json",
		"method":         "friends.getAppUsers",
	}, &friends); err != nil {
		return nil, err
	}
	users := make([]User, len(friends.Uids), len(friends.Uids))
	if err := api.apiRequest(data, map[string]string{
		"applicationKey": api.AppId,
		"uids":           strings.Join(friends.Uids, ","),
		"format":         "json",
		"method":         "users.getInfo",
	}, &friends); err != nil {
		return nil, err
	}
	return users, nil

}
