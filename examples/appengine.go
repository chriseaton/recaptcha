/**************************************************************
 * Copyright (c) 2016 Christopher Eaton
 * https://gitlab.com/chriseaton
 * This source code is subject to the terms of the MIT License.
 *************************************************************/

package main

import (
	"fmt"
	"gitlab.com/chriseaton/recaptcha"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"net/http"
)

func verify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ctx := appengine.NewContext(r)
	c := &recaptcha.Challenge{
		Secret:     "{MY_SECRET_HERE}",
		FormValue:  r.FormValue("g-recaptcha-response"),
		HttpClient: urlfetch.Client(ctx),
	}
	res, err := recaptcha.Verify(c)
	if err != nil {
		fmt.Fprintf(w, "Verify failed: %+v\n", err)
	} else {
		fmt.Fprintf(w, "Verified ok: %+v\n", res)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
				<head><script src='https://www.google.com/recaptcha/api.js'></script></head>
				<body>
					<form action="/verify" method="post">
						<div class="g-recaptcha" data-sitekey="%s"></div>
						<input type="submit" value="Submit Form">
					</form>
				</body>
			</html>`, "{MY_SITE_KEY}")
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/verify", verify)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Sprintf("Failed to start web server: %s", err)
	}
}
