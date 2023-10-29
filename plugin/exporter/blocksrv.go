package exporter

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/algorand/conduit/conduit/data"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
)

func (oe *deltaExporter) blksrvInit() (*http.Client, error) {

	ht := http.DefaultTransport.(*http.Transport).Clone()
	ht.MaxConnsPerHost = 100
	ht.MaxIdleConns = 100
	ht.MaxIdleConnsPerHost = 100

	return &http.Client{
		Timeout:   time.Second * 5,
		Transport: ht,
	}, nil
}

func (oe *deltaExporter) export(exportData data.BlockData) error {
	round := exportData.BlockHeader.Round
	buf := msgpack.Encode(exportData)

	url := fmt.Sprintf("%s/n2/conduit/blockdata/%d", oe.cfg.blocksrv, int(round))

	req, err := http.NewRequestWithContext(oe.ctx, http.MethodPut, url, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/msgpack")
	resp, err := oe.ht.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	oe.log.WithField("round", strconv.Itoa(int(round))).Infof("Block %dB exported with code %d", len(buf), resp.StatusCode)
	return nil
}
