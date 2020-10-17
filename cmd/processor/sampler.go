package main

import (
	"gitlab.com/Cacophony/go-kit/events"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/trace"
)

type customSampler struct {
	base trace.Sampler
}

func (s *customSampler) ShouldSample(parameters trace.SamplingParameters) trace.SamplingResult {
	if parameters.HasRemoteParent && parameters.ParentContext.IsSampled() {
		return trace.SamplingResult{
			Decision: trace.RecordAndSample,
		}
	}

	for _, attr := range parameters.Attributes {
		if attr.Key != events.SpanLabelEventingIsCommand {
			continue
		}

		if attr.Value.Type() == label.BOOL && attr.Value.AsBool() {
			return trace.SamplingResult{
				Decision: trace.RecordAndSample,
			}
		}
		break
	}

	return s.base.ShouldSample(parameters)
}

func (s *customSampler) Description() string {
	return "custom sampler"
}
