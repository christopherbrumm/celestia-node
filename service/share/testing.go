package share

import (
	"context"
	"math"
	"testing"

	mdutils "github.com/ipfs/go-merkledag/test"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/pkg/wrapper"

	"github.com/celestiaorg/nmt"
	"github.com/celestiaorg/rsmt2d"

	"github.com/celestiaorg/celestia-node/ipld"
	"github.com/celestiaorg/celestia-node/service/header"
)

// RandServiceWithSquare provides a share.Service filled with 'n' NMT trees of 'n' random shares, essentially storing a
// whole square.
func RandServiceWithSquare(t *testing.T, n int) (Service, *Root) {
	shares := RandShares(t, n*n)
	sharesSlices := make([][]byte, n*n)
	for i, share := range shares {
		sharesSlices[i] = share
	}
	dag, ctx := mdutils.Mock(), context.Background()
	na := ipld.NewNmtNodeAdder(ctx, dag)

	squareSize := uint32(math.Sqrt(float64(len(shares))))
	tree := wrapper.NewErasuredNamespacedMerkleTree(uint64(squareSize), nmt.NodeVisitor(na.Visit))
	eds, err := rsmt2d.ComputeExtendedDataSquare(sharesSlices, rsmt2d.NewRSGF8Codec(), tree.Constructor)
	require.NoError(t, err)

	err = na.Commit()
	require.NoError(t, err)

	dah, err := header.DataAvailabilityHeaderFromExtendedData(eds)
	require.NoError(t, err)

	return NewService(dag, NewLightAvailability(dag)), &dah
}

// RandShares provides 'n' randomized shares prefixed with random namespaces.
func RandShares(t *testing.T, n int) []Share {
	shares := make([]Share, n)
	for i, share := range ipld.RandNamespacedShares(t, n) {
		shares[i] = Share(share.Share)
	}
	return shares
}