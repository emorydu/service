package auditapp

import (
	"net/http"

	"github.com/emorydu/service/app/sdk/auth"
	"github.com/emorydu/service/app/sdk/authclient"
	"github.com/emorydu/service/app/sdk/mid"
	"github.com/emorydu/service/business/domain/auditbus"
	"github.com/emorydu/service/foundation/logger"
	"github.com/emorydu/service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	AuditBus   auditbus.ExtBusiness
	AuthClient authclient.Authenticator
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.AuditBus)

	app.HandlerFunc(http.MethodGet, version, "/audits", api.query, authen, ruleAdmin)
}
