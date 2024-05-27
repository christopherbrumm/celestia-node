package getters

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KYVENetwork/trustless-client-go/proof"
	"github.com/KYVENetwork/trustless-client-go/types"
	"github.com/celestiaorg/rsmt2d"
	"io/ioutil"
	"net/http"

	"github.com/celestiaorg/celestia-node/header"
	"github.com/celestiaorg/celestia-node/share"
)

// KYVEGetter implements custom share.Getter that uses the KYVE data archive to provide
// historical data.
type KYVEGetter struct {
	nodeEndpoint         string
	trustlessApiEndpoint string
}

// NewKYVEGetter instantiates a new KYVEGetter with the required Trustless API endpoint.
func NewKYVEGetter(trustlessApiEndpoint, nodeEndpoint string) *KYVEGetter {
	return &KYVEGetter{
		nodeEndpoint:         nodeEndpoint,
		trustlessApiEndpoint: trustlessApiEndpoint,
	}
}

// GetShare gets a share from any of registered share.Getters in cascading order.
func (kg *KYVEGetter) GetShare(
	ctx context.Context, header *header.ExtendedHeader, row, col int,
) (share.Share, error) {
	// TODO
	return nil, fmt.Errorf("not implemented yet")
}

// GetEDS gets a full EDS from any of registered share.Getters in cascading order.
func (kg *KYVEGetter) GetEDS(
	ctx context.Context, header *header.ExtendedHeader,
) (*rsmt2d.ExtendedDataSquare, error) {
	// TODO
	return nil, fmt.Errorf("not implemented yet")
}

// GetSharesByNamespace gets NamespacedShares from the provided KYVE Trustless API endpoint.
func (kg *KYVEGetter) GetSharesByNamespace(
	_ context.Context,
	header *header.ExtendedHeader,
	_ share.Namespace,
) (share.NamespacedShares, error) {
	// TODO: Replace last value with namespace object
	resp, err := http.Get(fmt.Sprintf("%vcelestia/GetSharesByNamespace?height=%v&namespace=%v", kg.trustlessApiEndpoint, header.Height(), "AAAAAAAAAAAAAAAAAAAAAAAAAIZiad33fbxA7Z0="))
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to request KYVE endpoint: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("got status code %d != 200", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var trustlessDataItem types.TrustlessDataItem
	if err = json.Unmarshal(data, &trustlessDataItem); err != nil {
		return nil, fmt.Errorf("failed to unmarshal KYVE Trustless API response: %v", err)
	}

	err = proof.DataItemInclusionProof(trustlessDataItem, kg.nodeEndpoint)
	if err != nil {
		// TODO: Improve error handling with typed out errors (e.g. namespace not supported, height not supported, DataItemInclusionProof failed)
		return nil, err
	}

	var sharesObject types.NamespacedShares
	if err = json.Unmarshal(trustlessDataItem.Value, &sharesObject); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sharesObject from KYVE Trustless API response: %w", err)
	}

	var shares share.NamespacedShares
	if err = json.Unmarshal(sharesObject.Data, &shares); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shares from KYVE Trustless API response: %w", err)
	}

	return shares, nil
}
