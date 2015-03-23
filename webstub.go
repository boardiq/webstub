package webstub

import (
  "net/http"
  "net/http/httptest"
  "io/ioutil"
)

type StubServer struct {
  Routes  *map[string]map[string]*RequestSpec
  Server  *httptest.Server
}

type RequestResult struct {
  Requested     bool
  Query         string
  Body          []byte
}

type RequestSpec struct {
  Status            int
  Response          []byte
  ResponseHeaders   http.Header
  Result            *RequestResult
}

func (s *StubServer) GET(path string, requestSpec *RequestSpec) {
  routes := *s.Routes
  requestSpec.Result = &RequestResult{}
  routes["GET"][path] = requestSpec
}

func (s *StubServer) POST(path string, requestSpec *RequestSpec) {
  routes := *s.Routes
  requestSpec.Result = &RequestResult{}
  routes["POST"][path] = requestSpec
}

func (s *StubServer) PUT(path string, requestSpec *RequestSpec) {
  routes := *s.Routes
  requestSpec.Result = &RequestResult{}
  routes["PUT"][path] = requestSpec
}

func (s *StubServer) HEAD(path string, requestSpec *RequestSpec) {
  routes := *s.Routes
  requestSpec.Result = &RequestResult{}
  routes["HEAD"][path] = requestSpec
}

func (s *StubServer) DELETE(path string, requestSpec *RequestSpec) {
  routes := *s.Routes
  requestSpec.Result = &RequestResult{}
  routes["DELETE"][path] = requestSpec
}

func (s *StubServer) PATCH(path string, requestSpec *RequestSpec) {
  routes := *s.Routes
  requestSpec.Result = &RequestResult{}
  routes["PATCH"][path] = requestSpec
}

func (s *StubServer) Reset() {
  routes := *s.Routes
  routes["GET"] = map[string]*RequestSpec{}
  routes["POST"] = map[string]*RequestSpec{}
  routes["PUT"] = map[string]*RequestSpec{}
  routes["HEAD"] = map[string]*RequestSpec{}
  routes["DELETE"] = map[string]*RequestSpec{}
  routes["PATCH"] = map[string]*RequestSpec{}
}

func NewStubServer() *StubServer {
  initialRoutes := &map[string]map[string]*RequestSpec{
    "GET": map[string]*RequestSpec{},
    "POST": map[string]*RequestSpec{},
    "PUT": map[string]*RequestSpec{},
    "HEAD": map[string]*RequestSpec{},
    "DELETE": map[string]*RequestSpec{},
    "PATCH": map[string]*RequestSpec{},
  }
  var handlerFunction http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
    routes := *initialRoutes
    if requestSpec, ok := routes[r.Method][r.URL.Path]; !ok {
      w.WriteHeader(500)
    } else {
      body, _ := ioutil.ReadAll(r.Body)
      requestSpec.Result = &RequestResult{
        Requested: true,
        Query: r.URL.RawQuery,
        Body: body,
      }
      writeResponseHeaders(w, requestSpec)
      w.WriteHeader(requestSpec.Status)
      w.Write(requestSpec.Response)
    }
  }
  server := httptest.NewServer(handlerFunction)
  return &StubServer{initialRoutes, server}
}

func writeResponseHeaders(w http.ResponseWriter, requestSpec *RequestSpec) {
  for key, values := range requestSpec.ResponseHeaders {
    for _, value := range values {
      w.Header().Add(key, value)
    }
  }
}
