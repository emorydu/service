package tranapp

import (
	"net/http"

	"github.com/emorydu/service/app/sdk/auth"
	"github.com/emorydu/service/app/sdk/authclient"
	"github.com/emorydu/service/app/sdk/mid"
	"github.com/emorydu/service/business/domain/productbus"
	"github.com/emorydu/service/business/domain/userbus"
	"github.com/emorydu/service/business/sdk/sqldb"
	"github.com/emorydu/service/foundation/logger"
	"github.com/emorydu/service/foundation/web"
	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	DB         *sqlx.DB
	UserBus    userbus.ExtBusiness
	ProductBus productbus.ExtBusiness
	AuthClient authclient.Authenticator
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	transaction := mid.BeginCommitRollback(cfg.Log, sqldb.NewBeginner(cfg.DB))
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.UserBus, cfg.ProductBus)

	app.HandlerFunc(http.MethodPost, version, "/tranexample", api.create, authen, ruleAdmin, transaction)
}
