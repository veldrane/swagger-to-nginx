package help

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	SwaggerUrl string
	OutFile    string
	Upstream   string
	Prefix     string
	Template   string
	Inner      string
	Paths      []string
	Patharg    string
	Skel       bool
	Debug      bool
}

func SetArgs() *Config {

	var config Config

	flag.StringVar(&config.SwaggerUrl, "swagger", "", "Swagger doc url point like http://simpleapi.lab.local/swagger-ui/openapi3.json")
	flag.StringVar(&config.OutFile, "out", "location.out", "Output nginx location file")
	flag.StringVar(&config.Upstream, "upstream", "", "Target service upstream on the backend")
	flag.StringVar(&config.Prefix, "prefix", "", "Path prefix for backend specification")
	flag.StringVar(&config.Template, "template", "", "Path for location template external templatefile")
	flag.StringVar(&config.Inner, "inner", "Default", "specify inner template name Default|NoJWT")
	flag.StringVar(&config.Patharg, "paths", "", "Create location just for specified paths (separated by columns)")
	flag.BoolVar(&config.Skel, "skel", false, "Generate template skell from default template and finish")
	flag.BoolVar(&config.Debug, "debug", false, "Enable debug mode")
	flag.Parse()

	// ok thats copy from stackowerflow :)
	// so whats going on...
	// first of all define all mandatory params
	// required := []string{"swagger","skell"}

	// first make a make of all parameters like key with true/false value
	seen := make(map[string]bool)

	// ...then label all present parameters with true...
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })

	if !seen["swagger"] && !seen["skel"] {
		fmt.Fprintf(os.Stderr, "-swagger or -skel parameter is required \n")
		os.Exit(2)
	}

	if seen["swagger"] {

		if !seen["upstream"] {
			fmt.Fprintf(os.Stderr, "-upstream parameter is required with swagger docs defintion \n")
			os.Exit(2)
		}
	}

	// Must be refactored - map must be generated from templates.go not hardcoded!

	if seen["inner"] {
		innerval := map[string]bool{"Default": true, "Jwt": true}
		if innerval[config.Inner] != true {
			fmt.Fprintf(os.Stderr, "Specified invalid name of internal template \n")
			os.Exit(2)
		}
	}

	// ..ok now lets a make a loop through all PRESENT parameters on the cmd line
	//	for _, req := range required {

	// ...and because we have map with the boolean we can make a check o
	// present params!
	//		if !seen[req] {

	// if the params doesnt exists throw the error code and exist
	//			fmt.Fprintf(os.Stderr, "missing required -%s argument \n", req)
	//			os.Exit(2) // the same exit code flag.Parse uses
	//		}
	//	}

	if config.Patharg != "" {
		config.Paths = strings.Split(config.Patharg, ",")
	}

	return &config
}
