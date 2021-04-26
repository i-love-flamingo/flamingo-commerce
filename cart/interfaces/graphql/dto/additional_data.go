package dto

type (

	// CustomAttributes represents map of custom attributes of cart and delivery info
	CustomAttributes struct {
		Attributes map[string]string
	}

	// CustomAttribute that represents a custom attribute of cart
	CustomAttribute interface {
		Key() string
		SetKey(key string)
	}

	// CustomAttributeBoolean that represents a boolean custom attribute of cart
	CustomAttributeBoolean struct {
		key   string
		Value bool
	}

	// CustomAttributeString that represents a string custom attribute of cart
	CustomAttributeString struct {
		key   string
		Value string
	}
)

// Get attribute by key
func (c *CustomAttributes) Get(key string) CustomAttribute {
	if c.Attributes == nil {
		return nil
	}

	if value, found := c.Attributes[key]; found {
		var attribute CustomAttribute

		switch value {
		case "true", "True":
			attribute = &CustomAttributeBoolean{Value: true}
		case "false", "False":
			attribute = &CustomAttributeBoolean{Value: false}
		default:
			attribute = &CustomAttributeString{Value: value}
		}

		attribute.SetKey(key)
		return attribute
	}

	return nil
}

// Key of custom attribute
func (c *CustomAttributeBoolean) Key() string {
	return c.key
}

// SetKey of custom attribute
func (c *CustomAttributeBoolean) SetKey(key string) {
	c.key = key
}

// Key of custom attribute
func (c *CustomAttributeString) Key() string {
	return c.key
}

// SetKey of custom attribute
func (c *CustomAttributeString) SetKey(key string) {
	c.key = key
}
