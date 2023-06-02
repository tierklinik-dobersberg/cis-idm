module github.com/tierklinik-dobersberg/cis-userd

go 1.20

require (
	github.com/bufbuild/connect-go v1.8.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/mitchellh/mapstructure v1.5.0
	github.com/rqlite/gorqlite v0.0.0-20230310040812-ec5e524a562e
	github.com/sethvargo/go-envconfig v0.9.0
	github.com/sirupsen/logrus v1.9.2
	github.com/tierklinik-dobersberg/apis v0.0.0-20230601061851-2b5e50954244
	golang.org/x/crypto v0.9.0
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1
	golang.org/x/net v0.10.0
)

require (
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/tierklinik-dobersberg/apis => ../apis
