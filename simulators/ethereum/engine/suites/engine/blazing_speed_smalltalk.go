package suite_engine

import (
	"time"

	"github.com/ethereum/hive/simulators/ethereum/engine/client/node"
	"github.com/ethereum/hive/simulators/ethereum/engine/clmock"
	"github.com/ethereum/hive/simulators/ethereum/engine/config"
	"github.com/ethereum/hive/simulators/ethereum/engine/test"
)

type BlazingSpeed struct {
	test.BaseSpec
	ProduceBlocksBeforePeering int
}

func (s BlazingSpeed) WithMainFork(fork config.Fork) test.Spec {
	specCopy := s
	specCopy.MainFork = fork
	return specCopy
}

func (ft BlazingSpeed) GetName() string {
	return "Blazing Speed Test"
}

func (s BlazingSpeed) GetForkConfig() *config.ForkConfig {
	forkConfig := s.BaseSpec.GetForkConfig()
	if forkConfig == nil {
		return nil
	}
	return forkConfig
}

func (ft BlazingSpeed) Execute(t *test.Env) {

	// Remove the original client so that it does not receive the payloads created on the canonical chain
	t.CLMock.RemoveEngineClient(t.Engine)

	// To allow having the invalid payload delivered via P2P, we need a second client to serve the payload
	starter := node.GethNodeEngineStarter{
		Config: node.GethNodeTestConfiguration{},
	}

	secondaryClient, err := starter.StartGethNode(t.T, t.TestContext, t.Genesis, t.ClientParams, t.ClientFiles)
	if err != nil {
		t.Fatalf("FAIL (%s): Unable to spawn a secondary client: %v", t.TestName, err)
	}
	t.CLMock.AddEngineClient(secondaryClient)
	test.NewTestEngineClient(t, secondaryClient)

	// Produce blocks before starting the test if required
	t.CLMock.ProduceBlocks(20, clmock.BlockProcessCallbacks{})

	// Add the main client as a peer of the secondary client so it is able to sync
	secondaryClient.AddPeer(t.Engine)
	// Add back the original client before side chain production
	t.CLMock.AddEngineClient(t.Engine)

	timeout := time.Now().Add(time.Second * 60)
	start := time.Now()

	is_synced := false
	for !is_synced {
		time.Sleep(time.Microsecond * 50)
		result := t.CLMock.BroadcastForkchoiceUpdated(&t.CLMock.LatestForkchoice, nil, 1)
		if len(result) > 0 {
			latest_result := result[len(result)-1]
			forkchoice_response := latest_result.ForkchoiceResponse
			is_synced = forkchoice_response.PayloadStatus.Status != "SYNCING"
		}

		if time.Now().After(timeout) {
			t.Fatalf("FAIL (%s): Timed out: %v")
		}
	}

	t.Logf("Info: We have finished start: %v, end: %v", start, time.Now())
}
