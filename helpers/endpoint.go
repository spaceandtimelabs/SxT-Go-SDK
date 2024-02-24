package helpers

import "strings"

func GetSqlEndpoint(subpath string) (endpoint string) {
	return getEndpointByType("sql", subpath)
}

func GetAuthenticationEndpoint(subpath string) (endpoint string) {
	return getEndpointByType("auth", subpath)
}

func GetDiscoverEndpoint(subpath string) (endpoint string) {
	return getEndpointByType("discover", subpath)
}

func getEndpointByType(endpointType, subpath string) (endpoint string) {
	apiEndPoint, _ := ReadEndPointDiscovery()
	segments := []string{apiEndPoint, endpointType, subpath}

	return strings.Join(segments, "/")
}
