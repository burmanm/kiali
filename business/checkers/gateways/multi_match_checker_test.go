package gateways

import (
	"testing"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/tests/data"
	"github.com/stretchr/testify/assert"
)

func TestCorrectGateways(t *testing.T) {
	conf := config.NewConfig()
	config.Set(conf)

	assert := assert.New(t)

	gwObject := data.AddServerToGateway(data.CreateServer([]string{"valid"}, 80, "http", "http"),
		data.CreateEmptyGateway("validgateway", map[string]string{
			"app": "real",
		}))

	gws := [][]kubernetes.IstioObject{[]kubernetes.IstioObject{gwObject}}

	validations := MultiMatchChecker{
		GatewaysPerNamespace: gws,
	}.Check()

	assert.Empty(validations)
	_, ok := validations[models.IstioValidationKey{"gateway", "validgateway"}]
	assert.False(ok)
}

func TestSameHostPortConfigInDifferentNamespace(t *testing.T) {
	conf := config.NewConfig()
	config.Set(conf)

	assert := assert.New(t)

	gwObject := data.AddServerToGateway(data.CreateServer([]string{"valid"}, 80, "http", "http"),
		data.CreateEmptyGateway("validgateway", map[string]string{
			"app": "real",
		}))

	// Another namespace
	gwObject2 := data.AddServerToGateway(data.CreateServer([]string{"valid"}, 80, "http", "http"),
		data.CreateEmptyGateway("stillvalid", map[string]string{
			"app": "someother",
		}))

	gws := [][]kubernetes.IstioObject{[]kubernetes.IstioObject{gwObject}, []kubernetes.IstioObject{gwObject2}}

	validations := MultiMatchChecker{
		GatewaysPerNamespace: gws,
	}.Check()

	assert.NotEmpty(validations)
	assert.Equal(1, len(validations))
	validation, ok := validations[models.IstioValidationKey{"gateway", "stillvalid"}]
	assert.True(ok)
	assert.False(validation.Valid)
}

func TestWildCardMatchingHost(t *testing.T) {
	conf := config.NewConfig()
	config.Set(conf)

	assert := assert.New(t)

	gwObject := data.AddServerToGateway(data.CreateServer([]string{"valid"}, 80, "http", "http"),
		data.CreateEmptyGateway("validgateway", map[string]string{
			"app": "real",
		}))

	// Another namespace
	gwObject2 := data.AddServerToGateway(data.CreateServer([]string{"*"}, 80, "http", "http"),
		data.CreateEmptyGateway("stillvalid", map[string]string{
			"app": "someother",
		}))

	gws := [][]kubernetes.IstioObject{[]kubernetes.IstioObject{gwObject}, []kubernetes.IstioObject{gwObject2}}

	validations := MultiMatchChecker{
		GatewaysPerNamespace: gws,
	}.Check()

	assert.NotEmpty(validations)
	assert.Equal(1, len(validations))
	validation, ok := validations[models.IstioValidationKey{"gateway", "stillvalid"}]
	assert.True(ok)
	assert.False(validation.Valid)

}

func TestTwoWildCardsMatching(t *testing.T) {
	conf := config.NewConfig()
	config.Set(conf)

	assert := assert.New(t)

	gwObject := data.AddServerToGateway(data.CreateServer([]string{"*"}, 80, "http", "http"),
		data.CreateEmptyGateway("validgateway", map[string]string{
			"app": "real",
		}))

	// Another namespace
	gwObject2 := data.AddServerToGateway(data.CreateServer([]string{"*"}, 80, "http", "http"),
		data.CreateEmptyGateway("stillvalid", map[string]string{
			"app": "someother",
		}))

	gws := [][]kubernetes.IstioObject{[]kubernetes.IstioObject{gwObject}, []kubernetes.IstioObject{gwObject2}}

	validations := MultiMatchChecker{
		GatewaysPerNamespace: gws,
	}.Check()

	assert.NotEmpty(validations)
	assert.Equal(1, len(validations))
	validation, ok := validations[models.IstioValidationKey{"gateway", "stillvalid"}]
	assert.True(ok)
	assert.False(validation.Valid)
}
