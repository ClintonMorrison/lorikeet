package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// See https://developers.google.com/recaptcha/docs/verify

type RecaptchaClient struct {
	secret string
}

type recaptchaValidationResponse struct {
	Success            bool   `json:"success"`
	ChallengeTimestamp string `json:"challenge_ts"`
	Hostname           string `json:"hostname"`
}

// VerifyRecaptcha returns true if recaptcha is valid
func (rc *RecaptchaClient) Verify(recaptchaResponse string, ip string) bool {
	data := url.Values{
		"secret":   {rc.secret},
		"response": {recaptchaResponse},
		"remoteip": {ip},
	}

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", data)
	if err != nil {
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	parsedResponse := &recaptchaValidationResponse{}
	err = json.Unmarshal(body, parsedResponse)
	if err != nil {
		return false
	}

	return parsedResponse.Success
}
