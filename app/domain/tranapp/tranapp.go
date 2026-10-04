// Package tranapp maintains the app layer api for the tran domain.
package tranapp

import (
	"context"
	"errors"
	"net/http"

	"github.com/emorydu/service/app/sdk/errs"
	"github.com/emorydu/service/app/sdk/mid"
	"github.com/emorydu/service/business/domain/productbus"
	"github.com/emorydu/service/business/domain/userbus"
	"github.com/emorydu/service/foundation/web"
)

type app struct {
	userBus    userbus.ExtBusiness
	productBus productbus.ExtBusiness
}

func newApp(userBus userbus.ExtBusiness, productBus productbus.ExtBusiness) *app {
	return &app{
		userBus:    userBus,
		productBus: productBus,
	}
}

// newWithTx constructs a new Handlers value with the domain apis
// using a store transaction that was created via middleware.
func (a *app) newWithTx(ctx context.Context) (*app, error) {
	tx, err := mid.GetTran(ctx)
	if err != nil {
		return nil, err
	}

	userBus, err := a.userBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	productBus, err := a.productBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	txApp := app{
		userBus:    userBus,
		productBus: productBus,
	}

	return &txApp, nil
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var txApp NewTran
	if err := web.Decode(r, &txApp); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	np, err := toBusNewProduct(txApp.Product)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nu, err := toBusNewUser(txApp.User)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Create(ctx, mid.GetSubjectID(ctx), nu)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return errs.Errorf(errs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	np.UserID = usr.ID

	prd, err := a.productBus.Create(ctx, np)
	if err != nil {
		return errs.Errorf(errs.Internal, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd)
}
