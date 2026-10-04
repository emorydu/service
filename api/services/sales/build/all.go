//go:build !crud && !reporting

// Package build manages different build options.
package build

import (
	"github.com/emorydu/service/app/domain/auditapp"
	"github.com/emorydu/service/app/domain/checkapp"
	"github.com/emorydu/service/app/domain/homeapp"
	"github.com/emorydu/service/app/domain/productapp"
	"github.com/emorydu/service/app/domain/rawapp"
	"github.com/emorydu/service/app/domain/tranapp"
	"github.com/emorydu/service/app/domain/userapp"
	"github.com/emorydu/service/app/domain/vproductapp"
	"github.com/emorydu/service/app/sdk/mux"
	"github.com/emorydu/service/foundation/web"
)

// Routes binds all the routes for the sales service.
func Routes() all {
	return all{}
}

type all struct{}

// Add implements the RouterAdder interface.
func (all) Add(app *web.App, cfg mux.Config) {
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	homeapp.Routes(app, homeapp.Config{
		Log:        cfg.Log,
		HomeBus:    cfg.BusConfig.HomeBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	productapp.Routes(app, productapp.Config{
		Log:        cfg.Log,
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	rawapp.Routes(app)

	tranapp.Routes(app, tranapp.Config{
		Log:        cfg.Log,
		DB:         cfg.DB,
		UserBus:    cfg.BusConfig.UserBus,
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	userapp.Routes(app, userapp.Config{
		Log:        cfg.Log,
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	auditapp.Routes(app, auditapp.Config{
		Log:        cfg.Log,
		AuditBus:   cfg.BusConfig.AuditBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	vproductapp.Routes(app, vproductapp.Config{
		Log:         cfg.Log,
		UserBus:     cfg.BusConfig.UserBus,
		VProductBus: cfg.BusConfig.VProductBus,
		AuthClient:  cfg.SalesConfig.AuthClient,
	})
}
