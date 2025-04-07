package libnetwork

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/docker/docker/libnetwork/config"
	"gotest.tools/v3/assert"
)

func TestEndpointStore(t *testing.T) {
	configOption := config.OptionDataDir(t.TempDir())
	c, err := New(configOption)
	assert.NilError(t, err)
	defer c.Stop()

	// Insert a first endpoint
	nw := &Network{id: "testNetwork"}
	ep1 := &Endpoint{network: nw, id: "testEndpoint1"}
	err = c.storeEndpoint(context.Background(), ep1)
	assert.NilError(t, err)

	// Then a second endpoint
	ep2 := &Endpoint{network: nw, id: "testEndpoint2"}
	err = c.storeEndpoint(context.Background(), ep2)
	assert.NilError(t, err)

	// Check that we can find both endpoints
	found := c.findEndpoints(filterEndpointByNetworkId("testNetwork"))
	slices.SortFunc(found, func(a, b *Endpoint) int { return strings.Compare(a.id, b.id) })
	assert.Equal(t, len(found), 2)
	assert.Equal(t, found[0], ep1)
	assert.Equal(t, found[1], ep2)

	// Delete the first endpoint
	err = c.deleteStoredEndpoint(ep1)
	assert.NilError(t, err)

	// Check that we can only find the second endpoint
	found = c.findEndpoints(filterEndpointByNetworkId("testNetwork"))
	assert.Equal(t, len(found), 1)
	assert.Equal(t, found[0], ep2)

	// Store the second endpoint again
	err = c.storeEndpoint(context.Background(), ep2)
	assert.NilError(t, err)
}
