# Webstub
## A Go library for stubbing web services

Webstub is for stubbing out web services in tests. It fulfils a similar
purpose to [webmock](https://github.com/bblimke/webmock).

To get started with Webstub you'll want to import it and in your test
you'll want to create a `*StubServer` with `NewStubServer()`.

    import (
      "github.com/boardiq/webstub"
      "testing"
    )

    func TestXxx(t *testing.T) {
      server := webstub.NewStubServer()
    }

You can then register requests you'd like to stub by creating a `RequestSpec`
and adding this to the server.

    func TestXxx(t *testing.T) {
      server := webstub.NewStubServer()
      requestSpec := &RequestSpec{
        Status: 200,
        Response: []byte("Hello, World!"),
        ResponseHeaders: http.Header{"Key": []string{"Values"}},
      }

      server.GET("/example", requestSpec)
    }

Now when you send a GET request to `/example` you'll get a status code `200`,
a response body `Hello, World!` and included in the response headers will include
the `Key` header with the value `Values`.

The root URL of the server is accessed like so:

    url := server.Server.URL

The `Server` attribute of the `StubServer` is a
[`httptest.Server`](http://golang.org/pkg/net/http/httptest/#Server).

Just like with `GET` we can register endpoints using the various HTTP verbs.

    func TestXxx(t *testing.T) {
      server := webstub.NewStubServer()
      requestSpec := &RequestSpec{
        Status: 201,
        Response: []byte("You created something!"),
      }

      server.POST("/example", requestSpec)
    }

So now if we `POST` to `/example` we will get back a `201` and the
response body `You created something!`

If we want to check that our stub actually received a request, we can
check this from the `RequestResult` that gets attached to the `RequestSpec`
when the server receives a request for that spec.

    func TestXxx(t *testing.T) {
      server := webstub.NewStubServer()
      requestSpec := &RequestSpec{
        Status: 201,
        Response: []byte("You created something!"),
      }

      server.POST("/example", requestSpec)

      // ... The bit where you do your test and hit the endpoint

      // requestSpecResult.Requested is true if the server received a
      // request for that request spec, and false if not.
      if !requestSpec.Result.Requested {
        t.Errorf("Request not received!") 
      }

      // requestSpec.Result.Body contains the request body that was 
      // received by the server.
      if !bytes.Equal(requestSpec.Result.Body, []byte("test body")) {
        t.Errorf("Body was wrong")
      }

      // requestSpec.Result.Query contains the query string that was
      // received by the server.
      if requestSpec.Result.Query != "foo=bar&baz=1" {
        t.Errorf("Query string was wrong")
      }
    }
