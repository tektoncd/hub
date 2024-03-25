package httpcheck

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestingT is an interface wrapper around *testing.T.
type TestingT interface {
	Errorf(format string, args ...any)
	FailNow()
}

// Tester represents the HTTP tester having TestingT.
type Tester struct {
	t TestingT
	*Checker
}

// T returns the TestingT.
func (tt *Tester) T() TestingT {
	return tt.t
}

// Check - Will make request to built request object.
// After request is made, it will save response object for future assertions
// Responsibility of this method is also to start and stop HTTP server
func (tt *Tester) Check() *Tester {
	// start server in new goroutine
	tt.run()
	defer tt.stop()

	if tt.debug {
		log.Println("==", tt.request.Method, tt.request.URL)
		log.Println(">> header", tt.request.Header)
		body := "nil"
		if tt.request.Body != nil {
			b, err := io.ReadAll(tt.request.Body)
			require.NoError(tt.t, err, "failed to read request body")
			tt.request.Body = io.NopCloser(bytes.NewReader(b))
			body = string(b)
		}
		log.Println(">> body:", string(body))
	}

	newJar, _ := cookiejar.New(nil)
	for name := range tt.pcookies {
		for _, oldCookie := range tt.client.Jar.Cookies(tt.request.URL) {
			if name == oldCookie.Name {
				newJar.SetCookies(tt.request.URL, []*http.Cookie{oldCookie})
				break
			}
		}
	}

	tt.client.Jar = newJar
	response, err := tt.client.Do(tt.request)
	if err != nil {
		assert.FailNow(tt.t, err.Error())
	}
	defer func() {
		_ = response.Body.Close()
	}()

	// save response for assertion checks
	b, err := io.ReadAll(response.Body)
	require.NoError(tt.t, err)

	tt.response = response
	tt.response.Body = io.NopCloser(bytes.NewReader(b))

	if tt.debug {
		log.Println("<< status:", tt.response.Status)
		log.Println("<< body:", string(b))
	}

	return tt
}

// Cb - Will call provided callback function with current response
func (tt *Tester) Cb(callback func(*http.Response)) {
	callback(tt.response)
}

// Will run HTTP server
func (tt *Tester) run() {
	// log.Println("running server")
	if tt.external {
		return
	}
	tt.server.Start()
}

// Will stop HTTP server
func (tt *Tester) stop() {
	// log.Println("stopping server")
	if tt.external {
		return
	}
	tt.server.Close()
	tt.server = createServer(tt.handler)
}
