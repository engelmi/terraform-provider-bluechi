package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	blueChiProvider "github.com/engelmi/terraform-provider-bluechi/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"bluechi": providerserver.NewProtocol6WithError(blueChiProvider.New("1.0.0")()),
}

func testAccPreCheck(t *testing.T) {}
