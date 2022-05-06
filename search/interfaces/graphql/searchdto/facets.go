package searchdto

import searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"

// CommerceSearchFacet interface for facets
type CommerceSearchFacet interface {
	Name() string
	Label() string
	Position() int
	// Items() []CommerceSearchFacetItem
	HasSelectedItem() bool
}

// CommerceSearchFacetItem interface for facet items
type CommerceSearchFacetItem interface {
	Label() string
	Value() string
	Selected() bool
	Count() int
}

// WrapListFacet wraps the list facet into the graphql dto
func WrapListFacet(facet searchdomain.Facet) *CommerceSearchListFacet {
	items := make([]*CommerceSearchListFacetItem, len(facet.Items))
	for i, item := range facet.Items {
		items[i] = &CommerceSearchListFacetItem{
			label:    item.Label,
			value:    item.Value,
			selected: item.Selected,
			count:    int(item.Count),
		}
	}

	return &CommerceSearchListFacet{
		name:     facet.Name,
		label:    facet.Label,
		position: facet.Position,
		items:    items,
	}
}

// CommerceSearchListFacet dto for list facets
type CommerceSearchListFacet struct {
	name     string
	label    string
	position int
	items    []*CommerceSearchListFacetItem
}

// Name getter
func (c *CommerceSearchListFacet) Name() string {
	return c.name
}

// Label getter
func (c *CommerceSearchListFacet) Label() string {
	return c.label
}

// Position getter
func (c *CommerceSearchListFacet) Position() int {
	return c.position
}

// Items getter
func (c *CommerceSearchListFacet) Items() []*CommerceSearchListFacetItem {
	return c.items
}

// HasSelectedItem getter
func (c *CommerceSearchListFacet) HasSelectedItem() bool {
	for _, item := range c.items {
		if item.selected {
			return true
		}
	}
	return false
}

// CommerceSearchListFacetItem dto for list facet items
type CommerceSearchListFacetItem struct {
	label    string
	value    string
	selected bool
	count    int
}

// Label getter
func (c *CommerceSearchListFacetItem) Label() string {
	return c.label
}

// Value getter
func (c *CommerceSearchListFacetItem) Value() string {
	return c.value
}

// Selected getter
func (c *CommerceSearchListFacetItem) Selected() bool {
	return c.selected
}

// Count getter
func (c *CommerceSearchListFacetItem) Count() int {
	return c.count
}

func mapTreeFacetItems(facetItems []*searchdomain.FacetItem) []*CommerceSearchTreeFacetItem {
	items := make([]*CommerceSearchTreeFacetItem, len(facetItems))
	for i, item := range facetItems {
		items[i] = &CommerceSearchTreeFacetItem{
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
func WrapTreeFacet(facet searchdomain.Facet) *CommerceSearchTreeFacet {
	return &CommerceSearchTreeFacet{
		name:     facet.Name,
		label:    facet.Label,
		position: facet.Position,
		items:    mapTreeFacetItems(facet.Items),
	}
}

// CommerceSearchTreeFacet dto for tree facets
type CommerceSearchTreeFacet struct {
	name     string
	label    string
	position int
	items    []*CommerceSearchTreeFacetItem
}

// Name getter
func (c *CommerceSearchTreeFacet) Name() string {
	return c.name
}

// Label getter
func (c *CommerceSearchTreeFacet) Label() string {
	return c.label
}

// Position getter
func (c *CommerceSearchTreeFacet) Position() int {
	return c.position
}

// Items getter
func (c *CommerceSearchTreeFacet) Items() []*CommerceSearchTreeFacetItem {
	return c.items
}

func hasSelectedItem(items []*CommerceSearchTreeFacetItem) bool {
	for _, item := range items {
		if item.selected || hasSelectedItem(item.items) {
			return true
		}
	}
	return false
}

// HasSelectedItem getter
func (c *CommerceSearchTreeFacet) HasSelectedItem() bool {
	return hasSelectedItem(c.items)
}

// CommerceSearchTreeFacetItem dto for tree facet items
type CommerceSearchTreeFacetItem struct {
	label    string
	value    string
	selected bool
	count    int
	active   bool
	items    []*CommerceSearchTreeFacetItem
}

// Label getter
func (c *CommerceSearchTreeFacetItem) Label() string {
	return c.label
}

// Value getter
func (c *CommerceSearchTreeFacetItem) Value() string {
	return c.value
}

// Selected getter
func (c *CommerceSearchTreeFacetItem) Selected() bool {
	return c.selected
}

// Count getter
func (c *CommerceSearchTreeFacetItem) Count() int {
	return c.count
}

// Active getter
func (c *CommerceSearchTreeFacetItem) Active() bool {
	return c.active
}

// Items getter
func (c *CommerceSearchTreeFacetItem) Items() []*CommerceSearchTreeFacetItem {
	return c.items
}

// WrapRangeFacet wraps the range facet into the graphql dto
func WrapRangeFacet(facet searchdomain.Facet) *CommerceSearchRangeFacet {
	items := make([]*CommerceSearchRangeFacetItem, len(facet.Items))
	for i, item := range facet.Items {
		items[i] = &CommerceSearchRangeFacetItem{
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

	return &CommerceSearchRangeFacet{
		name:     facet.Name,
		label:    facet.Label,
		position: facet.Position,
		items:    items,
	}
}

// CommerceSearchRangeFacet dto for range facets
type CommerceSearchRangeFacet struct {
	name     string
	label    string
	position int
	items    []*CommerceSearchRangeFacetItem
}

// Name getter
func (c *CommerceSearchRangeFacet) Name() string {
	return c.name
}

// Label getter
func (c *CommerceSearchRangeFacet) Label() string {
	return c.label
}

// Position getter
func (c *CommerceSearchRangeFacet) Position() int {
	return c.position
}

// Items getter
func (c *CommerceSearchRangeFacet) Items() []*CommerceSearchRangeFacetItem {
	return c.items
}

// HasSelectedItem getter
func (c *CommerceSearchRangeFacet) HasSelectedItem() bool {
	for _, item := range c.items {
		if item.selected {
			return true
		}
	}
	return false
}

// CommerceSearchRangeFacetItem dto for range facet items
type CommerceSearchRangeFacetItem struct {
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
func (c *CommerceSearchRangeFacetItem) Label() string {
	return c.label
}

// Value getter
func (c *CommerceSearchRangeFacetItem) Value() string {
	return c.value
}

// Selected getter
func (c *CommerceSearchRangeFacetItem) Selected() bool {
	return c.selected
}

// Count getter
func (c *CommerceSearchRangeFacetItem) Count() int {
	return c.count
}

// Min getter
func (c *CommerceSearchRangeFacetItem) Min() int {
	return c.min
}

// Max getter
func (c *CommerceSearchRangeFacetItem) Max() int {
	return c.max
}

// SelectedMin getter
func (c *CommerceSearchRangeFacetItem) SelectedMin() int {
	return c.selectedMin
}

// SelectedMax getter
func (c *CommerceSearchRangeFacetItem) SelectedMax() int {
	return c.selectedMax
}
