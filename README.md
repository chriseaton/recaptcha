# [recaptcha](https://gitlab.com/chriseaton/recaptcha)

[recaptcha](https://gitlab.com/chriseaton/recaptcha) is an easy to use Google ReCaptcha library written in Go, 
supporting the latest v2 API. It supports simple verification by a single function call, or can be used with Google
AppEngine by replacing the http client with AppEngine's urlfetch client.

Report issues [here](https://gitlab.com/chriseaton/recaptcha/issues).

### Getting Started

In your console:
````
go get gitlab.com/chriseaton/recaptcha
````
In your go file:
````
import "gitlab.com/chriseaton/recaptcha"
````

### Examples
**Basic Usage**
````
func pageHandler(w http.ResponseWriter, r *http.Request) {
	res, err := recaptcha.VerifyRequest(r, "{MY_SECRET_HERE}")
	if err != nil {
		fmt.Fprintf(w, "Verify failed: %+v\n", err)
	} else {
		fmt.Fprintf(w, "Verified ok: %+v\n", res)
	}
}
````

**Google AppEngine**
````
func pageHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ctx := appengine.NewContext(r)
	c := &recaptcha.Challenge{
		Secret: "{MY_SECRET_HERE}",
		FormValue: r.FormValue("g-recaptcha-response")
		HttpClient: urlfetch.Client(ctx),
	}
	res, err := recaptcha.Verify(c)
	if err != nil {
		fmt.Fprintf(w, "Verify failed: %+v\n", err)
	} else {
		fmt.Fprintf(w, "Verified ok: %+v\n", res)
	}
}
````