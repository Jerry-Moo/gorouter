package gorouter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var errorFormat, expected string

func init() {
	expected = "hi, gorouter"
	errorFormat = "handler returned unexpected body: got %v want %v"
}

func TestRouter_GET(t *testing.T) {
	router := New()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/hi", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.GET("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// 测试 URL 后缀
func TestRouter_URL_SUFFIX(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/hello/", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// Test POST
func TestRouter_POST(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPost, "/hi", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.POST("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// Test PATCH
func TestRouter_PATCH(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPatch, "/hi", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.PATCH("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// Test DELETE
func TestRouter_DELETE(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodDelete, "/hi", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.DELETE("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// Test PUT
func TestRouter_PUT(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPut, "/hi", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.PUT("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// Test Group
func TestRouter_Group(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	prefix := "/api"

	req, err := http.NewRequest(http.MethodGet, prefix+"/hi", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.Group(prefix).GET("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// Test CustomHandleNotFound
func TestRouter_CustomHandleNotFound(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/xxx", nil)
	if err != nil {
		t.Fatal(err)
	}

	customNotFoundStr := "404 page !"
	router.NotFoundFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, customNotFoundStr)
	})

	router.GET("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != customNotFoundStr {
		t.Errorf(errorFormat, rr.Body.String(), customNotFoundStr)
	}
}

// Test HandleNotFound
func TestRouter_HandleNotFound(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/xxx", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.GET("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String()[:3] != "404" {
		t.Errorf(errorFormat, rr.Body.String(), "404 page not found\n")
	}
}


// Test CustomPanicHandler
func TestRouter_CustomPanicHandler(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req ,err := http.NewRequest(http.MethodPost, "/xxx", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		t.Log("received a panic", err)
	}

	router.POST("/xxx", func(w http.ResponseWriter, r *http.Request) {
		panic("err")
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)
}

// Test NotFoundMethod
func TestRouter_NotFoundMethod(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPost, "/aaa", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.GET("/aaa", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)
}


// Test GetParam
func TestGetParam(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	param := "1"
	req, err := http.NewRequest(http.MethodGet, "/test/" + param, nil)
	if err != nil {
		t.Fatal(err)
	}

	router.GET("/test/:id", func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r, "id")
		if id != param {
			t.Fatal("TestGetParam test fail")
		}
	})
	router.ServeHTTP(rr, req)
}