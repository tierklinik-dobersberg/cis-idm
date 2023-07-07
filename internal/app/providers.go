package app

import (
	"github.com/bufbuild/protovalidate-go"
	"github.com/tierklinik-dobersberg/cis-idm/internal/common"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/sms"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type Providers struct {
	TemplateEngine *tmpl.Engine
	SMSSender      sms.Sender
	Datastore      *repo.Repo
	Config         config.Config
	Common         *common.Service
	ProtoRegistry  *protoregistry.Files
	Validator      *protovalidate.Validator
}
