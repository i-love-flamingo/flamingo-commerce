package controller

import (
	"bytes"
	"flamingo/core/cart/domain"
	"flamingo/framework/event"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"fmt"
	"net/http"
)

type (
	TestLoginController struct {
		*responder.RenderAware `inject:""`

		EventRouter event.Router `inject:""`

		//pageservice interfaces.PageService
	}
)

// Testlogin - throws an event
func (lc *TestLoginController) Get(c web.Context) web.Response {
	fmt.Println("Test login yeah :-)")

	lc.EventRouter.Dispatch(
		domain.LoginSucessEvent{
			"U213213",
		},
	)

	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        bytes.NewReader([]byte("login sucess test")),
		ContentType: "text/html; charset=utf-8",
	}

}
