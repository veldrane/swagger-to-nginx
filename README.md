### Swag2Nginx tool

This tool provides converting openapi definition into the stuff based on the golang template. It was primarly used for generating nginx configuration
for custom API Gateway, but int can be used for any cfg file, if know how to describe the conversion by templates.


#### Build

```
$ cd src
$ make install
$ cd ../bin
$ ./swag2nginx --help
Usage of ./swag2nginx:
  -debug
    	Enable debug mode
  -inner string
    	specify inner template name Default|NoJWT (default "Default")
  -out string
    	Output nginx location file (default "location.out")
  -paths string
    	Create location just for specified paths (separated by columns)
  -prefix string
    	Path prefix for backend specification
  -skel
    	Generate template skell from default template and finish
  -swagger string
    	Swagger doc url point like http://simpleapi.lab.local/swagger-ui/openapi3.json
  -template string
    	Path for location template external templatefile
  -upstream string
    	Target service upstream on the backend
```

#### Usage


```
$ ./swag2nginx -swagger http://simpleapi.lab.local/swagger-ui/openapi3.json -prefix /simple -upstream svc-simple-api.simple.svc.cluster.local -out /tmp/location.out
time="2023-08-18T16:14:11+02:00" level=info msg="Using url http://simpleapi.lab.local/swagger-ui/openapi3.json to get swagger data"
time="2023-08-18T16:14:11+02:00" level=info msg="Setting up the output file /tmp/location.out..."
time="2023-08-18T16:14:11+02:00" level=info msg="Used internal template Default"
time="2023-08-18T16:14:11+02:00" level=info msg="Done!"
```

takes the default template (its part of the binary) and convert them to piece of the nginx location:

```
location = /simple/articles {
        proxy_read_timeout 45s;
        set $upstream svc-simple-api.simple.svc.cluster.local;
        set $locid 3016b1a32453def5;
        include /etc/nginx/conf.d/include/default-cors.loc;
        include /etc/nginx/conf.d/include/default-headers.loc;
        limit_except OPTIONS GET {
          deny all;
        }
        if ($request_uri ~* "/simple/(.*))" {
         proxy_pass http://$upstream/$1;
        }
}
```

the default template is present on the ../template subdir:

```
{{ range . -}}
location {{- if .IsRegex}} ~ {{- else}} = {{- end }} {{ .Prefix}}{{.Path}} {
        proxy_read_timeout 45s;
        set $upstream {{.Upstream}};
        set $locid {{.Id}};
        include /etc/nginx/conf.d/include/default-cors.loc;
        include /etc/nginx/conf.d/include/default-headers.loc;
        limit_except OPTIONS{{range .Limits}} {{.}}{{end}} {
          deny all;
        }
        if ($request_uri ~* "{{.Prefix}}/(.*))" {
         proxy_pass http://$upstream/$1;
        }
}
{{ end -}}
```

but fell free to write own more sophisticated one.

#### API

Inside of the template you can use lots of openapi property like:

- .Upstream - upstream server: defined by the command line
- .Prefix - prefix path 
- .Limits - used methods 
etc - tbd.

Api keys use the nginx language. 

