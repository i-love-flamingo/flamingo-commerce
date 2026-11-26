package domain_test

import (
	"errors"
	"testing"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

func TestConfigurableProduct_Variant_NotFound(t *testing.T) {
	t.Parallel()

	_, err := domain.ConfigurableProduct{}.Variant("missing")

	if !errors.Is(err, domain.ErrVariantNotFound) {
		t.Fatalf("expected ErrVariantNotFound, got %v", err)
	}

	if want := "No Variant with code missing found "; err.Error() != want {
		t.Errorf("message changed: got %q, want %q", err.Error(), want)
	}
}
