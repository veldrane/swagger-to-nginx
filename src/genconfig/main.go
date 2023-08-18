package main

import (
	"encoding/base64"
	"os"
	"text/template"
	"fmt"

	help "bitbucket.org/veldrane/golibs/help"
	swagger "bitbucket.org/veldrane/golibs/swagger"

	log "github.com/sirupsen/logrus"

)

func setLog(config *help.Config) {

	log.SetOutput(os.Stdout)

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{DisableColors: true})

}

func setSwagger(target *swagger.Swagger_t, config *help.Config) {

	target.Prefix = config.Prefix
	target.Upstream = config.Upstream
	return

}

func setTemplate(config *help.Config) *template.Template {

	// Needs custumoziation for more internal templates

	log.Debugf("Entering into setTemplate function") 

	if config.Template != "" {
		log.Infof ("Using template file %s", config.Template)
		t, err := template.ParseFiles(config.Template)

		if err != nil {
			log.Fatalf("Unable to parse template %s ", config.Template)
		}

		return t
	}

	log.Debugf("No template file found, using internal template %s", config.Inner)

	if swagger.Templates[config.Inner] == "" {
		log.Fatalf("No internal template %s found", config.Inner)
	}

	DefaultTemplate, err := base64.StdEncoding.DecodeString(swagger.Templates[config.Inner])

	if err != nil {
		log.Fatalf("Unable to decode default template via base64 %s", swagger.Templates[config.Inner])
	}

	t := template.New("Default")
	
	log.Infof("Used internal template %s", config.Inner)

	if config.Debug {
		fmt.Printf("%s", string (DefaultTemplate))
	}

	t, err = t.Parse(string(DefaultTemplate[:]))

	if err != nil {
		log.Fatalf("Unable to parse default template")
	}

	return t

}

func writeSkel (f *os.File, config *help.Config) {

	log.Infof("Writing tmplate to output file %s", config.OutFile)

	if swagger.Templates[config.Inner] == "" {
		log.Fatalf("No internal template %s found", config.Inner)
	}

	t, _ := base64.StdEncoding.DecodeString(swagger.Templates[config.Inner])
	f.Write(t)
	os.Exit(0)

}

func getSwagger(config *help.Config) *swagger.Swagger_t {

	swag := swagger.GetSwaggerFromUrl(config.SwaggerUrl)
	swag.SetSwaggerPrefix(config.Prefix)
	swag.SetSwaggerUpstream(config.Upstream)

	return swag
}

func main() {

	//var f *os.File

	config := help.SetArgs()
	setLog(config)

	f, err := os.OpenFile(config.OutFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Errorf("Cannon open output file: %s", config.OutFile)
	}
	defer f.Close()

	if config.Skel {
		writeSkel(f, config)
	}

	log.Infof("Using url %s to get swagger data", config.SwaggerUrl)
	log.Infof("Setting up the output file %s...", config.OutFile)

	swag := getSwagger(config)

	endpoints := swag.GetEndpointsFromSwagger()

	if config.Paths != nil {
		endpoints = swag.GetEndpointsFromPaths(config.Paths)
	}

	if endpoints == nil {
		log.Fatalf("No endpoints found on %s", config.SwaggerUrl)
	}

    t := setTemplate(config)

    err = t.Execute(f, endpoints)
    if err != nil {
        log.Fatalf("Error during template occurred: %s", err)
    }	

	log.Infoln("Done!")
}
