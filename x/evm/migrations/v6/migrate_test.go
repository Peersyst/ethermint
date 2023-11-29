package v6_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/evmos/ethermint/app"
	"github.com/evmos/ethermint/encoding"
	v6 "github.com/evmos/ethermint/x/evm/migrations/v6"
	v6types "github.com/evmos/ethermint/x/evm/migrations/v6/types"
	"github.com/evmos/ethermint/x/evm/types"
)

func TestMigrate(t *testing.T) {
	var oldParams v6types.Params
	encCfg := encoding.MakeConfig(app.ModuleBasics)
	cdc := encCfg.Codec

	storeKey := sdk.NewKVStoreKey(types.ModuleName)
	tKey := sdk.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)
	kvStore := ctx.KVStore(storeKey)

	chainConfig := types.DefaultChainConfig()

	// Set the params in the
	oldParams.EvmDenom = "aphoton"
	oldParams.ExtraEIPs = types.AvailableExtraEIPs
	oldParams.AllowUnprotectedTxs = true
	oldParams.EnableCall = true
	oldParams.EnableCreate = true
	oldParams.ChainConfig = v6types.ChainConfig{
		HomesteadBlock:      chainConfig.HomesteadBlock,
		DAOForkBlock:        chainConfig.DAOForkBlock,
		DAOForkSupport:      chainConfig.DAOForkSupport,
		EIP150Block:         chainConfig.EIP150Block,
		EIP150Hash:          chainConfig.EIP150Hash,
		EIP155Block:         chainConfig.EIP155Block,
		EIP158Block:         chainConfig.EIP158Block,
		ByzantiumBlock:      chainConfig.ByzantiumBlock,
		ConstantinopleBlock: chainConfig.ConstantinopleBlock,
		PetersburgBlock:     chainConfig.PetersburgBlock,
		IstanbulBlock:       chainConfig.IstanbulBlock,
		MuirGlacierBlock:    chainConfig.MuirGlacierBlock,
		BerlinBlock:         chainConfig.BerlinBlock,
		LondonBlock:         chainConfig.LondonBlock,
		ArrowGlacierBlock:   chainConfig.ArrowGlacierBlock,
		GrayGlacierBlock:    chainConfig.GrayGlacierBlock,
		MergeNetsplitBlock:  chainConfig.MergeNetsplitBlock,
		ShanghaiBlock:       nil,
		CancunBlock:         nil,
	}
	oldParamsBz := cdc.MustMarshal(&oldParams)
	kvStore.Set(types.KeyPrefixParams, oldParamsBz)

	err := v6.MigrateStore(ctx, storeKey, cdc)
	require.NoError(t, err)

	paramsBz := kvStore.Get(types.KeyPrefixParams)
	var params types.Params
	cdc.MustUnmarshal(paramsBz, &params)

	// test that the params have been migrated correctly
	require.Equal(t, oldParams.EvmDenom, params.EvmDenom)
	require.True(t, params.EnableCreate)
	require.True(t, params.EnableCall)
	require.True(t, params.AllowUnprotectedTxs)
	require.Equal(t, oldParams.ExtraEIPs, params.ExtraEIPs)
	chainConfig.ShanghaiTime = nil
	chainConfig.CancunTime = nil
	require.Equal(t, chainConfig, params.ChainConfig)
}
