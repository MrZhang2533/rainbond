package controller

import (
	"net/http"

	"github.com/goodrain/rainbond/api/handler"
	"github.com/goodrain/rainbond/api/middleware"
	"github.com/goodrain/rainbond/api/model"
	"github.com/goodrain/rainbond/api/util/bcode"
	"github.com/goodrain/rainbond/db"
	httputil "github.com/goodrain/rainbond/util/http"
	"github.com/sirupsen/logrus"
)

// AddConfigGroup -
func (a *ApplicationStruct) AddConfigGroup(w http.ResponseWriter, r *http.Request) {
	var configReq model.ApplicationConfigGroup
	appID := r.Context().Value(middleware.ContextKey("app_id")).(string)
	if !httputil.ValidatorRequestStructAndErrorResponse(r, w, &configReq, nil) {
		return
	}

	// Get the application bound serviceIDs
	availableServices := db.GetManager().TenantServiceDao().GetServiceIDsByAppID(appID)
	// Judge whether the requested service ID is correct
	set := make(map[string]struct{})
	for _, s := range availableServices {
		set[s.ServiceID] = struct{}{}
	}
	for _, sid := range configReq.ServiceIDs {
		_, ok := set[sid]
		if !ok {
			logrus.Infof("The serviceID [%s] is not under the application or does not exist", sid)
			httputil.ReturnBcodeError(r, w, bcode.ErrServiceNotFound)
			return
		}
	}

	// create app ConfigGroups
	app, err := handler.GetApplicationHandler().AddConfigGroup(appID, &configReq)
	if err != nil {
		httputil.ReturnBcodeError(r, w, err)
		return
	}
	httputil.ReturnSuccess(r, w, app)
}
