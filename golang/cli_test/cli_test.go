package clitest

import (
	"testing"
	"fmt"
	"encoding/json"
	"encoding/hex"

	"github.com/stretchr/testify/require"
	rtypes "github.com/summa-tx/relays/golang/x/relay/types"
)


func TestRelayCLIIsAncestor(t *testing.T) {
	// Get data needed for transaction
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)

	var newDifficultyHeaders []rtypes.BitcoinHeader
	newDiffJSON := readJSONFile(t, "0_new_difficulty")
	err = json.Unmarshal([]byte(newDiffJSON), &newDifficultyHeaders)
	require.NoError(t, err)

	// Query Chain for Actual Value
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	// Set transaction parameter values
	fooAddr := f.KeyAddress(keyFoo)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	ancestor := hex.EncodeToString(newDifficultyHeaders[0].HashLE[:])
	digest := hex.EncodeToString(newDifficultyHeaders[1].HashLE[:])
	limit := "5"

	// must ingest headers in order to perform query
	success, _, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	require.True(t, success, stderr)

	isancestor := f.QueryIsAncestor(digest, ancestor, limit)

	// Condition
	expected := true
	actual := isancestor.Res
	require.Equal(t, expected, actual)

	//Cleanup
	f.Cleanup()
}

func TestRelayCLIGetRelayGenesis(t *testing.T) {
	// Get Expected Value
	fmt.Println("tired of commenting and uncommenting fmt")
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)
	expected := genesisHeaders[1].HashLE


	// Query Chain for Actual Value
    f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	genesisRelay := f.QueryGetRelayGenesis(fooAddr)
	actual := genesisRelay.Res

	// Condition
	require.Equal(t, expected, actual)

	//Cleanup
	f.Cleanup()
}

func TestRelayCLIGetLastReorgLCA(t *testing.T) {
	// Get Expected Value
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)
	expected := genesisHeaders[1].HashLE

	// Query Chain for Actual Value
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	lastReorgLCA := f.QueryGetLastReorgLCA(fooAddr)
	actual := lastReorgLCA.Res

	// Condition
	require.Equal(t, expected, actual)

	//Cleanup
	f.Cleanup()
}

func TestRelayCLIGetBestDigest(t *testing.T) {
	// Get Expected Value
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)
	expected := genesisHeaders[1].HashLE

	// Query Chain for Actual Value
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	bestDigest := f.QueryGetBestDigest(fooAddr)
	actual := bestDigest.Res

	// Condition
	require.Equal(t, expected, actual)

	//Cleanup
	f.Cleanup()
}

func TestRelayCLIQueryFindAncestor(t *testing.T) {
	// Extract data for transactions
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)

	var newDifficultyHeaders []rtypes.BitcoinHeader
	newDiffJSON := readJSONFile(t, "0_new_difficulty")
	err = json.Unmarshal([]byte(newDiffJSON), &newDifficultyHeaders)
	require.NoError(t, err)

	// Transact with Chain for Actual Value
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])

	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	// Require findancestor fails if ancestor does not exist on relay
	digest := hex.EncodeToString(newDifficultyHeaders[1].HashLE[:])
	offset := "5"
	f.QueryFindAncestorInvalid("could not find ancestor", digest, offset)

	// Require findancestor returns ancestor if valid query
	offset = "1"
	findancestor := f.QueryFindAncestor(digest, offset)
	expected := hex.EncodeToString(newDifficultyHeaders[0].HashLE[:])
	actual := hex.EncodeToString(findancestor.Res[:])
	require.Equal(t, expected, actual)
}

// func TestRelayCLIIsMostRecentCommonAncestor(t *testing.T) {
// 	// Get data needed for transaction
// 	var genesisHeaders []rtypes.BitcoinHeader
// 	genesisJSON := readJSONFile(t, "genesis")
// 	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
// 	require.NoError(t, err)
//
// 	var newDifficultyHeaders []rtypes.BitcoinHeader
// 	newDiffJSON := readJSONFile(t, "0_new_difficulty")
// 	err = json.Unmarshal([]byte(newDiffJSON), &newDifficultyHeaders)
// 	require.NoError(t, err)
//
// 	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])
//
// 	// Query Chain for Actual Value
// 	f := InitFixtures(t)
// 	proc := f.RelayDStart()
// 	defer proc.Stop(false)
// 	fooAddr := f.KeyAddress(keyFoo)
// 	// must ingest headers in order to perform query
// 	success, _, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
// 	require.True(t, success, stderr)
//
// 	//perform query
// 	ancestor := hex.EncodeToString(newDifficultyHeaders[0].HashLE[:])
// 	left := hex.EncodeToString(newDifficultyHeaders[0].HashLE[:])
// 	right := hex.EncodeToString(newDifficultyHeaders[1].HashLE[:])
// 	limit := "3"
// 	ismostrecentcommonancestor := f.QueryIsMostRecentCommonAncestor(ancestor, left, right, limit)
// 	actual := ismostrecentcommonancestor.Res
//
// 	// Condition
// 	expected := true
// 	require.Equal(t, expected, actual)
//
// 	// Require query returns error if invalid
// 	ancestor = hex.EncodeToString(newDifficultyHeaders[1].HashLE[:])
// 	ismostrecentcommonancestor = f.QueryIsMostRecentCommonAncestor(ancestor, left, right, limit)
// 	fmt.Println(ismostrecentcommonancestor)
//
// 	//Cleanup
// 	f.Cleanup()
// }

func TestRelayCLIQueryHeaviestFromAncestor(t *testing.T) {
	// Extract data for transactions
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)

	var newDifficultyHeaders []rtypes.BitcoinHeader
	newDiffJSON := readJSONFile(t, "0_new_difficulty")
	err = json.Unmarshal([]byte(newDiffJSON), &newDifficultyHeaders)
	require.NoError(t, err)

	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])

	// Transact with Chain for Actual Value
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)

	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	// Require heaviestfromancestor returns valid newBest
	ancestor := hex.EncodeToString(genesisHeaders[1].HashLE[:])
	currentBest := hex.EncodeToString(genesisHeaders[1].HashLE[:])
	newBest := hex.EncodeToString(newDifficultyHeaders[1].HashLE[:])
	limit := "10"
	heaviestfromancestor := f.QueryHeaviestFromAncestor(ancestor, currentBest, newBest, limit)
	expected := newBest
	actual := hex.EncodeToString(heaviestfromancestor.Res[:])
	require.Equal(t, expected, actual)

	// Require heaviestfromancestor fails with invalid params
	invalidNewBest := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	f.QueryHeaviestFromAncestorInvalid("could not determine if", ancestor, currentBest, invalidNewBest, limit)
}

func TestRelayCLIQueryCheckProof(t *testing.T) {
	// Extract data for transactions
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)

	var newDifficultyHeaders []rtypes.BitcoinHeader
	newDiffJSON := readJSONFile(t, "0_new_difficulty")
	err = json.Unmarshal([]byte(newDiffJSON), &newDifficultyHeaders)
	require.NoError(t, err)

	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])

	// Transact with Chain for Actual Value
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)

	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	// Require checkproof fails without associated headers with
	checkProof := f.QueryCheckProof("1_check_proof.json", "--inputfile")
	require.Equal(t, false, checkProof.Valid)

	success, stdout, stderr = f.TxIngestHeaders(fooAddr, "2_ingest_headers.json", "--inputfile -y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	// Require proof is valid when associated header exists with valid transaction
	checkProof = f.QueryCheckProof("1_check_proof.json", "--inputfile")
	require.Equal(t, true, checkProof.Valid)
}

func TestRelayCLITXIngestHeaders(t *testing.T) {
	// Extract data for transaction
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)

	var newDifficultyHeaders []rtypes.BitcoinHeader
	newDiffJSON := readJSONFile(t, "0_new_difficulty")
	err = json.Unmarshal([]byte(newDiffJSON), &newDifficultyHeaders)
	require.NoError(t, err)

	var newHeaders []rtypes.BitcoinHeader
	ingestHeadersJSON := readJSONFile(t, "2_ingest_headers")
	err = json.Unmarshal([]byte(ingestHeadersJSON), &newHeaders)
	require.NoError(t, err)

	// Transact with Chain
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)

	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	success, stdout, stderr = f.TxIngestHeaders(fooAddr, "2_ingest_headers.json", "--inputfile -y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	//Cleanup
	f.Cleanup()
}

func TestRelayCLITXIngestDiffChange(t *testing.T) {
	// Extract data for transaction
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)

	var newDifficultyHeaders []rtypes.BitcoinHeader
	newDiffJSON := readJSONFile(t, "0_new_difficulty")
	err = json.Unmarshal([]byte(newDiffJSON), &newDifficultyHeaders)
	require.NoError(t, err)

	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])

	// Transact with Chain
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")

	// require successful transaction
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	//Cleanup
	f.Cleanup()
}

func TestRelayCLITXProvideProof(t *testing.T) {
	// Extract data for transactions
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)

	var newDifficultyHeaders []rtypes.BitcoinHeader
	newDiffJSON := readJSONFile(t, "0_new_difficulty")
	err = json.Unmarshal([]byte(newDiffJSON), &newDifficultyHeaders)
	require.NoError(t, err)

	// Transact with Chain for Actual Value
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])

	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	success, stdout, stderr = f.TxIngestHeaders(fooAddr, "2_ingest_headers.json", "--inputfile -y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	// checkproof fails given invalid requests
	success, stdout, stderr = f.TxProvideProof(fooAddr, "1_check_proof.json", "3_filled_requests.json", "--inputfile -y")
	require.Contains(t, stdout, `"Request not found`)

	// submit request
	spends   := "0x"
	pays     := "0x17a91423737cd98bb6b2da5a11bcd82e5de36591d69f9f87"
	value    := "0"
	numConfs := "1"
	success, stdout, stderr = f.TxNewRequest(fooAddr, spends, pays, value, numConfs, "-y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	// checkproof succeeds given valid proof and requests
	success, stdout, stderr = f.TxProvideProof(fooAddr, "1_check_proof.json", "3_filled_requests.json", "--inputfile -y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)
}

func TestRelayCLITxMarkNewHeaviest(t *testing.T) {
	// Extract data for transactions
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)

	var newDifficultyHeaders []rtypes.BitcoinHeader
	newDiffJSON := readJSONFile(t, "0_new_difficulty")
	err = json.Unmarshal([]byte(newDiffJSON), &newDifficultyHeaders)
	require.NoError(t, err)

	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	ancestor := hex.EncodeToString(genesisHeaders[1].HashLE[:])
	bestKnown := hex.EncodeToString(genesisHeaders[1].Raw[:])
	newBest := hex.EncodeToString(newDifficultyHeaders[1].Raw[:])
	limit := "10"

	// Get Expected Value
	expected := hex.EncodeToString(newDifficultyHeaders[1].HashLE[:])

	// Transact with Chain for Actual Value
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)

	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	success, stdout, stderr = f.TxMarkNewHeaviest(fooAddr, ancestor, bestKnown, newBest, limit, "-y")
	require.True(t, success, stderr)
	require.Contains(t, stdout, `"success":true`)

	bestDigest := f.QueryGetBestDigest(fooAddr)
	actual := hex.EncodeToString(bestDigest.Res[:])

	// Condition
	require.Equal(t, expected, actual)

	//Cleanup
	f.Cleanup()
}