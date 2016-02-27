/**************************************************************
 * Copyright (c) 2016 Christopher Eaton
 * https://gitlab.com/chriseaton
 * This source code is subject to the terms of the MIT License.
 *************************************************************/
package recaptcha

import (
	"fmt"
	"github.com/joho/godotenv"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

//To test, it's recommended to create a .env file using Google's automated testing secret value, which is
//available here: https://developers.google.com/recaptcha/docs/faq
//
//Sample .env file:
//GOOGLE_RECAPTCHA_TEST_KEY={ReCaptcha Site Key}
//GOOGLE_RECAPTCHA_TEST_SECRET={ReCaptcha Secret}
//GOOGLE_RECAPTCHA_TEST_RESPONSE={Any valid or invalid response value, depending on what you want to test}

//Returns the key specified as an environmental variable (or loaded from a .env file, see above), or if not found
//uses Google's public key for automated testing.
func getTestKey() string {
	key := os.Getenv("GOOGLE_RECAPTCHA_TEST_KEY")
	if key == "" {
		key = "6LeIxAcTAAAAAJcZVRqyHh71UMIEGNQ_MXjiZKhI" //Google's public key for automated testing.
	}
	return key
}

//Returns the secret specified as an environmental variable (or loaded from a .env file, see above), or if not found
//uses Google's public secret for automated testing.
func getTestSecret() string {
	secret := os.Getenv("GOOGLE_RECAPTCHA_TEST_SECRET")
	if secret == "" {
		secret = "6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe" //Google's public secret for automated testing.
	}
	return secret
}

//Test the http call using Google's testing secret
func TestVerify(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Skip("No .env file for loading test variables. Skipping...")
	}
	c := &Challenge{
		Secret:    getTestSecret(),
		FormValue: os.Getenv("GOOGLE_RECAPTCHA_TEST_RESPONSE"),
	}
	t.Logf("Sending challenge: %+v", c)
	res, err := Verify(c)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", res)
	}
}

func TestVerifyRequest(t *testing.T) {
	s := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Errorf("Failed to start web server for testing, error: %s", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
				<head><script src='https://www.google.com/recaptcha/api.js'></script></head>
				<body>
					<form action="/test" method="post">
						<div class="g-recaptcha" data-sitekey="%s"></div>
						<input type="submit" value="Send Test">
					</form>
				</body>
			</html>`, getTestKey())
	})
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		res, err := VerifyRequest(r, getTestSecret())
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("%+v", res)
		}
		//close the webserver
		l.Close()
	})
	s.Serve(l)
}
