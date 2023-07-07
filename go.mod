module github.com/tierklinik-dobersberg/cis-idm

go 1.20

require (
	github.com/Masterminds/sprig/v3 v3.2.3
	github.com/bufbuild/connect-go v1.8.0
	github.com/bufbuild/protovalidate-go v0.1.1
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-webauthn/webauthn v0.8.2
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/hashicorp/go-multierror v1.1.1
	github.com/kevinburke/twilio-go v0.0.0-20221122012537-65f3dd7539e2
	github.com/mdp/qrterminal v1.0.1
	github.com/mennanov/fmutils v0.2.0
	github.com/mileusna/useragent v1.3.3
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pquerna/otp v1.4.0
	github.com/rqlite/gorqlite v0.0.0-20230310040812-ec5e524a562e
	github.com/rs/cors v1.9.0
	github.com/sethvargo/go-envconfig v0.9.0
	github.com/sirupsen/logrus v1.9.2
	github.com/spf13/cobra v1.7.0
	github.com/tierklinik-dobersberg/apis v0.0.0-20230601061851-2b5e50954244
	github.com/vincent-petithory/dataurl v1.0.0
	golang.org/x/crypto v0.9.0
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1
	golang.org/x/net v0.10.0
	golang.org/x/sync v0.1.0
	google.golang.org/protobuf v1.30.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.30.0-20230530223247-ca37dc8895db.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.0 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.0.0-20230512164433-5d1fd1a340c9 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/fxamacker/cbor/v2 v2.4.0 // indirect
	github.com/go-webauthn/revoke v0.1.9 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/cel-go v0.16.0 // indirect
	github.com/google/go-tpm v0.3.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kevinburke/go-types v0.0.0-20210723172823-2deba1f80ba7 // indirect
	github.com/kevinburke/rest v0.0.0-20210506044642-5611499aa33c // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.16.0 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/ttacon/builder v0.0.0-20170518171403-c099f663e1c2 // indirect
	github.com/ttacon/libphonenumber v1.2.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/term v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/qr v0.2.0 // indirect
)

replace github.com/tierklinik-dobersberg/apis => ../apis
