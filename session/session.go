package session

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/shoppehub/fastapi/base"
	"github.com/shoppehub/fastapi/crud"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var SidKey = "sid"
var ChemballUserKey = "c_u"

var MaxAge int

var collectionName = "user_sessions"

// 获取 session 的登录id
func GetUserId(r *http.Request) string {
	s := GetUserSession(r)
	if s != nil {
		return s.Uid
	}
	return ""
}

// 获取dession
func GetUserSession(r *http.Request) *UserSession {
	var ctx = r.Context()
	val := ctx.Value(SidKey)
	if val == nil {
		return &UserSession{}
	}
	return val.(*UserSession)
}

func wrapUserSession(resource *crud.Resource, r *http.Request) {

	var sid = getCookieSid(r)
	if sid == "" {
		return
	}
	var session UserSession
	resource.FindById(sid, &session, crud.FindOneOptions{CollectionName: &collectionName})
	if session.Status == "login" && (session.Expires.After(time.Now())) {

		setUserSession(r, &session)
	}
}

func setUserSession(r *http.Request, session *UserSession) {
	var ctx = r.Context()
	*r = *r.WithContext(context.WithValue(ctx, SidKey, session))
}

func getCookieSid(r *http.Request) string {
	var sid string
	if c, errCookie := r.Cookie(SidKey); errCookie == nil {
		sid = c.Value
	}
	return sid
}

func NewUserSession(resource *crud.Resource, userSession UserSession, r *http.Request, w http.ResponseWriter) (*UserSession, string) {

	cookie := writeCookie(r, w, userSession)

	oid, _ := primitive.ObjectIDFromHex(cookie.Value)

	session := UserSession{
		BaseId: base.BaseId{
			Id: &oid,
		},
		Uid:      userSession.Uid,
		Agent:    r.UserAgent(),
		Expires:  &cookie.Expires,
		Ip:       GetIP(r),
		Status:   "login",
		Avatar:   userSession.Avatar,
		NickName: userSession.NickName,
	}

	if cookie.Expires.IsZero() {
		d := time.Duration(24) * time.Hour
		expires := time.Now().Add(d)
		session.Expires = &expires
	}

	resource.SaveOrUpdateOne(session, &crud.UpdateOption{
		CollectionName: &collectionName,
	})

	return &session, cookie.Value
}

func writeCookie(r *http.Request, w http.ResponseWriter, session UserSession) *http.Cookie {

	domain := GetDomain(r.Host)

	sid := primitive.NewObjectID().Hex()

	// sid 是服务端判断 sessionId
	sidCookie := http.Cookie{
		Name:     SidKey,
		Value:    sid,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
	}
	if r.TLS != nil {
		sidCookie.Secure = true
		sidCookie.SameSite = http.SameSiteNoneMode
	}
	setCookieExpires(&sidCookie, session.MaxAge)
	http.SetCookie(w, &sidCookie)

	userInfo := map[string]interface{}{
		"avatar":   session.Avatar,
		"nickName": session.NickName,
		"email":    session.Email,
	}
	userInfoStr, err := json.Marshal(userInfo)
	if err == nil {
		logrus.Error("Error: ", err)
		userInfoStr, encryptError := Encrypt(string(userInfoStr), sid, "chemball")
		if encryptError == nil {
			userCookie := http.Cookie{
				Name:   ChemballUserKey,
				Value:  userInfoStr,
				Path:   "/",
				Domain: domain,
			}
			if r.TLS != nil {
				userCookie.Secure = true
				userCookie.SameSite = http.SameSiteNoneMode
			}
			setCookieExpires(&userCookie, session.MaxAge)
			http.SetCookie(w, &userCookie)
		}
	} else {
		logrus.Error("json.Marshal Error: ", err)
	}

	return &sidCookie
}

/*设置cookie的过期时间*/
func setCookieExpires(cookie *http.Cookie, sessionMaxAge int64) {
	if sessionMaxAge > 0 {
		d := time.Duration(sessionMaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else if MaxAge > 0 {
		d := time.Duration(MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	}
}

// 退出登录
func RemoveUserSession(resource *crud.Resource, r *http.Request) {
	var sid = getCookieSid(r)
	if sid == "" {
		return
	}
	var session UserSession
	resource.FindById(sid, &session, crud.FindOneOptions{CollectionName: &collectionName})

	session.Status = "logout"
	resource.SaveOrUpdateOne(session, &crud.UpdateOption{
		CollectionName: &collectionName,
	})
}
