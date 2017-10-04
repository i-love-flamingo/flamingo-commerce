package controller

import (
	"bytes"
	"fmt"
	"net/http"

	"go.aoe.com/flamingo/core/cart/domain"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// TestLoginController for testing
	TestLoginController struct {
		responder.RenderAware `inject:""`

		EventRouter event.Router `inject:""`

		//pageservice interfaces.PageService
	}
)

// Get Testlogin - throws an event
func (lc *TestLoginController) Get(c web.Context) web.Response {
	fmt.Println("Test login yeah :-)")

	lc.EventRouter.Dispatch(
		domain.LoginSucessEvent{
			UserID: "U213213",
		},
	)

	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        bytes.NewReader([]byte("login sucess test")),
		ContentType: "text/html; charset=utf-8",
	}

}
