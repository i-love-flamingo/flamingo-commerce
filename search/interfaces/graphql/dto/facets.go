package dto

import searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"

// Commerce_Search_Facet interface for facets
type Commerce_Search_Facet interface {
	Name() string
	Label() string
	Position() int
	//Items() []Commerce_Search_FacetItem
	HasSelectedItem() bool
}

// Commerce_Search_Facet interface for facet items
type Commerce_Search_FacetItem interface {
	Label() string
	Value() string
	Selected() bool
	Count() int
}

// WrapListFacet wraps the list facet into the graphql dto
func WrapListFacet(facet searchdomain.Facet) *Commerce_Search_ListFacet {
	items := make([]*Commerce_Search_ListFacetItem, len(facet.Items))
	for i, item := range facet.Items {
		items[i] = &Commerce_Search_ListFacetItem{
			label:    item.Label,
			value:    item.Value,
			selected: item.Selected,
			count:    int(item.Count),
		}
	}

	return &Commerce_Search_ListFacet{
		name:     facet.Name,
		label:    facet.Label,
		position: facet.Position,
		items:    items,
	}
}

// Commerce_Search_ListFacet dto for list facets
type Commerce_Search_ListFacet struct {
	name     string
	label    string
	position int
	items    []*Commerce_Search_ListFacetItem
}

// Name getter
func (c *Commerce_Search_ListFacet) Name() string {
	return c.name
}

// Label getter
func (c *Commerce_Search_ListFacet) Label() string {
	return c.label
}

// Position getter
func (c *Commerce_Search_ListFacet) Position() int {
	return c.position
}

// Items getter
func (c *Commerce_Search_ListFacet) Items() []*Commerce_Search_ListFacetItem {
	return c.items
}

// HasSelectedItem getter
func (c *Commerce_Search_ListFacet) HasSelectedItem() bool {
	for _, item := range c.items {
		if item.selected {
			return true
		}
	}
	return false
}

// Commerce_Search_ListFacetItem dto for list facet items
type Commerce_Search_ListFacetItem struct {
	label    string
	value    string
	selected bool
	count    int
}

// Label getter
func (c *Commerce_Search_ListFacetItem) Label() string {
	return c.label
}

// Value getter
func (c *Commerce_Search_ListFacetItem) Value() string {
	return c.value
}

// Selected getter
func (c *Commerce_Search_ListFacetItem) Selected() bool {
	return c.selected
}

// Count getter
func (c *Commerce_Search_ListFacetItem) Count() int {
	return c.count
}

func mapTreeFacetItems(facetItems []*searchdomain.FacetItem) []*Commerce_Search_TreeFacetItem {
	items := make([]*Commerce_Search_TreeFacetItem, len(facetItems))
	for i, item := range facetItems {
		items[i] = &Commerce_Search_TreeFacetItem{
			label:    item.Label,
			value:    item.Value,
			selected: item.Selected,
			count:    int(item.Count),
			active:   item.Active,
			items:    mapTreeFacetItems(item.Items),
		}
	}
	return items
}

// WrapTreeFacet wraps the tree facet into the graphql dto
func WrapTreeFacet(facet searchdomain.Facet) *Commerce_Search_TreeFacet {
	return &Commerce_Search_TreeFacet{
		name:     facet.Name,
		label:    facet.Label,
		position: facet.Position,
		items:    mapTreeFacetItems(facet.Items),
	}
}

// Commerce_Search_TreeFacet dto for tree facets
type Commerce_Search_TreeFacet struct {
	name     string
	label    string
	position int
	items    []*Commerce_Search_TreeFacetItem
}

// Name getter
func (c *Commerce_Search_TreeFacet) Name() string {
	return c.name
}

// Label getter
func (c *Commerce_Search_TreeFacet) Label() string {
	return c.label
}

// Position getter
func (c *Commerce_Search_TreeFacet) Position() int {
	return c.position
}

// Items getter
func (c *Commerce_Search_TreeFacet) Items() []*Commerce_Search_TreeFacetItem {
	return c.items
}

func hasSelectedItem(items []*Commerce_Search_TreeFacetItem) bool {
	for _, item := range items {
		if item.selected || hasSelectedItem(item.items) {
			return true
		}
	}
	return false
}

// HasSelectedItem getter
func (c *Commerce_Search_TreeFacet) HasSelectedItem() bool {
	if hasSelectedItem(c.items) {
		return true
	}
	return false
}

// Commerce_Search_TreeFacetItem dto for tree facet items
type Commerce_Search_TreeFacetItem struct {
	label    string
	value    string
	selected bool
	count    int
	active   bool
	items    []*Commerce_Search_TreeFacetItem
}

// Label getter
func (c *Commerce_Search_TreeFacetItem) Label() string {
	return c.label
}

// Value getter
func (c *Commerce_Search_TreeFacetItem) Value() string {
	return c.value
}

// Selected getter
func (c *Commerce_Search_TreeFacetItem) Selected() bool {
	return c.selected
}

// Count getter
func (c *Commerce_Search_TreeFacetItem) Count() int {
	return c.count
}

// Active getter
func (c *Commerce_Search_TreeFacetItem) Active() bool {
	return c.active
}

// Items getter
func (c *Commerce_Search_TreeFacetItem) Items() []*Commerce_Search_TreeFacetItem {
	return c.items
}

// WrapRangeFacet wraps the range facet into the graphql dto
func WrapRangeFacet(facet searchdomain.Facet) *Commerce_Search_RangeFacet {
	items := make([]*Commerce_Search_RangeFacetItem, len(facet.Items))
	for i, item := range facet.Items {
		items[i] = &Commerce_Search_RangeFacetItem{
			label:       item.Label,
			value:       item.Value,
			selected:    item.Selected,
			count:       int(item.Count),
			min:         int(item.Min),
			max:         int(item.Max),
			selectedMin: int(item.SelectedMin),
			selectedMax: int(item.SelectedMax),
		}
	}

	return &Commerce_Search_RangeFacet{
		name:     facet.Name,
		label:    facet.Label,
		position: facet.Position,
		items:    items,
	}
}

// Commerce_Search_RangeFacet dto for range facets
type Commerce_Search_RangeFacet struct {
	name     string
	label    string
	position int
	items    []*Commerce_Search_RangeFacetItem
}

// Name getter
func (c *Commerce_Search_RangeFacet) Name() string {
	return c.name
}

// Label getter
func (c *Commerce_Search_RangeFacet) Label() string {
	return c.label
}

// Position getter
func (c *Commerce_Search_RangeFacet) Position() int {
	return c.position
}

// Items getter
func (c *Commerce_Search_RangeFacet) Items() []*Commerce_Search_RangeFacetItem {
	return c.items
}

// HasSelectedItem getter
func (c *Commerce_Search_RangeFacet) HasSelectedItem() bool {
	for _, item := range c.items {
		if item.selected {
			return true
		}
	}
	return false
}

// Commerce_Search_RangeFacetItem dto for range facet items
type Commerce_Search_RangeFacetItem struct {
	label       string
	value       string
	selected    bool
	count       int
	min         int
	max         int
	selectedMin int
	selectedMax int
}

// Label getter
func (c *Commerce_Search_RangeFacetItem) Label() string {
	return c.label
}

// Value getter
func (c *Commerce_Search_RangeFacetItem) Value() string {
	return c.value
}

// Selected getter
func (c *Commerce_Search_RangeFacetItem) Selected() bool {
	return c.selected
}

// Count getter
func (c *Commerce_Search_RangeFacetItem) Count() int {
	return c.count
}

// Min getter
func (c *Commerce_Search_RangeFacetItem) Min() int {
	return c.min
}

// Max getter
func (c *Commerce_Search_RangeFacetItem) Max() int {
	return c.max
}

// SelectedMin getter
func (c *Commerce_Search_RangeFacetItem) SelectedMin() int {
	return c.selectedMin
}

// SelectedMax getter
func (c *Commerce_Search_RangeFacetItem) SelectedMax() int {
	return c.selectedMax
}
