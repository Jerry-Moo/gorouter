package gorouter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	ErrGenerateParameters = errors.New("params contains wrong parameters")

	ErrNotFoundRouter = errors.New("can't find route in tree")

	ErrNotFoundMethod = errors.New("can't find method in tree")

	ErrPatternGrammar = errors.New("pattern grammar error")

	defaultPattern = `[\w]+`
	idPattern      = `[\d]+`
	idKey          = `id`

	methods = map[string]struct{}{
		http.MethodGet:    {},
		http.MethodPost:   {},
		http.MethodPut:    {},
		http.MethodDelete: {},
		http.MethodPatch:  {},
	}
)

type (
	// 定义中间件
	MiddlewareType func(next http.HandlerFunc) http.HandlerFunc

	// 自定义路由 多路复用解析http请求
	// 记录所有 URL 参数 并执行路由到函数的转发
	Router struct {
		// Router 的前缀
		prefix string
		// 中间件列表
		middleware []MiddlewareType
		// 树结构
		trees      map[string]*Tree
		parameters Parameters
		// Custom route not found handler
		notFound http.HandlerFunc
		// PanicHandler for handling panic. 恐慌路由
		PanicHandler func(w http.ResponseWriter, r *http.Request, err interface{})
	}
	// 参数记录 - 记录参数
	Parameters struct {
		routeName string
	}
)

// New returns a newly initialized Router object that implements the Router
func New() *Router {
	return &Router{
		trees: make(map[string]*Tree),
	}
}

// GET adds the route `path` that matches a GET http method to
// execute the `handle` http.HandlerFunc.
func (r *Router) GET(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handle)
}

// POST adds the route `path` that matches a POST http method to
// execute the `handle` http.HandlerFunc.
func (r *Router) POST(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handle)
}

// DELETE adds the route `path` that matches a DELETE http method to
// execute the `handle` http.HandlerFunc.
func (r *Router) DELETE(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handle)
}

// PUT adds the route `path` that matches a DELETE http method to
// execute the `handle` http.HandlerFunc.
func (r *Router) PUT(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handle)
}

// PATCH adds the route `path` that matches a DELETE http method to
// execute the `handle` http.HandlerFunc.
func (r *Router) PATCH(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodPatch, path, handle)
}

// GETAndName is short for `GET` and Named routeName
func (r *Router) GETAndName(path string, handle http.HandlerFunc, routeName string) {
	r.parameters.routeName = routeName
	r.GET(path, handle)
}

// POSTAndName is short for `Post` and Named routeName
func (r *Router) POSTAndName(path string, handle http.HandlerFunc, routeName string) {
	r.parameters.routeName = routeName
	r.POST(path, handle)
}

// DELETEAndName is short for `DELETE` and Named routeName
func (r *Router) DELETEAndName(path string, handle http.HandlerFunc, routeName string) {
	r.parameters.routeName = routeName
	r.DELETE(path, handle)
}

// PUTAndName is short for `PUT` and Named routeName
func (r *Router) PUTAndName(path string, handle http.HandlerFunc, routeName string) {
	r.parameters.routeName = routeName
	r.PUT(path, handle)
}

// PATCHAndName is short for `PUT` and Named routeName
func (r *Router) PATCHAndName(path string, handle http.HandlerFunc, routeName string) {
	r.parameters.routeName = routeName
	r.PATCH(path, handle)
}

// Group define routes groups if there is a path prefix that uses `prefix`
func (r *Router) Group(prefix string) *Router {
	return &Router{
		prefix:     prefix,
		trees:      r.trees,
		middleware: r.middleware,
	}
}

// Generate returns reverse routing by method, routeName and params
// 通过method，routeName和params生成返回反向路由
func (r *Router) Generate(method string, routeName string, params map[string]string) (string, error) {
	tree, ok := r.trees[method]
	if !ok {
		return "", ErrNotFoundMethod
	}
	route, ok := tree.routes[routeName]
	if !ok {
		return "", ErrNotFoundRouter
	}

	var segments []string
	res := splitPattern(route.path)
	for _, segment := range res {
		if string(segment[0]) == ":" {
			key := params[string(segment[1:])]
			re := regexp.MustCompile(defaultPattern)
			if one := re.Find([]byte(key)); one == nil {
				return "", ErrGenerateParameters
			}
			segments = append(segments, key)
			continue
		}

		if string(segment[0]) == "{" {
			segmentLen := len(segment)
			if string(segment[segmentLen-1]) == "}" {
				splitRes := strings.Split(string(segment[1:segmentLen-1]), ":")
				re := regexp.MustCompile(splitRes[1])
				key := params[splitRes[0]]
				if one := re.Find([]byte(key)); one == nil {
					return "", ErrGenerateParameters
				}
				segments = append(segments, key)
				continue
			}
			return "", ErrPatternGrammar
		}
		if string(segment[len(segment)-1]) == "}" && string(segment[0]) != "{" {
			return "", ErrPatternGrammar
		}
		segments = append(segments, segment)
		continue
	}
	return "/" + strings.Join(segments, "/"), nil
}

// NotFoundFunc registers a handler when the request route is not found
func (r *Router) NotFoundFunc(handler http.HandlerFunc) {
	r.notFound = handler
}

// Handle register a new request handler with the given path and method.
func (r *Router) Handle(method string, path string, handle http.HandlerFunc) {
	if _, ok := methods[method]; !ok {
		panic(fmt.Errorf("invalid method"))
	}

	// 新增路由的时候 以请求方式获取 树结构
	tree, ok := r.trees[method]
	if !ok {
		tree = NewTree()
		r.trees[method] = tree
	}
	// 判断前缀是否为空 不为空把前缀添加到 新路由前缀
	if r.prefix != "" {
		path = r.prefix + "/" + path
	}

	if routeName := r.parameters.routeName; routeName != "" {
		tree.parameters.routeName = routeName
	}

	tree.Add(path, handle, r.middleware...)
}

// GetParam returns route param stored in http.request.
func GetParam(r *http.Request, key string) string {
	return GetAllParams(r)[key]
}

// contextKeyType is a private struct that is used for storing values in net.Context
// contextKeyType是一个私有结构，用于在net.Context中存储值
type contextKeyType struct{}

// contextKey is the key that is used to store values in net.Context for each request
// contextKey是用于在每个请求的net.Context中存储值的键
var contextKey = contextKeyType{}

// paramsMapType is a private type that is used to store route params
type paramsMapType map[string]string

// GetAllParams returns all route params stored in http.Request.
func GetAllParams(r *http.Request) paramsMapType {
	if values, ok := r.Context().Value(contextKey).(paramsMapType); ok {
		return values
	}
	return nil
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestUrl := req.URL.Path

	// goroutine 异常捕获
	if r.PanicHandler != nil {
		defer func() {
			if err := recover(); err != nil {
				r.PanicHandler(w, req, err)
			}
		}()
	}

	if _, ok := r.trees[req.Method]; !ok {
		r.HandleNotFound(w, req, r.middleware)
		return
	}

	nodes := r.trees[req.Method].Find(requestUrl, false)
	if len(nodes) > 0 {
		node := nodes[0]
		if node.handle != nil {
			if node.path == requestUrl {
				handle(w, req, node.handle, node.middleware)
				return
			}
			if node.path == requestUrl[1:] {
				handle(w, req, node.handle, node.middleware)
				return
			}
		}
	}

	if len(nodes) == 0 {
		res := strings.Split(requestUrl, "/")
		prefix := res[1]
		nodes := r.trees[req.Method].Find(prefix, true)
		for _, node := range nodes {
			if handler := node.handle; handler != nil && node.path != requestUrl {
				if matchParamsMap, ok := r.matchAndParse(requestUrl, node.path); ok {
					ctx := context.WithValue(req.Context(), contextKey, matchParamsMap)
					req = req.WithContext(ctx)
					handle(w, req, handler, node.middleware)
					return
				}
			}
		}
	}
	r.HandleNotFound(w, req, r.middleware)
}

// HandleNotFound registers a handler when the request route is not found
func (r *Router) HandleNotFound(w http.ResponseWriter, req *http.Request, middleware []MiddlewareType) {
	if r.notFound != nil {
		handle(w, req, r.notFound, middleware)
		return
	}
	http.NotFound(w, req)
}

// handle executes middleware chain 执行中间件
func handle(w http.ResponseWriter, req *http.Request, handler http.HandlerFunc, middleware []MiddlewareType) {
	var basehandler = handler
	for _, m := range middleware {
		basehandler = m(basehandler)
	}
	basehandler(w, req)
}

// Match checks if the request matches the route pattern
// 匹配检查请求是否与路由模式匹配
func (r *Router) Match(requestUrl string, path string) bool {
	_, ok := r.matchAndParse(requestUrl, path)
	return ok
}

// matchAndParse checks if the request matches the route path and returns a map of the parsed
// 检查请求是否与路径路径匹配，并返回已解析的映射
func (r *Router) matchAndParse(requestUrl string, path string) (matchParams paramsMapType, b bool) {
	var (
		matchName []string
		pattern   string
	)

	b = true
	matchParams = make(paramsMapType)
	res := strings.Split(path, "/")
	for _, str := range res {
		if str == "" {
			continue
		}

		strLen := len(str)
		firstChar := str[0]
		lastChar := str[strLen-1]
		if string(firstChar) == "{" && string(lastChar) == "}" {
			matchStr := string(str[1 : strLen-1])
			res := strings.Split(matchStr, ":")
			matchName = append(matchName, res[0])
			pattern = pattern + "/" + "(" + res[1] + ")"
		} else if string(firstChar) == ":" {
			matchStr := str
			res := strings.Split(matchStr, ":")
			matchName = append(matchName, res[1])
			if res[1] == idKey {
				pattern = pattern + "/" + "(" + idPattern + ")"
			} else {
				pattern = pattern + "/" + "(" + defaultPattern + ")"
			}
		} else {
			pattern = pattern + "/" + str
		}
	}
	if strings.HasSuffix(requestUrl, "/") {
		pattern = pattern + "/"
	}
	re := regexp.MustCompile(pattern)
	if subMatch := re.FindSubmatch([]byte(requestUrl)); subMatch != nil {
		if string(subMatch[0]) == requestUrl {
			subMatch = subMatch[1:]
			for k, v := range subMatch {
				matchParams[matchName[k]] = string(v)
			}
			return
		}
	}
	return nil, false
}
