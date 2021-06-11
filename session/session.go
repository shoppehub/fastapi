package session

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shoppehub/conf"
	"github.com/shoppehub/fastapi/base"
	"github.com/shoppehub/fastapi/crud"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Sid = "sid"

var MaxAge int

// 用户session模型
type UserSession struct {
	base.BaseId `bson,inline`
	Uid         primitive.ObjectID `bson:"uid" json:"uid"`
	Expires     int                `bson:"expires" json:"expires"`
	Agent       string             `bson:"agent" json:"agent"`
	Ip          string             `bson:"ip" json:"ip"`
}

func Init() {
	sid := conf.GetString("http.sid")
	if sid != "" {
		Sid = sid
	}

	maxAge := conf.GetInt("http.maxAge")
	if maxAge != 0 {
		MaxAge = int(maxAge)
	}
}

func NewUserSession(resource *crud.Resource, uid string, r *http.Request, w *http.ResponseWriter) {
	var sid string
	if c, errCookie := r.Cookie(Sid); errCookie == nil {
		sid = c.Value
		if !primitive.IsValidObjectID(sid) {
			sid = ""
		}
	}
	var cookie *http.Cookie
	if sid == "" {
		cookie = SaveSessionId(r, w)
	}
	oid, _ := primitive.ObjectIDFromHex(uid)
	session := UserSession{
		Expires: cookie.MaxAge,
		Uid:     oid,
	}

	collectionName := "user_sessions"
	resource.SaveOrUpdateOne(session, &crud.UpdateOption{
		CollectionName: &collectionName,
	})

}

func SaveSessionId(r *http.Request, w *http.ResponseWriter) *http.Cookie {
	domain := GetDomain(r.Host)
	uid := primitive.NewObjectID().Hex()

	cookie := http.Cookie{
		Name:     Sid,
		Value:    uid,
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	if r.URL.Scheme == "http" {
		cookie.Secure = false
	}

	NewCookie(&cookie)

	http.SetCookie(*w, &cookie)
	return &cookie
}

func NewCookie(cookie *http.Cookie) {
	if MaxAge > 0 {
		d := time.Duration(MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else if MaxAge < 0 {
		// Set it to the past to expire now.
		cookie.Expires = time.Unix(1, 0)
	}
}

// 获取域名
func GetDomain(host string) string {
	host = strings.TrimSpace(host)
	hostParts := strings.Split(host, ".")
	lengthOfHostParts := len(hostParts)

	if lengthOfHostParts == 1 {
		return hostParts[0] // scenario C
	} else {
		_, err := strconv.ParseFloat(hostParts[0], 64)
		if err == nil {
			return host
		} else {
			return strings.Join([]string{hostParts[lengthOfHostParts-2], hostParts[lengthOfHostParts-1]}, ".")
		}
	}
}
