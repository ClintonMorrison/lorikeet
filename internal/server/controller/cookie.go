package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/server/session"
	"github.com/ClintonMorrison/lorikeet/internal/utils"
)

type CookieHelper struct {
	localDev bool
}

func NewCookieHelper(localDev bool) *CookieHelper {
	return &CookieHelper{localDev}
}

const sessionCookieName = "session"

func (ch *CookieHelper) setCookieHeader(name string, value string, lifespan time.Duration) ResponseHeader {
	cookie := utils.FormatCookie(name, value, lifespan, !ch.localDev)
	return ResponseHeader{"Set-Cookie", cookie}
}

func (ch *CookieHelper) SetSessionCookieHeader(sessionToken string) ResponseHeader {
	return ch.setCookieHeader(sessionCookieName, sessionToken, session.Lifespan)
}

func (ch *CookieHelper) ClearSessionCookieHeader() ResponseHeader {
	return ch.setCookieHeader(sessionCookieName, "", 0)
}

func ParseBasicContext(r *http.Request) model.RequestContext {
	username, password, _ := r.BasicAuth()

	username = strings.TrimSpace(strings.ToLower(username))
	ip := utils.GetIpFromRequest(r)

	cookies := utils.ParseCookies(r)
	sesionToken := cookies["session"]

	return model.RequestContext{Username: username, Ip: ip, Password: password, SessionToken: sesionToken}
}
