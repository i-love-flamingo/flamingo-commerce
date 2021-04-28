package dto

type (

	// CustomAttributes represents map of custom attributes of cart and delivery info
	CustomAttributes struct {
		Attributes map[string]string
	}

	// KeyValue for cart and delivery
	KeyValue struct {
		Key   string
		Value string
	}

	// DeliveryAdditionalData of delivery
	DeliveryAdditionalData struct {
		DeliveryCode   string
		AdditionalData []KeyValue
	}
)

// Get attribute by key
func (c *CustomAttributes) Get(key string) *KeyValue {
	if c.Attributes == nil {
		return nil
	}

	if value, found := c.Attributes[key]; found {
		return &KeyValue{Key: key, Value: value}
	}

	return nil
}
