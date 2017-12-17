package script

import (
	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/client"
	"github.com/cockroachdb/cockroach/pkg/rpc"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/stop"
	"github.com/cockroachdb/cockroach/pkg/util/tracing"
)

// NewDB creates a client to the CockroachDB at <address>.
func NewDB(address string) (*client.DB, error) {
	clock := hlc.NewClock(hlc.UnixNano, 0)
	stopper := stop.NewStopper()
	rpcContext := rpc.NewContext(log.AmbientContext{Tracer: tracing.NewTracer()}, &base.Config{
		User:     security.NodeUser,
		Insecure: true,
	}, clock, stopper)
	conn, err := rpcContext.GRPCDial(address)
	if err != nil {
		return nil, err
	}
	db := client.NewDB(client.NewSender(conn), clock)
	return db, nil
}
