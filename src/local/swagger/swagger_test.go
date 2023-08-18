package swagger

import (
	"testing"
)

var swag2 = Swagger_t{
	Swagger:  "2.0",
	Host:     "api.test.lab",
	Upstream: "svc-simple.simple.cluster.local",
	Prefix:   "/simple",
	Paths: map[string]map[string]*Method_t{
		"/api/articles": {"get": {
			Summary:     "get articles",
			OperationId: "getarticlesGET",
			Produces: []string{
				"*/*",
			},
		},
		},
	},
}

func TestGetEndpoint(t *testing.T) {

	endpoint := swag2.GetEndpoint("/api/articles{articlesId}")

	if endpoint.ApiPath != "/api/articles/{articlesId}" {
		t.Error("Wrong api path generated", endpoint.Path, "instead of /api/articles{articlesId}")
	}

}

func TestGetEndpointsFromPathsPath(t *testing.T) {

	paths := []string{
		"/api/articles{articlesId}",
	}

	endpoints := swag2.GetEndpointsFromPaths(paths)

	if endpoints[0].ApiPath != "/api/articles{articlesId}" {
		t.Error("Wrong api path generated", endpoints[0].ApiPath, "instead of/api/articles{articlesId}")
	}
}

func TestGetEndpointsFromPathsApiPath(t *testing.T) {

	paths := []string{
		"/api/articles{articlesId}",
	}

	endpoints := swag2.GetEndpointsFromPaths(paths)

	if endpoints[0].Path != "/api/articles/(?articlesId>([a-zA-Z0-9])+)$" {
		t.Error("Wrong conversion from swagger path definition to nginx location")
	}

}
