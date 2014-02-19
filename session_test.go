package revel

import (
	"net/http"
	"testing"
	"time"
)

func TestSessionRestore(t *testing.T) {
	expireAfterDuration = 0
	originSession := make(Session)
	originSession["user"] = "Tom"
	originSession["info"] = "hi"
	cookie := originSession.cookie()
	if !cookie.Expires.IsZero() {
		t.Error("incorrect cookie expire", cookie.Expires)
	}

	restoredSession := getSessionFromCookie(cookie)
	for k, v := range originSession {
		if restoredSession[k] != v {
			t.Errorf("session restore failed session[%s] != %s", k, v)
		}
	}
}

func TestSessionExpire(t *testing.T) {
	expireAfterDuration = time.Hour
	session := make(Session)
	session["user"] = "Tom"
	var cookie *http.Cookie
	for i := 0; i < 3; i++ {
		cookie = session.cookie()
		time.Sleep(time.Second)
		session = getSessionFromCookie(cookie)
	}
	expectExpire := time.Now().Add(expireAfterDuration)
	if cookie.Expires.Unix() < expectExpire.Add(-time.Second).Unix() {
		t.Error("expect expires", cookie.Expires, "after", expectExpire.Add(-time.Second))
	}
	if cookie.Expires.Unix() > expectExpire.Unix() {
		t.Error("expect expires", cookie.Expires, "before", expectExpire)
	}

	session.SetSessionNoExpiration()
	for i := 0; i < 3; i++ {
		cookie = session.cookie()
		session = getSessionFromCookie(cookie)
	}
	cookie = session.cookie()
	if !cookie.Expires.IsZero() {
		t.Error("expect cookie expires is zero")
	}

	session.SetSessionDefaultExpiration()
	cookie = session.cookie()
	expectExpire = time.Now().Add(expireAfterDuration)
	if cookie.Expires.Unix() < expectExpire.Add(-time.Second).Unix() {
		t.Error("expect expires", cookie.Expires, "after", expectExpire.Add(-time.Second))
	}
	if cookie.Expires.Unix() > expectExpire.Unix() {
		t.Error("expect expires", cookie.Expires, "before", expectExpire)
	}
}
