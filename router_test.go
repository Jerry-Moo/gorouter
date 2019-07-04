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

	req, err := http.NewRequest(http.MethodPost, "/xxx", nil)
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
	req, err := http.NewRequest(http.MethodGet, "/test/"+param, nil)
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

// Test GetAllParams
func TestGetAllParams(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	param1 := "1"
	param2 := "2"

	req, err := http.NewRequest(http.MethodGet, "/param1/"+param1+"/param2/"+param2, nil)
	if err != nil {
		t.Fatal(err)
	}

	router.GET("/param1/:param1/param2/:param2", func(w http.ResponseWriter, r *http.Request) {
		params := GetAllParams(r)
		if params["param1"] != param1 {
			t.Fatal("TestGetAllParams test fail")
		}
		if params["param2"] != param2 {
			t.Fatal("TestGetAllParams test fail")
		}
	})
	router.ServeHTTP(rr, req)
}

// Test ParamsMiss
func TestGetAllParamsMiss(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/param1", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.GET("/param1", func(w http.ResponseWriter, r *http.Request) {
		params := GetAllParams(r)

		if params != nil {
			t.Fatal("TestGetAllParams test fail")
		}
	})
	router.ServeHTTP(rr, req)
}

// logHandlfunc
func withLoggint(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("Logged connection from %s,", r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}

// Test Use
func TestRouter_Use(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/hi", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.Use(withLoggint)
	router.GET("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)
	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// Test UserFroRoot
func TestRouter_UseForRoot(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.Use(withLoggint)
	expected := "hi index"
	router.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)
	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// Test Regex
func TestRouter_Regex(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/param/1", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.GET("/param/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
		id := GetParam(r, "id")
		if id != "1" {
			t.Fatal("TestGetAllParams test fail")
		}
	})

	router.ServeHTTP(rr, req)
	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// Test HandleRoot
func TestRouter_HandleRoot(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := "hi,root"
	router.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// Test HandlePanic
func TestRouter_HandlePanic(t *testing.T) {
	router := New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered in f", r)
		}
	}()

	router.Handle("", "/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})

	router.ServeHTTP(rr, req)
	if rr.Body.String() != expected {
		t.Errorf(errorFormat, rr.Body.String(), expected)
	}
}

// Test Match
func TestRouter_Match(t *testing.T) {
	router := New()
	requestUrl := "/xxx/1/yyy/2"

	ok := router.Match(requestUrl, "/xxx/:param1/yyy/:param2")
	if !ok {
		t.Fatal("TestRouter_match test fail")
	}

	errorRequestUrl := "#xxx#1#yyy#2"
	ok = router.Match(errorRequestUrl, "/xxx/:param1/yyy/:param2")
	if ok {
		t.Fatal("TestRouter_Match test fail")
	}

	errorPath := "#xxx#1#yyy#2"
	ok = router.Match(errorRequestUrl, errorPath)
	if ok {
		t.Fatal("TestRouter_Match test fail")
	}

	missRequestUrl := "/xxx/1/yyy/###"
	ok = router.Match(missRequestUrl, "/xxx/:param1/yyy/:param2")

	if ok {
		t.Fatal("TestRouter_match test fail")
	}
}

// Test Generate
func TestRouter_Generate(t *testing.T) {
	mux := New()

	routeName1 := "user_event"
	params := make(map[string]string)
	params["user"] = "jerrywu"

	// GETAndName
	mux.GETAndName("/users/:user/events", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/:user/events"))
	}, routeName1)

	if url, _ := mux.Generate(http.MethodGet, routeName1, params); url != "/users/jerrywu/events" {
		t.Fatal("TestRouter_Generate test fail")
	}

	routeName2 := "user_repos_keys"
	params = make(map[string]string)
	params["owner"] = "jerrywu"
	params["repo"] = "jerrywu_repo"

	// POSTAndName
	mux.POSTAndName("/repos/{owner:\\w+}/{repo:\\w+}/keys", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/:user/repos"))
	}, routeName2)

	if url, _ := mux.Generate(http.MethodPost, routeName2, params); url != "/repos/jerrywu/jerrywu_repo/keys" {
		t.Fatal("TestRouter_Generate test fail")
	}
	//DELETEAndName
	routeName3 := "repos_releases"
	mux.DELETEAndName("/repos/{owner:\\w+}/{repo:\\w+}/releases/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/repos/{owner:\\w+}/{repo:\\w+}/releases/{id:[0-9]+}"))
	}, routeName3)
	params = make(map[string]string)
	params["owner"] = "xujiajun"
	params["repo"] = "xujiajun_repo"
	params["id"] = "100"
	if url, _ := mux.Generate(http.MethodDelete, routeName3, params); url != "/repos/xujiajun/xujiajun_repo/releases/100" {
		t.Fatal("TestRouter_Generate test fail")
	}

	//PUTAndName
	routeName4 := "user_following"
	params = make(map[string]string)
	params["user"] = "xujiajun001"
	mux.PUTAndName("/user/following/{user:\\w+}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/user/following/{user:\\w+}"))
	}, routeName4)

	if url, _ := mux.Generate(http.MethodPut, routeName4, params); url != "/user/following/xujiajun001" {
		t.Fatal("TestRouter_Generate test fail")
	}

	//PATCHAndName
	routeName6 := "repos_keys"
	params = make(map[string]string)
	params["owner"] = "xujiajun001"
	params["repo"] = "xujiajun_repo"
	params["id"] = "100"
	mux.PATCHAndName("/repos/:owner/:repo/keys/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/repos/:owner/:repo/keys/{id:[0-9]+}"))
	}, routeName6)

	if url, _ := mux.Generate(http.MethodPatch, routeName6, params); url != "/repos/xujiajun001/xujiajun_repo/keys/100" {
		t.Fatal("TestRouter_Generate test fail")
	}

	//params contains wrong parameters
	routeName5 := "user_event3"
	mux.GETAndName("/users/{user:\\w+}/events", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/{user:\\w+}/events"))
	}, routeName5)
	params = make(map[string]string)
	params["user"] = "@@@@"
	if _, err := mux.Generate(http.MethodGet, routeName5, params); err == nil {
		t.Fatal("TestRouter_Generate test fail")
	}
	mux.GETAndName("/users/:user/events", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/:user/events"))
	}, routeName5)
	params = make(map[string]string)
	params["user"] = "@@@@"
	if _, err := mux.Generate(http.MethodGet, routeName5, params); err == nil {
		t.Fatal("TestRouter_Generate test fail")
	}

	//pattern grammar error
	routeName7 := "user_event4"
	mux.GETAndName("/users/user:\\w+}/events", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/user:\\w+}/events"))
	}, routeName7)
	params = make(map[string]string)
	params["user"] = "xujiajun"
	if _, err := mux.Generate(http.MethodGet, routeName7, params); err == nil {
		t.Fatal("TestRouter_Generate test fail")
	}

	mux.GETAndName("/users/{user:\\w+/events", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/{user:\\w+/events"))
	}, routeName7)
	params = make(map[string]string)
	params["user"] = "xujiajun"
	if _, err := mux.Generate(http.MethodGet, routeName7, params); err == nil {
		t.Fatal("TestRouter_Generate test fail")
	}

	//cannot found route in tree
	if _, err := mux.Generate("GET", "notFoundRouteName", params); err == nil {
		t.Fatal("TestRouter_Generate test fail")
	}

	//cannot found method in tree
	if _, err := mux.Generate("METHOD", routeName5, params); err == nil {
		t.Fatal("TestRouter_Generate test fail")
	}
}
