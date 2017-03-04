/**************************************************************
 * Copyright (c) 2016 Christopher Eaton
 * https://github.com/chriseaton
 * This source code is subject to the terms of the MIT License.
 *************************************************************/
package recaptcha

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

//The REST URL to Google's ReCaptcha verification call.
var recaptchaURL = "https://www.google.com/recaptcha/api/siteverify"

//Sent to the ReCaptcha service
type Challenge struct {
	//The shared key between your site and ReCAPTCHA.
	Secret string
	//The user's response token provided by the reCAPTCHA. This is automatically populated from the form
	//value for key "g-recaptcha-response"
	FormValue string
	//*Optional* Remote IP address of the user. This is automatically populated from the HTTP request if left blank,
	//the IP address can be omitted entirely by setting the value to "none"
	RemoteIP string
	//*Optional* Override the http client object with a custom one (see examples/ for Google AppEngine).
	HttpClient *http.Client
}

//Returned from the ReCaptcha service
type Response struct {
	//True if the ReCaptcha verified ok. False if not.
	Success bool
	//Timestamp of the challenge
	ChallengeTimestamp string `json:"challenge_ts"`
	//Hostname of the site where the reCAPTCHA was solved
	Hostname string
	//Slice of errors returned from the ReCaptcha verification service.
	Errors []string `json:"error-codes"`
}

func getClientIPAddress(r *http.Request) string {
	var ip string
	ipProxy := r.Header.Get("X-FORWARDED-FOR")
	if len(ipProxy) > 0 {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	} else if r.RemoteAddr != "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	//disallow IPs that don't parse correctly
	if net.ParseIP(ip) == nil {
		ip = ""
	}
	return ip
}

//Verify a ReCaptcha challenge by extracting the values from the http.Request.
func VerifyRequest(r *http.Request, secret string) (*Response, error) {
	c := &Challenge{
		Secret:    secret,
		FormValue: r.FormValue("g-recaptcha-response"),
		RemoteIP:  getClientIPAddress(r),
	}
	return Verify(c)
}

//Verify a ReCaptcha challenge.
func Verify(c *Challenge) (*Response, error) {
	if c == nil {
		return nil, fmt.Errorf("The challenge object argument must be specified.")
	} else if c.Secret == "" {
		return nil, fmt.Errorf("Your site's challenge secret must be non-empty.")
	} else if c.FormValue == "" {
		return nil, fmt.Errorf("The ReCaptcha response form value must be non-empty.")
	}
	var client *http.Client
	var rr = &Response{}
	if c.HttpClient != nil {
		//use custom http client if provided
		client = c.HttpClient
	} else {
		//use default http.Client
		client = &http.Client{}
	}
	httpResp, err := client.PostForm(recaptchaURL, url.Values{
		"secret":   {c.Secret},
		"response": {c.FormValue},
		"remoteip": {c.RemoteIP},
	})
	defer httpResp.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(httpResp.Body).Decode(rr)
	if err != nil {
		return nil, err
	} else if rr.Success == false {
		//failed verification, set error message as first error if present
		if rr != nil && len(rr.Errors) > 0 {
			err = fmt.Errorf(rr.Errors[0])
		}
	}
	return rr, err
}
