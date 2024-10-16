// +heroku goVersion go1.18

module github.com/heroku/cnb-shim

go 1.18

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/buildpack/libbuildpack v1.25.11
	github.com/google/uuid v1.3.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/heroku/rollrus v0.2.0
	github.com/joeshaw/envdecode v0.0.0-20200121155833-099f1fc765bd
	github.com/joho/godotenv v1.4.0
	github.com/rollbar/rollbar-go v1.4.5
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.1
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
