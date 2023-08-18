module genconfig

go 1.14

require (
	bitbucket.org/veldrane/golibs/help v0.0.0
	bitbucket.org/veldrane/golibs/swagger v0.0.0
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/objx v0.3.0 // indirect
)

replace (
	bitbucket.org/veldrane/golibs/help => ./local/help
	bitbucket.org/veldrane/golibs/swagger => ./local/swagger
)
