package webstub

import (
  . "github.com/smartystreets/goconvey/convey"
  "github.com/parnurzeal/gorequest"
  "testing"
  "net/http"
)

func TestWebstub(t *testing.T) {
  server := NewStubServer()
  request := gorequest.New()

  Convey("Stubbing HTTP requests", t, func() {
    Convey("GET", func() {
      requestSpec := &RequestSpec{
        Status: 200,
        Response: []byte("test"),
      }

      server.GET("/something", requestSpec)

      Convey("Correct request made", func() {
        resp, body, _ := request.Get(server.Server.URL + "/something?a=b&c=d").End()

        So(resp.StatusCode, ShouldEqual, 200)
        So(body, ShouldResemble, string(requestSpec.Response))
      })

      Convey("Query string is wrong", func() {
        resp, body, _ := request.Get(server.Server.URL + "/something").End()

        So(resp.StatusCode, ShouldEqual, 200)
        So(body, ShouldResemble, string(requestSpec.Response))

        expectedResult := &RequestResult{
          Requested: true,
          Body: []byte{},
        }

        So(expectedResult, ShouldResemble, requestSpec.Result)
      })

      Convey("Path is wrong", func() {
        resp, _, _:= request.Get(server.Server.URL + "/somethingElse").End()

        So(resp.StatusCode, ShouldEqual, 500)
        So(requestSpec.Result.Requested, ShouldEqual, false)
      })
    })

    Convey("POST", func() {
      requestSpec := &RequestSpec{
        Status: 201,
        Response: []byte("test"),
      }

      server.POST("/example", requestSpec)

      resp, body,  _ := request.Post(server.Server.URL + "/example").
        Set("Content-Type", "text/plain").
        Send("example=value").
        End()

      So(resp.StatusCode, ShouldEqual, 201)
      So(requestSpec.Result.Requested, ShouldEqual, true)
      So(requestSpec.Result.Body, ShouldResemble, []byte("example=value"))
      So(body, ShouldEqual, "test")
    })

    Convey("PUT", func() {
      requestSpec := &RequestSpec{
        Status: 200,
        Response: []byte("test"),
      }

      server.PUT("/example", requestSpec)

      resp, body,  _ := request.Put(server.Server.URL + "/example").
        SendString("example=value").
        End()

      So(resp.StatusCode, ShouldEqual, 200)
      So(requestSpec.Result.Requested, ShouldEqual, true)
      So(requestSpec.Result.Body, ShouldResemble, []byte("example=value"))
      So(body, ShouldEqual, "test")
    })

    Convey("HEAD", func() {
      requestSpec := &RequestSpec{
        Status: 200,
        Response: []byte("test"),
        ResponseHeaders: http.Header{
          "Key": []string{"Value", "List"},
        },
      }

      server.HEAD("/example", requestSpec)

      resp, body, _ := request.Head(server.Server.URL + "/example?a=b").End()
      So(resp.StatusCode, ShouldEqual, 200)
      So(body, ShouldEqual, "")
      So(requestSpec.Result.Query, ShouldEqual, "a=b")
      So(resp.Header["Key"], ShouldResemble, []string{"Value", "List"})
    })

    Convey("DELETE", func() {
      requestSpec := &RequestSpec{
        Status: 200,
        Response: []byte("test"),
      }

      server.DELETE("/example", requestSpec)

      resp, body, _ := request.Delete(server.Server.URL + "/example?a=b").End()
      So(resp.StatusCode, ShouldEqual, 200)
      So(body, ShouldEqual, "test")
      So(requestSpec.Result.Query, ShouldEqual, "a=b")
    })

    Convey("PATCH", func() {
      requestSpec := &RequestSpec{
        Status: 200,
        Response: []byte("test"),
      }

      server.PATCH("/example", requestSpec)

      resp, body, _ := request.Patch(server.Server.URL + "/example?a=b").End()
      So(resp.StatusCode, ShouldEqual, 200)
      So(body, ShouldEqual, "test")
      So(requestSpec.Result.Query, ShouldEqual, "a=b")
    })
  })

  server.Server.Close()
}
