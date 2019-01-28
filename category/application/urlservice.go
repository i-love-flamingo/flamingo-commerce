package application

// URLWithName points to a category with a given name
func URLWithName(code, name string) (string, map[string]string) {
	return "category.view", map[string]string{"code": code, "name": name}
}

// URL to page with name
func URL(code string) (string, map[string]string) {
	return "category.view", map[string]string{"code": code}
}
