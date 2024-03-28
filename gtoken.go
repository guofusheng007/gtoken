package gtoken

import (
	"encoding/json"
	"math/rand"
	"time"
)

type Gtoken struct {
	UserID     int         `json:"userid"`
	UserName   string      `json:"username"`
	UserEmail  string      `json:"email"`
	UserMobile string      `json:"mobile"`
	ExpiresAt time.Time    `json:"expire"`
	IssuedAt time.Time     `json:"issuedat"`
	Issuer string          `json:"issuer"`
	Subject string         `json:"subject"`
	TokenID string         `json:"tokenid"`  //动态随机字串,12个随机字母组成的字串，在校验token时识别客户端提交的token和tokenID是否一致。
}


//creae token encText
func CreateToken(gtoken *Gtoken,key []byte) (string){
	token, _ := json.Marshal(gtoken)
	enc,_ := EncryptCBC(string(token),key)
	return enc
}

//解析token
//其中返回值int是token TTL过期时间(分钟)，零时表示已过期。
func CheckToken(encToken string, key []byte) (*Gtoken,int64,error) {
	var token Gtoken
    dec_text,err := DecryptCBC(encToken,key)
	if err != nil {
		return nil,0,err
	} else {
		//fmt.Println(dec_text)
	    json.Unmarshal([]byte(dec_text), &token)
	    t := int64(token.ExpiresAt.Sub(time.Now()).Minutes())
	    if (t <= 0) {
		    t = 0
	    }
	    return &token,t,nil
	}
}

//产生随机字串，供tokenID使用
func RandomString(length int) string {
    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    b := make([]rune, length)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}