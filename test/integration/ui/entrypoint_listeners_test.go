package ui

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_client"
	"github.com/supergiant/supergiant/test/fake_core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEntrypointListenersList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI EntrypointListeners List works correctly", t, func() {

		table := []struct {
			// Mocks
			// mockClientListItems []*model.EntrypointListener
			// mockClientListError error
			mockAuthenticated bool
			// Expectations
			responseStatusCode int
			responseURL        string
		}{
			// A successful example
			{
				mockAuthenticated:  true,
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/entrypoint_listeners",
			},
			// Unauthenticated
			{
				mockAuthenticated:  false,
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/sessions/new",
			},
		}

		for _, item := range table {

			// For unauthenticated Session-based routes
			srv.Core.APIClient = func(authType string, authToken string) *client.Client {
				return new(client.Client)
			}

			srv.Core.Sessions = &fake_core.Sessions{
				ClientFn: func(sessionID string) *client.Client {
					if item.mockAuthenticated {
						return &client.Client{}
					}
					return nil
				},
			}

			req, _ := http.NewRequest("GET", "http://localhost:10000/ui/entrypoint_listeners", nil)

			// As long as we have a cookie with the right name, it will trigger the
			// use of our fake_core.Sessions above.
			cookie := &http.Cookie{
				Name:  core.SessionCookieName,
				Value: "fake-session-id",
				Path:  "/",
			}
			req.AddCookie(cookie)

			resp, _ := http.DefaultClient.Do(req)

			So(resp.StatusCode, ShouldEqual, item.responseStatusCode)
			So(resp.Request.URL.String(), ShouldEqual, item.responseURL)
		}
	})
}

//------------------------------------------------------------------------------

func TestEntrypointListenersNew(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI EntrypointListeners New works correctly", t, func() {

		table := []struct {
			// Mocks
			mockAuthenticated bool
			// Expectations
			responseStatusCode int
			responseURL        string
		}{
			// A successful example
			{
				mockAuthenticated:  true,
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/entrypoint_listeners/new",
			},
			// Unauthenticated
			{
				mockAuthenticated:  false,
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/sessions/new",
			},
		}

		for _, item := range table {

			// For unauthenticated Session-based routes
			srv.Core.APIClient = func(authType string, authToken string) *client.Client {
				return new(client.Client)
			}

			srv.Core.Sessions = &fake_core.Sessions{
				ClientFn: func(sessionID string) *client.Client {
					if item.mockAuthenticated {
						return &client.Client{}
					}
					return nil
				},
			}

			req, _ := http.NewRequest("GET", "http://localhost:10000/ui/entrypoint_listeners/new", nil)

			// As long as we have a cookie with the right name, it will trigger the
			// use of our fake_core.Sessions above.
			cookie := &http.Cookie{
				Name:  core.SessionCookieName,
				Value: "fake-session-id",
				Path:  "/",
			}
			req.AddCookie(cookie)

			resp, _ := http.DefaultClient.Do(req)

			So(resp.StatusCode, ShouldEqual, item.responseStatusCode)
			So(resp.Request.URL.String(), ShouldEqual, item.responseURL)
		}
	})
}

//------------------------------------------------------------------------------

func TestEntrypointListenersCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI EntrypointListeners Create works correctly", t, func() {

		table := []struct {
			// Input
			jsonInput string
			// Mocks
			mockAuthenticated     bool
			mockClientCreateError error
			// Expectations
			redirected          bool // how to distinguish between a successful create and a failure
			modelPassedToCreate *model.EntrypointListener
			responseStatusCode  int
			responseURL         string
		}{
			// A successful example
			{
				jsonInput: `{
          "name": "test",
          "entrypoint_name": "test",
          "entrypoint_port": 80,
          "node_port": 8080
        }`,
				mockAuthenticated:     true,
				mockClientCreateError: nil,
				redirected:            true,
				modelPassedToCreate: &model.EntrypointListener{
					Name:           "test",
					EntrypointName: "test",
					EntrypointPort: 80,
					NodePort:       8080,
				},
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/entrypoint_listeners",
			},
			// Bad JSON input
			{
				jsonInput: `{
          "name": "te..
        }`,
				mockAuthenticated:     true,
				mockClientCreateError: nil,
				modelPassedToCreate:   nil,
				responseStatusCode:    200,
				responseURL:           "http://localhost:10000/ui/entrypoint_listeners",
			},

			// Unauthenticated
			{
				jsonInput: `{
          "name": "test"
        }`,
				mockAuthenticated:     false,
				mockClientCreateError: nil,
				redirected:            true, // to login page
				modelPassedToCreate:   nil,
				responseStatusCode:    200,
				responseURL:           "http://localhost:10000/ui/sessions/new",
			},
			// Client Create error
			{
				jsonInput: `{
          "name": "test"
        }`,
				mockAuthenticated:     true,
				mockClientCreateError: errors.New("something bad"),
				modelPassedToCreate: &model.EntrypointListener{
					Name: "test",
				},
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/entrypoint_listeners",
			},
		}

		for _, item := range table {

			var modelPassedToCreate *model.EntrypointListener
			var redirected bool

			// For unauthenticated Session-based routes
			srv.Core.APIClient = func(authType string, authToken string) *client.Client {
				return new(client.Client)
			}

			srv.Core.Sessions = &fake_core.Sessions{
				ClientFn: func(sessionID string) *client.Client {
					if item.mockAuthenticated {
						return &client.Client{
							EntrypointListeners: &fake_client.EntrypointListeners{
								Collection: fake_client.Collection{
									CreateFn: func(m model.Model) error {
										modelPassedToCreate = m.(*model.EntrypointListener)
										return item.mockClientCreateError
									},
								},
							},
						}
					}
					return nil
				},
			}

			form := url.Values{}
			form.Add("json_input", item.jsonInput)
			body := strings.NewReader(form.Encode())

			req, _ := http.NewRequest("POST", "http://localhost:10000/ui/entrypoint_listeners", body)

			// As long as we have a cookie with the right name, it will trigger the
			// use of our fake_core.Sessions above.
			cookie := &http.Cookie{
				Name:  core.SessionCookieName,
				Value: "fake-session-id",
				Path:  "/",
			}
			req.AddCookie(cookie)

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			client := http.DefaultClient

			// https://github.com/golang/go/issues/4800
			client.CheckRedirect = func(r *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return errors.New("too many redirects")
				}
				if len(via) == 0 {
					return nil
				}

				redirected = true

				for attr, val := range via[0].Header {
					if _, ok := r.Header[attr]; !ok {
						r.Header[attr] = val
					}
				}
				return nil
			}

			resp, _ := client.Do(req)

			// defer resp.Body.Close()
			// b, _ := ioutil.ReadAll(resp.Body)
			// fmt.Println(string(b))

			So(resp.StatusCode, ShouldEqual, item.responseStatusCode)
			So(resp.Request.URL.String(), ShouldEqual, item.responseURL)
			So(redirected, ShouldEqual, item.redirected)
			So(modelPassedToCreate, ShouldResemble, item.modelPassedToCreate)
		}
	})
}

//------------------------------------------------------------------------------

func TestEntrypointListenersGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI EntrypointListeners Get works correctly", t, func() {

		table := []struct {
			// Mocks
			mockAuthenticated bool
			mock404           bool
			// Expectations
			responseStatusCode int
			responseURL        string
		}{
			// A successful example
			{
				mockAuthenticated:  true,
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/entrypoint_listeners/1",
			},
			// Unauthenticated
			{
				mockAuthenticated:  false,
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/sessions/new",
			},
			// 404
			{
				mockAuthenticated:  true,
				mock404:            true,
				responseStatusCode: 404,
				responseURL:        "http://localhost:10000/ui/entrypoint_listeners/1",
			},
		}

		for _, item := range table {

			// For unauthenticated Session-based routes
			srv.Core.APIClient = func(authType string, authToken string) *client.Client {
				return new(client.Client)
			}

			srv.Core.Sessions = &fake_core.Sessions{
				ClientFn: func(sessionID string) *client.Client {
					if item.mockAuthenticated {
						return &client.Client{
							EntrypointListeners: &fake_client.EntrypointListeners{
								Collection: fake_client.Collection{
									GetFn: func(id interface{}, m model.Model) error {
										if item.mock404 {
											return errors.New("404")
										}
										return nil
									},
								},
							},
						}
					}
					return nil
				},
			}

			req, _ := http.NewRequest("GET", "http://localhost:10000/ui/entrypoint_listeners/1", nil)

			// As long as we have a cookie with the right name, it will trigger the
			// use of our fake_core.Sessions above.
			cookie := &http.Cookie{
				Name:  core.SessionCookieName,
				Value: "fake-session-id",
				Path:  "/",
			}
			req.AddCookie(cookie)

			resp, _ := http.DefaultClient.Do(req)

			So(resp.StatusCode, ShouldEqual, item.responseStatusCode)
			So(resp.Request.URL.String(), ShouldEqual, item.responseURL)
		}
	})
}
