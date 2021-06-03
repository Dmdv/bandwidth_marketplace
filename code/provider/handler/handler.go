package handler

import (
	"context"
	"net/http"

	"github.com/0chain/gosdk/zcncore"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/build"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/chain"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/datastore"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/log"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/node"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// SetupHandlers sets up the necessary API end points.
func SetupHandlers(r *mux.Router) {
	r.HandleFunc("/", HomePageHandler)
}

// HomePageHandler provides basic info when accessing the home page of the server.
func HomePageHandler(w http.ResponseWriter, _ *http.Request) {
	selfNode := node.GetSelfNode()

	mc := chain.GetServerChain()
	page := "<div>Running since " + selfNode.StartTime().String() + " ...\n" +
		"<div>Working on the chain: " + mc.ID + "</div>\n" +
		"<div>I am a Magma Provider with<ul>\n" +
		"<li>id: " + selfNode.ID() + "</li>\n" +
		"<li>public_key: " + selfNode.PublicKey() + "</li>\n" +
		"<li>build_tag: " + build.Tag + "</li>\n" +
		"</ul></div>\n" +
		"<div>Miners ...\n"

	network := zcncore.GetNetwork()
	for _, miner := range network.Miners {
		page += miner + "\n"
	}
	page += "</div><div>Sharders ...\n"
	for _, sharder := range network.Sharders {
		page += sharder + "\n"
	}
	page += "</div"

	if _, err := w.Write([]byte(page)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//nolint:gosimple // need more time to verify
func HandleShutdown(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			log.Logger.Info("Closing database")
			if err := datastore.GetStore().Close(); err != nil {
				log.Logger.Error("Error while closing database: ", zap.Error(err))
			}
		}
	}()
}
