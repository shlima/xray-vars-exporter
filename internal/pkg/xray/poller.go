package xray

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"xray-vars-exporter/internal/pkg/metrics"

	"github.com/samber/lo"
)

type Poller struct {
	client                       IClient
	stdout                       io.Writer
	stderr                       io.Writer
	requestTimeout               time.Duration
	overrideOfflineDelay         time.Duration
	metrics                      metrics.IMetrics
	previousObservatoryOutbounds []string
	mx                           sync.Mutex
}

func NewPoller(client IClient, metrics metrics.IMetrics) *Poller {
	return &Poller{
		client:  client,
		stdout:  os.Stdout,
		stderr:  os.Stderr,
		metrics: metrics,
	}
}

func (p *Poller) SetRequestTimeoput(input time.Duration) {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.requestTimeout = input
}

func (p *Poller) SetOverrideOfflineDelay(input time.Duration) {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.overrideOfflineDelay = input
}

func (p *Poller) AsyncStartPolling(ctx context.Context, interval time.Duration) {
	go p.SyncStartPolling(ctx, interval)
}

func (p *Poller) SyncStartPolling(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.poll(ctx); err != nil {
				_, _ = p.stderr.Write([]byte(time.Now().String() + " " + err.Error() + "\n"))
			}
		}
	}
}

func (p *Poller) poll(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, p.requestTimeout)
	defer cancel()

	got, err := p.client.GetDebugVars(ctx)
	if err != nil {
		p.metrics.IncScrapeErrors()
		return fmt.Errorf("failed to client.GetDebugVars: %w", err)
	}

	current := lo.Keys(got.Observatories)
	p.cleanPreviousObservatories(current)

	for outbound, observatory := range got.Observatories {
		p.gaugeObservatory(outbound, observatory)
	}

	return nil
}

func (p *Poller) cleanPreviousObservatories(current []string) {
	p.mx.Lock()
	defer p.mx.Unlock()

	rotten := lo.Trim(p.previousObservatoryOutbounds, current)
	for _, outbound := range rotten {
		p.metrics.GetObservatory().
			DeleteLabelValues(metrics.LabelsObservatory{Outbound: outbound})
	}

	p.previousObservatoryOutbounds = current
}

func (p *Poller) gaugeObservatory(outbound string, input XrayObservatory) {
	labels := metrics.LabelsObservatory{Outbound: outbound}
	obs := p.metrics.GetObservatory()

	if !input.Alive && p.overrideOfflineDelay != 0 {
		obs.SetDelay(labels, p.overrideOfflineDelay.Milliseconds())
	} else {
		obs.SetAlive(labels, input.Alive)
	}

	obs.ObserveDelayHist(labels, input.Delay)
	obs.SetLastSeen(labels, input.LastSeenTime)
	obs.SetLastTry(labels, input.LastTryTime)
}
