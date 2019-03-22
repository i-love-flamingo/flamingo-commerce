package domain

type (

	// Tree domain model
	Tree interface {
		//Code - returns the idendifier of the category
		Code() string
		//Name - returns the name of the category
		Name() string
		//Path returns the Path as string
		Path() string
		//Active - should return true if the category is in the rootpath of the current category
		Active() bool
		//Subtrees returns a list of subtrees of the current node
		SubTrees() []Tree
		//HasChilds returns true if the node is no leaf node
		HasChilds() bool
		//DocumentCount - the amount of documents (products) in the category
		DocumentCount() int
	}

	// TreeData defines the default domain tree data model
	TreeData struct {
		CategoryCode          string
		CategoryName          string
		CategoryPath          string
		CategoryDocumentCount int
		SubTreesData          []*TreeData
		IsActive              bool
	}
)

// Active gets the node (category) active state
func (c TreeData) Active() bool {
	return c.IsActive
}

// Code gets the category code represented by this node in the tree
func (c TreeData) Code() string {
	return c.CategoryCode
}

// Name gets the category name
func (c TreeData) Name() string {
	return c.CategoryName
}

// Path gets the Node (category) path
func (c TreeData) Path() string {
	return c.CategoryPath
}

// DocumentCount gets the amount of documents in that node
func (c TreeData) DocumentCount() int {
	return c.CategoryDocumentCount
}

// SubTrees gets the child Trees
func (c TreeData) SubTrees() []Tree {
	result := make([]Tree, len(c.SubTreesData))
	for i, child := range c.SubTreesData {
		result[i] = Tree(child)
	}

	return result
}

// HasChilds - true if subTrees exist
func (c TreeData) HasChilds() bool {
	return len(c.SubTreesData) > 0
}
