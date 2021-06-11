package session

import (
	"context"
	"net/http"
	"time"

	"github.com/shoppehub/conf"
	"github.com/shoppehub/fastapi/base"
	"github.com/shoppehub/fastapi/crud"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var SidKey = "sid"

var MaxAge int

var collectionName string

func init() {
	collectionName = "user_sessions"
}

// 用户session模型
type UserSession struct {
	base.BaseId `bson,inline`
	Uid         *primitive.ObjectID `bson:"uid" json:"uid"`
	Expires     *time.Time          `bson:"expires,omitempty" json:"expires,omitempty"`
	Agent       string              `bson:"agent,omitempty" json:"agent,omitempty"`
	Ip          string              `bson:"ip,omitempty" json:"ip,omitempty"`
	Status      string              `bson:"status" json:"status"`
}

func (u *UserSession) GetUid() string {
	return u.Uid.Hex()
}

func Init() {
	sid := conf.GetString("http.sid")
	if sid != "" {
		SidKey = sid
	}

	maxAge := conf.GetInt("http.maxAge")
	if maxAge != 0 {
		MaxAge = int(maxAge)
	}
}

// 获取 session
func GetUserSession(r *http.Request) *UserSession {
	var ctx = r.Context()
	val := ctx.Value(SidKey)
	if val == nil {
		return nil
	}
	return val.(*UserSession)
}

func WrapUserSession(resource *crud.Resource, r *http.Request) {

	var sid = getCookieSid(r)
	if sid == "" {
		return
	}
	var session UserSession
	resource.FindById(sid, &session, crud.FindOneOptions{CollectionName: &collectionName})
	if session.Status == "login" && (session.Expires == nil || session.Expires.After(time.Now())) {
		SetUserSession(r, &session)
	}
}

func SetUserSession(r *http.Request, session *UserSession) {
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

func NewUserSession(resource *crud.Resource, uid string, r *http.Request, w http.ResponseWriter) string {

	cookie := SaveSessionId(r, w)
	oid, _ := primitive.ObjectIDFromHex(cookie.Value)
	ouid, _ := primitive.ObjectIDFromHex(uid)

	session := UserSession{
		BaseId: base.BaseId{
			Id: &oid,
		},
		Uid:    &ouid,
		Agent:  r.UserAgent(),
		Ip:     GetIP(r),
		Status: "login",
	}

	if cookie.Expires.Second() != 0 {
		session.Expires = &cookie.Expires
	}

	resource.SaveOrUpdateOne(session, &crud.UpdateOption{
		CollectionName: &collectionName,
	})

	return cookie.Value
}

func SaveSessionId(r *http.Request, w http.ResponseWriter) *http.Cookie {

	domain := GetDomain(r.Host)

	uid := getCookieSid(r)
	if uid == "" {
		uid = primitive.NewObjectID().Hex()
	}

	cookie := http.Cookie{
		Name:     SidKey,
		Value:    uid,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
	}

	if r.TLS != nil {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	}

	newCookie(&cookie)

	http.SetCookie(w, &cookie)
	return &cookie
}

func newCookie(cookie *http.Cookie) {
	if MaxAge > 0 {
		d := time.Duration(MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else if MaxAge < 0 {
		// Set it to the past to expire now.
		cookie.Expires = time.Unix(1, 0)
	}
}

func RemoveUserSession(resource *crud.Resource, r *http.Request) {
	var sid = getCookieSid(r)
	if sid == "" {
		return
	}
	var session UserSession
	resource.FindById(sid, &session, crud.FindOneOptions{CollectionName: &collectionName})

	if session.Uid != nil {
		session.Status = "logout"
		resource.SaveOrUpdateOne(session, &crud.UpdateOption{
			CollectionName: &collectionName,
		})
	}

}
