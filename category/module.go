package category

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"
	flamingographql "flamingo.me/graphql"

	"flamingo.me/flamingo-commerce/v3/category/application"
	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/category/infrastructure"
	"flamingo.me/flamingo-commerce/v3/category/infrastructure/fake"
	"flamingo.me/flamingo-commerce/v3/category/interfaces/controller"
	categoryGraphql "flamingo.me/flamingo-commerce/v3/category/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/product"
	productApplication "flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/search"
)

// Module registers our profiler
type Module struct {
	useCategoryFixedAdapter bool
	useFakeService          bool
}

// URL to category
func URL(code string) (string, map[string]string) {
	return application.URL(code)
}

// URLWithName to category
func URLWithName(code, name string) (string, map[string]string) {
	return application.URLWithName(code, web.URLTitle(name))
}

// Inject dependencies
func (m *Module) Inject(
	config *struct {
		UseCategoryFixedAdapter bool `inject:"config:commerce.category.useCategoryFixedAdapter,optional"`
		UseFakeService          bool `inject:"config:commerce.category.fakeService.enabled,optional"`
	},
) {
	if config != nil {
		m.useCategoryFixedAdapter = config.UseCategoryFixedAdapter
		m.useFakeService = config.UseFakeService
	}
}

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(controller.QueryHandler)).To(controller.QueryHandlerImpl{})
	injector.Bind(new(controller.ProductSearchService)).To(productApplication.ProductSearchService{})

	if m.useCategoryFixedAdapter {
		injector.Bind((*domain.CategoryService)(nil)).To(infrastructure.CategoryServiceFixed{})
	}
	if m.useFakeService {
		injector.Override((*domain.CategoryService)(nil), "").To(fake.CategoryService{}).In(dingo.ChildSingleton)
	}
	web.BindRoutes(injector, new(routes))
	injector.Bind(new(application.RouterRouter)).To(new(web.Router))
	injector.BindMulti(new(flamingographql.Service)).To(categoryGraphql.Service{})
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(product.Module),
		new(search.Module),
	}
}

type routes struct {
	view   *controller.ViewController
	entity *controller.Entity
	tree   *controller.Tree
}

func (r *routes) Inject(view *controller.ViewController, entity *controller.Entity, tree *controller.Tree) {
	r.view = view
	r.entity = entity
	r.tree = tree
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleGet("category.view", r.view.Get)
	handler, _ := registry.Route("/category/:code/:name.html", "category.view(code, name, *)")
	handler.Normalize("name")
	registry.MustRoute("/category/:code", "category.view(code, *)")

	registry.HandleData("category.tree", r.tree.Data)
	registry.HandleData("category", r.entity.Data)
}

// CueConfig defines the category module configuration
func (*Module) CueConfig() string {
	return `
commerce: {
	CategoryTree :: {
		[string]: CategoryTreeNode
	}
	CategoryTreeNode :: {
		code: string
		name: string
		sort?: number
		childs?: CategoryTree
	}

	category: {
		view:  {
			template: *"category/category" | !=""
			teaserTemplate: *"category/teaser" | !=""
		}
		useCategoryFixedAdapter: bool | *false
		if useCategoryFixedAdapter {
			categoryServiceFixed: {
      			tree: CategoryTree
			}
		}
		fakeService: {
			enabled: bool | *false
			if enabled {
			  testDataFolder?: string & !=""
			}
		}
	}
}`
}
