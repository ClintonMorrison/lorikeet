package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// See https://developers.google.com/recaptcha/docs/verify

type RecaptchaClient struct {
	debugLogger *log.Logger
	secret      string
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

	rc.debugLogger.Println("validating recaptcha using secret: " + rc.secret)

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", data)
	if err != nil {
		rc.debugLogger.Println("error from recaptcha post request: " + err.Error())
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		rc.debugLogger.Println("error reading recaptcha body: " + err.Error())
		return false
	}

	parsedResponse := &recaptchaValidationResponse{}
	err = json.Unmarshal(body, parsedResponse)
	if err != nil {
		rc.debugLogger.Println("error unmarshalling recaptcha body: " + err.Error())
		return false
	}

	return parsedResponse.Success
}
