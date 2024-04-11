package loader

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/plugins/log"
	"github.com/grafana/grafana/pkg/plugins/manager/pipeline/bootstrap"
	"github.com/grafana/grafana/pkg/plugins/manager/pipeline/discovery"
	"github.com/grafana/grafana/pkg/plugins/manager/pipeline/initialization"
	"github.com/grafana/grafana/pkg/plugins/manager/pipeline/termination"
	"github.com/grafana/grafana/pkg/plugins/manager/pipeline/validation"
)

type Loader struct {
	discovery   discovery.Discoverer
	bootstrap   bootstrap.Bootstrapper
	initializer initialization.Initializer
	termination termination.Terminator
	validation  validation.Validator
	log         log.Logger
}

func New(discovery discovery.Discoverer, bootstrap bootstrap.Bootstrapper, validation validation.Validator,
	initializer initialization.Initializer, termination termination.Terminator) *Loader {
	return &Loader{
		discovery:   discovery,
		bootstrap:   bootstrap,
		validation:  validation,
		initializer: initializer,
		termination: termination,
		log:         log.New("plugin.loader"),
	}
}

func (l *Loader) Load(ctx context.Context, src plugins.PluginSource) ([]*plugins.Plugin, error) {
	end := l.instrumentLoad(ctx, src)

	discoveredPlugins, err := l.discovery.Discover(ctx, src)
	if err != nil {
		return nil, err
	}

	bootstrappedPlugins := []*plugins.Plugin{}
	for _, foundBundle := range discoveredPlugins {
		bootstrappedPlugin, err := l.bootstrap.Bootstrap(ctx, src, foundBundle)
		if err != nil {
			// TODO: Add error to registry
			continue
		}
		bootstrappedPlugins = append(bootstrappedPlugins, bootstrappedPlugin...)
	}

	validatedPlugins := []*plugins.Plugin{}
	for _, bootstrappedPlugin := range bootstrappedPlugins {
		err := l.validation.Validate(ctx, bootstrappedPlugin)
		if err != nil {
			// TODO: Add error to registry
			continue
		}
		validatedPlugins = append(validatedPlugins, bootstrappedPlugin)
	}

	initializedPlugins := []*plugins.Plugin{}
	for _, validatedPlugin := range validatedPlugins {
		initializedPlugin, err := l.initializer.Initialize(ctx, validatedPlugin)
		if err != nil {
			// TODO: Add error to registry
			continue
		}
		initializedPlugins = append(initializedPlugins, initializedPlugin)
	}

	end(initializedPlugins)

	return initializedPlugins, nil
}

func (l *Loader) Unload(ctx context.Context, p *plugins.Plugin) (*plugins.Plugin, error) {
	return l.termination.Terminate(ctx, p)
}

func (l *Loader) instrumentLoad(ctx context.Context, src plugins.PluginSource) func([]*plugins.Plugin) {
	start := time.Now()
	sourceLogger := l.log.New("source", src.PluginClass(ctx)).FromContext(ctx)
	sourceLogger.Debug("Loading plugin source...")

	return func(logger log.Logger, start time.Time) func([]*plugins.Plugin) {
		return func(plugins []*plugins.Plugin) {
			if len(plugins) == 0 {
				logger.Debug("Plugin source loaded, though no plugins were found")
				return
			}
			names := make([]string, len(plugins))
			for i, p := range plugins {
				names[i] = p.ID
			}
			sort.Strings(names)
			pluginsStr := strings.Join(names, ", ")

			logger.Debug("Plugin source loaded", "plugins", pluginsStr, "duration", time.Since(start))
		}
	}(sourceLogger, start)
}
