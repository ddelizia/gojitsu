package gojitsu

type Route struct {
	Id              string
	RequestMatcher  RequestMatcher
	ResponseBuilder ResponseBuilder
}
