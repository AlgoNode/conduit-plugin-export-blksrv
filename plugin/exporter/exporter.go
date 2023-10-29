package exporter

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/algorand/conduit/conduit/data"
	"github.com/algorand/conduit/conduit/plugins"
	"github.com/algorand/conduit/conduit/plugins/exporters"
)

//go:embed sample.yaml
var sampleConfig string

// metadata contains information about the plugin used for CLI helpers.
var metadata = plugins.Metadata{
	Name:         "delta_exporter",
	Description:  "Exports block deltas into custom Nodely block server.",
	Deprecated:   false,
	SampleConfig: sampleConfig,
}

func init() {
	exporters.Register(metadata.Name, exporters.ExporterConstructorFunc(func() exporters.Exporter {
		return &deltaExporter{}
	}))
}

// deltaExporter is the object which implements the exporter plugin interface.
type deltaExporter struct {
	log *logrus.Logger
	cfg Config
	ip  data.InitProvider
	ht  *http.Client
	ctx context.Context
}

func (oe *deltaExporter) Metadata() plugins.Metadata {
	return metadata
}

func (oe *deltaExporter) Config() string {
	ret, _ := yaml.Marshal(oe.cfg)
	return string(ret)
}

func (oe *deltaExporter) Close() error {
	oe.log.Infof("Shutting down")
	return nil
}

func (oe *deltaExporter) Init(ctx context.Context, ip data.InitProvider, cfg plugins.PluginConfig, logger *logrus.Logger) error {
	var err error
	oe.log = logger
	oe.ctx = ctx
	oe.ip = ip
	if err := cfg.UnmarshalConfig(&oe.cfg); err != nil {
		return fmt.Errorf("unable to read configuration: %w", err)
	}

	if oe.ht, err = oe.blksrvInit(); err != nil {
		return err
	}

	oe.setGenesis(ip.GetGenesis())

	return nil
}

func (oe *deltaExporter) Receive(exportData data.BlockData) error {
	round := exportData.BlockHeader.Round

	oe.log.Infof("Processing block %d ", round)

	if err := oe.export(exportData); err != nil {
		return err
	}

	return nil
}
