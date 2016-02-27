/**************************************************************
 * Copyright (c) 2016 Christopher Eaton
 * https://gitlab.com/chriseaton
 * This source code is subject to the terms of the MIT License.
 *************************************************************/
package recaptcha

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
)

//These vars are loaded from the recaptcha_test.data file. Populated using Google's automated testing values
//which are available here: https://developers.google.com/recaptcha/docs/faq
var (
	testSiteKey  string
	testSecret   string
	testResponse string
)

//load the test data
func TestMain(m *testing.M) {
	file, err := os.Open("recaptcha_test.data")
	if err != nil {
		os.Exit(5)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var index int
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] != '#' {
			switch index {
			case 0:
				testSiteKey = line
			case 1:
				testSecret = line
			case 2:
				testResponse = line
			}
			index++
			if index > 2 {
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		os.Exit(10)
	}
	os.Exit(m.Run())
}

//Test the http call using Google's testing secret
func TestVerify(t *testing.T) {
	c := &Challenge{
		Secret:    testSecret,
		FormValue: testResponse,
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
	s := &http.Server{}
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
			</html>`, testSiteKey)
	})
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		res, err := VerifyRequest(r, testSecret)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("%+v", res)
		}
		//close the webserver
		l.Close()
	})
	t.Logf("Started test webserver.")
	s.Serve(l)
}
