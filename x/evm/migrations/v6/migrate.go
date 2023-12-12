package v6

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
	v6types "github.com/evmos/ethermint/x/evm/migrations/v6/types"
	"github.com/evmos/ethermint/x/evm/types"
)

// MigrateStore migrates the x/evm module state from the consensus version 5 to
// version 6. Specifically, it takes the parameters that are currently stored
// and managed by the Cosmos SDK params module and stores them directly into the x/evm module state.
func MigrateStore(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) error {
	var (
		v6params v6types.Params
		params   types.Params
	)

	store := ctx.KVStore(storeKey)
	paramsBz := store.Get(types.KeyPrefixParams)
	cdc.MustUnmarshal(paramsBz, &v6params)

	params.EvmDenom = v6params.EvmDenom
	params.EnableCreate = v6params.EnableCreate
	params.EnableCall = v6params.EnableCall
	shanghaiTime := sdkmath.NewUint(0)
	params.ChainConfig = types.ChainConfig{
		HomesteadBlock:      v6params.ChainConfig.HomesteadBlock,
		DAOForkBlock:        v6params.ChainConfig.DAOForkBlock,
		DAOForkSupport:      v6params.ChainConfig.DAOForkSupport,
		EIP150Block:         v6params.ChainConfig.EIP150Block,
		EIP150Hash:          v6params.ChainConfig.EIP150Hash,
		EIP155Block:         v6params.ChainConfig.EIP155Block,
		EIP158Block:         v6params.ChainConfig.EIP158Block,
		ByzantiumBlock:      v6params.ChainConfig.ByzantiumBlock,
		ConstantinopleBlock: v6params.ChainConfig.ConstantinopleBlock,
		PetersburgBlock:     v6params.ChainConfig.PetersburgBlock,
		IstanbulBlock:       v6params.ChainConfig.IstanbulBlock,
		MuirGlacierBlock:    v6params.ChainConfig.MuirGlacierBlock,
		BerlinBlock:         v6params.ChainConfig.BerlinBlock,
		LondonBlock:         v6params.ChainConfig.LondonBlock,
		ArrowGlacierBlock:   v6params.ChainConfig.ArrowGlacierBlock,
		GrayGlacierBlock:    v6params.ChainConfig.GrayGlacierBlock,
		MergeNetsplitBlock:  v6params.ChainConfig.MergeNetsplitBlock,
		ShanghaiTime:        &shanghaiTime,
		CancunTime:          nil,
	}
	params.ExtraEIPs = v6params.ExtraEIPs
	params.AllowUnprotectedTxs = v6params.AllowUnprotectedTxs

	store.Delete(types.KeyPrefixParams)

	if err := params.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&params)

	store.Set(types.KeyPrefixParams, bz)
	return nil
}
