// Package userapp maintains the app layer api for the user domain.
package userapp

import (
	"context"
	"errors"
	"net/http"

	"github.com/emorydu/service/app/sdk/errs"
	"github.com/emorydu/service/app/sdk/mid"
	"github.com/emorydu/service/app/sdk/query"
	"github.com/emorydu/service/business/domain/userbus"
	"github.com/emorydu/service/business/sdk/order"
	"github.com/emorydu/service/business/sdk/page"
	"github.com/emorydu/service/foundation/web"
)

type app struct {
	userBus userbus.ExtBusiness
}

func newApp(userBus userbus.ExtBusiness) *app {
	return &app{userBus: userBus}
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var nu NewUser
	if err := web.Decode(r, &nu); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	bnu, err := toBusNewUser(nu)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Create(ctx, mid.GetSubjectID(ctx), bnu)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return errs.Errorf(errs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	return toAppUser(usr)
}

func (a *app) update(ctx context.Context, r *http.Request) web.Encoder {
	var uu UpdateUser
	if err := web.Decode(r, &uu); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	buu, err := toBusUpdateUser(uu)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := mid.GetUser(ctx)
	if err != nil {
		return errs.Errorf(errs.Internal, "user missing in context: %s", err)
	}

	updUsr, err := a.userBus.Update(ctx, mid.GetSubjectID(ctx), usr, buu)
	if err != nil {
		return errs.Errorf(errs.Internal, "update: userID[%s] uu[%+v]: %s", usr.ID, buu, err)
	}

	return toAppUser(updUsr)
}

func (a *app) updateRole(ctx context.Context, r *http.Request) web.Encoder {
	var uur UpdateUserRole
	if err := web.Decode(r, &uur); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	buu, err := toBusUpdateUserRole(uur)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := mid.GetUser(ctx)
	if err != nil {
		return errs.Errorf(errs.Internal, "user missing in context: %s", err)
	}

	updUsr, err := a.userBus.Update(ctx, mid.GetSubjectID(ctx), usr, buu)
	if err != nil {
		return errs.Errorf(errs.Internal, "updaterole: userID[%s] uu[%+v]: %s", usr.ID, buu, err)
	}

	return toAppUser(updUsr)
}

func (a *app) delete(ctx context.Context, _ *http.Request) web.Encoder {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return errs.Errorf(errs.Internal, "userID missing in context: %s", err)
	}

	if err := a.userBus.Delete(ctx, mid.GetSubjectID(ctx), usr); err != nil {
		return errs.Errorf(errs.Internal, "delete: userID[%s]: %s", usr.ID, err)
	}

	return nil
}

func (a *app) query(ctx context.Context, r *http.Request) web.Encoder {
	qp, err := parseQueryParams(r)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, userbus.DefaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	usrs, err := a.userBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query: %s", err)
	}

	total, err := a.userBus.Count(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppUsers(usrs), total, page)
}

func (a *app) queryByID(ctx context.Context, _ *http.Request) web.Encoder {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return errs.Errorf(errs.Internal, "querybyid: %s", err)
	}

	return toAppUser(usr)
}
