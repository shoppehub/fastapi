package session

import (
	"context"
	"net/http"
	"time"

	"github.com/shoppehub/fastapi/base"
	"github.com/shoppehub/fastapi/crud"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var SidKey = "sid"

var MaxAge int

var collectionName = "user_sessions"

var defaultExpires = time.Unix(1, 0)

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

	cookie := saveSessionId(r, w, userSession.MaxAge)
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

	if cookie.Expires.String() == defaultExpires.String() {
		d := time.Duration(24) * time.Hour
		expires := time.Now().Add(d)
		session.Expires = &expires
	}

	resource.SaveOrUpdateOne(session, &crud.UpdateOption{
		CollectionName: &collectionName,
	})

	return &session, cookie.Value
}

func saveSessionId(r *http.Request, w http.ResponseWriter, maxAge int64) *http.Cookie {

	domain := GetDomain(r.Host)

	sid := primitive.NewObjectID().Hex()

	cookie := http.Cookie{
		Name:     SidKey,
		Value:    sid,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
	}

	if r.TLS != nil {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	}

	newCookie(maxAge, &cookie)
	http.SetCookie(w, &cookie)
	return &cookie
}

func newCookie(sessionMaxAge int64, cookie *http.Cookie) {
	if sessionMaxAge > 0 {
		d := time.Duration(sessionMaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else if MaxAge > 0 {
		d := time.Duration(MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else {
		cookie.Expires = defaultExpires
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
