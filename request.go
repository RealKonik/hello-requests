package request

import (
	"encoding/json"
	"github.com/hunterbdm/hello-requests/compress"
	"github.com/hunterbdm/hello-requests/http"
	"github.com/hunterbdm/hello-requests/http/cookiejar"
	"github.com/hunterbdm/hello-requests/utils"
	"io/ioutil"
	"strings"
	"time"
)

// Features:
//
// - Matching ClientHello fingerprints (done)
// - Matching http/2 fingerprints (done)
// - Matching http/2 header order (done)
// - Custom normal header ordering (done)
// - Trusted certificate checks (done)
// - JSON response parsing (done)
// - JSON body building (done)
// - Custom idle connection timeouts (done)
// - Custom request timeouts (done)
// - Brotli decompression (done)

// utls additions/fixes:
//
// - PreSharedKey extension support added
// - Fixed the same value being used on both GREASE extensions
//   causing "tls: error decoding message"

// TODO
//
// - Add PSK toggle

func Do(opts Options) (*Response, error) {
	return request(opts)
}

func Jar() *cookiejar.Jar {
	jar, _ := cookiejar.New(nil)
	return jar
}

func request(opts Options) (*Response, error) {
	if opts.ClientSettings == nil {
		opts.ClientSettings = &defaultClientSettings
	} else {
		opts.ClientSettings.AddDefaults()
	}

	// Pull http.Client with ClientSettings options
	httpClient := GetHttpClient(opts.ClientSettings)

	// Check for errors in options provided
	parsedUrl, err := opts.Validate()
	if err != nil {
		return nil, err
	}

	// Add cookie header from Jar
	if opts.Jar != nil {
		cookieHeader := ""

		for i, cookie := range opts.Jar.Cookies(parsedUrl) {
			if i > 0 {
				cookieHeader += " "
			}
			cookieHeader += cookie.String() + ";"
		}

		if cookieHeader != "" {
			opts.Headers["Cookie"] = cookieHeader
		}
	}

	// Build http.Request to pass into the http.Client
	req, err := http.NewRequest(opts.Method, opts.URL, strings.NewReader(opts.Body))
	if err != nil {
		return nil, err
	}

	for name, value := range opts.Headers {
		req.Header.Set(name, value)
	}
	// Add HeaderOrder onto request to be used later in the h2_bundle
	req.HeaderOrder = opts.HeaderOrder

	if req.Header.Get("Host") != "" {
		req.Host = req.Header.Get("Host")
	}

	start := time.Now().UnixNano() / int64(time.Millisecond)
	resp, err := httpClient.Do(req)
	end := time.Now().UnixNano() / int64(time.Millisecond)

	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	var body string
	if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		if encoding, ok := resp.Header["Content-Encoding"]; ok {
			body = compress.Decompress(bodyBytes, encoding[0])
		} else {
			body = string(bodyBytes)
		}
	}

	// Add response cookies to jar
	if opts.Jar != nil {
		opts.Jar.SetCookies(parsedUrl, utils.ReadSetCookies(resp.Header))
	}

	var jsonParsed JSON
	if opts.ParseJSONResponse {
		_ = json.Unmarshal([]byte(body), &jsonParsed)
	} else if contentType, ok := resp.Header["Content-Type"]; ok && strings.Contains(contentType[0], "application/json") {
		// Attempt to parse JSON body if response content-type is "application/json"
		_ = json.Unmarshal([]byte(body), &jsonParsed)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
		Json:       jsonParsed,
		Request:    &opts,
		Time:       int(end - start),
	}, nil
}