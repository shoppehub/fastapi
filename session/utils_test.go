package session

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	userInfo := map[string]interface{}{
		"avatar":   "session.Avatar",
		"nickName": "session.NickName",
		"email":    "session.Email",
	}
	userInfoStr, err := json.Marshal(userInfo)
	if err == nil {
		userInfoStr, encryptError := Encrypt(string(userInfoStr), "chemball.com", "chemball")
		if encryptError == nil {
			fmt.Println("", userInfoStr)
		}
	}
}
