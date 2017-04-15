package controller

import (
	"flamingo/framework/web"
	"fmt"
	"net/http"
	"flamingo/framework/web/responder"
	"flamingo/core/event2"
	"bytes"
	"flamingo/core/cart/domain"
)

type (
	TestLoginController struct {
		*responder.RenderAware `inject:""`

		EventDispatcher event2.EventDispatcher `inject:""`

		//pageservice interfaces.PageService
	}


)



// Testlogin - throws an event
func (lc *TestLoginController) Get(c web.Context) web.Response {
	fmt.Println("Test login yeah :-)")

	lc.EventDispatcher.Dispatch(
		domain.NewLoginSucessEvent("U123213"),
	)

	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        bytes.NewReader([]byte("login sucess test")),
		ContentType: "text/html; charset=utf-8",
	}

}


