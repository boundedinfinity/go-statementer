module github.com/boundedinfinity/docsorter

go 1.18

require (
	github.com/boundedinfinity/go-commoner v1.0.23
	github.com/boundedinfinity/rfc3339date v1.0.1
	github.com/oriser/regroup v0.0.0-20210730155327-fca8d7531263
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gocarina/gocsv v0.0.0-20230123225133-763e25b40669 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/boundedinfinity/commons => ../commons

replace github.com/boundedinfinity/rfc3339date => ../rfc3339date
