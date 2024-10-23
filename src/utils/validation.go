package utils

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"net/url"
)

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	Score       float64  `json:"score"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}

func ValidateQueryParams(context echo.Context, fields []string) bool {
	for _, field := range fields {
		if context.QueryParam(field) == "" {
			return false
		}
	}

	return true
}

func VerifyRecaptcha(token string) (bool, error) {
	verifyURL := "https://www.google.com/recaptcha/api/siteverify"
	data := url.Values{
		"secret":   {""},
		"response": {token},
	}

	resp, err := http.PostForm(verifyURL, data)
	if err != nil {
		return false, fmt.Errorf("failed to verify reCAPTCHA: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read reCAPTCHA response: %w", err)
	}

	var recaptchaResponse RecaptchaResponse
	err = json.Unmarshal(body, &recaptchaResponse)
	if err != nil {
		return false, fmt.Errorf("failed to parse reCAPTCHA response: %w", err)
	}

	print(recaptchaResponse.Score)

	// Reject if score is less than 0.5
	if recaptchaResponse.Score < 0.5 {
		return false, nil
	}

	print(recaptchaResponse.Success)

	return recaptchaResponse.Success, nil
}
