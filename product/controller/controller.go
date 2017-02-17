package controller

import (
	"encoding/json"
	"flamingo/core/app/web"
	"flamingo/core/app/web/responder"
	"flamingo/core/product/interfaces"
	"io/ioutil"
)

type ViewController struct {
	*responder.RenderAware `inject:""`

	interfaces.ProductService `inject:""`
}

func (vc *ViewController) Get(c web.Context) web.Response {
	//product := vc.ProductService.Get(c.Param1("sku"))

	//return vc.Render(c, "pages/product/view", product)

	product := make(map[string]interface{})

	p, _ := ioutil.ReadFile("frontend/src/mocks/product.json")
	json.Unmarshal(p, &product)

	return vc.Render(c, "pages/product/view", map[string]interface{}{"Product": product})
}
