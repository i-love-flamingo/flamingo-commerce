package models

type (
	ProductImage struct {
		Type  string
		Alt   string
		Title string
		Urls  struct {
			Base string
			Xss  string
			Xs   string
			Sm   string
			Md   string
			Lg   string
			Xl   string
		}
	}

	Product struct {
		Id    string
		Name  string
		Brand struct {
			Id   string
			Name string
		}
		Retailer struct {
			Id   string
			Name string
		}
		Description string
		Images      []ProductImage
		Attributes  []struct {
			ID    string
			Name  string
			Value string
		}
		Prices struct {
			Base float64
		}
		Shipping []struct {
			Id        string
			Available bool
			Title     string
			Duration  string
		}
		Categories []struct {
			ID   string
			Name string
		}
	}
)
