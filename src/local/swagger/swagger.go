package swagger

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Endpoint_t struct {
	Id        string
	Path      string
	ApiPath   string
	Prefix    string
	Upstream  string
	IsRegex   bool
	IsSecured bool
	Limits    []string
	Methods   map[string]*Method_t
}

type Swagger_t struct {
	Swagger  string `json:"swagger"`
	Host     string `json:"host"`
	Upstream string
	Prefix   string
	Paths    map[string]map[string]*Method_t `json:"paths"`
}

type Method_t struct {
	Tags        []string        `json:"tags"`
	Summary     string          `json:"summary"`
	OperationId string          `json:"operationId"`
	Produces    []string        `json:"produces"`
	Parameters  []*Parameters_t `json:"parameters"`
	Deprecated  bool            `json:"deprecated"`
	Security    []*interface{}  `json:"security"`
}

type Parameters_t struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type Security_t struct{}

func (endpoint *Endpoint_t) SetPath() {

	var apivars []string

	re := regexp.MustCompile(`{([a-zA-Z0-9]+)}`)
	revars := re.FindAllStringSubmatch(endpoint.ApiPath, -1)

	for k := range revars {
		apivars = append(apivars, (revars[k])[1])
	}

	if apivars == nil {
		endpoint.IsRegex = false
		log.Debugf("No vars found for endpoint %s", endpoint.ApiPath)
		endpoint.Path = (endpoint.Prefix + endpoint.ApiPath)
		return
	}

	endpoint.IsRegex = true
	isEol, _ := regexp.MatchString(`{([a-zA-Z0-9]+)}$`, endpoint.ApiPath)

	if isEol {
		endpoint.Path = (endpoint.Prefix + (re.ReplaceAllString(endpoint.ApiPath, "(?<$1>([a-zA-Z0-9])+)")) + "$")
	} else {
		endpoint.Path = (endpoint.Prefix + (re.ReplaceAllString(endpoint.ApiPath, "(?<$1>([a-zA-Z0-9])+)")))
	}

}

func (endpoint *Endpoint_t) SetId() {

	k := md5.Sum([]byte(endpoint.ApiPath))
	endpoint.Id = hex.EncodeToString(k[:8])

}

func (endpoint *Endpoint_t) SetPrefix(prefix string) {

	endpoint.Prefix = prefix

}

func (endpoint *Endpoint_t) SetUpstream(target string) {

	endpoint.Upstream = target

}

func (endpoint *Endpoint_t) SetIsSecured() {

	log.Debugf("Processing security flag %s", strconv.FormatBool(endpoint.IsSecured))

	endpoint.IsSecured = true

	for k, v := range endpoint.Methods {

		endpoint.IsSecured = true

		if v.Security == nil {

			log.Debugf("Processing true security flag %+v %s", v.Security, k)
			endpoint.IsSecured = true
		} else {
			log.Debugf("Processing false security flag %+v %s", v.Security, k)
			endpoint.IsSecured = false
		}
	}

}

func (endpoint *Endpoint_t) SetLimits() {

	for i := range endpoint.Methods {
		endpoint.Limits = append(endpoint.Limits, strings.ToUpper(i))
	}

}

func GetSwaggerFromUrl(url string) *Swagger_t {

	var target Swagger_t

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	var myClient = &http.Client{Timeout: 10 * time.Second, Transport: tr}

	r, err := myClient.Get(url)
	if err != nil {
		log.Errorf("Unable to connect %s", url)
		return &target
	}
	defer r.Body.Close()

	resbody, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(resbody, &target)

	return &target
}

func (target *Swagger_t) TestSwaggerObject() {

	for i := range target.Paths { //Go through all paths
		fmt.Println("Key is: ", i)

		for l, m := range target.Paths[i] { //... and all methods

			fmt.Println("Method is: ", l)

			for o := range m.Parameters { //... and all parameters
				println("Name: ", m.Parameters[o].Name)
				println("In: ", m.Parameters[o].In)
				println("Description: ", m.Parameters[o].Description)
			}
		}
	}

}

func (target *Swagger_t) GetEndpoint(path string) *Endpoint_t {

	var endpoint Endpoint_t

	for i := range target.Paths {

		if i != path {
			continue
		}

		endpoint.ApiPath = i
		endpoint.Methods = make(map[string]*Method_t)

		for l, m := range target.Paths[i] {
			endpoint.Methods[l] = m
		}

	}

	return &endpoint

}

func (target *Swagger_t) GetEndpointsFromSwagger() []Endpoint_t {

	var endpoints []Endpoint_t

	for i := range target.Paths {
		endpoint := target.GetEndpoint(i)
		endpoint.SetPath()
		endpoint.SetUpstream(target.Upstream)
		endpoint.SetPrefix(target.Prefix)
		endpoint.SetLimits()
		endpoint.SetId()
		endpoint.SetIsSecured()
		endpoints = append(endpoints, *endpoint)
	}

	return endpoints

}

func (target *Swagger_t) SetSwaggerPrefix(prefix string) {

	target.Prefix = prefix

}

func (target *Swagger_t) SetSwaggerUpstream(upstream string) {

	target.Upstream = upstream

}

func (target *Swagger_t) GetEndpointsFromPaths(paths []string) []Endpoint_t {

	var endpoints []Endpoint_t

	log.Debugf("Paths definition found: %s", paths)

	for _, k := range paths {

		endpoint := target.GetEndpoint(k)
		//		log.Debugf("Endpoint processed %s", *endpoint)

		endpoint.SetPath()

		if endpoint.Path == "" {
			log.Infof("Endpoint %s was not found in swagger docs", k)
			continue
		}

		endpoint.SetUpstream(target.Upstream)
		endpoint.SetPrefix(target.Prefix)
		endpoint.SetLimits()
		endpoint.SetId()

		log.Debugf("Adding endpoint %s to the endpoint array", k)
		endpoints = append(endpoints, *endpoint)
	}

	//	log.Debugf ("Endpoints created %s", endpoints)

	return endpoints

}
