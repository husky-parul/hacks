// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Example using the OTLP exporter + collector + third-party backends. For
// information about using the exporter, see:
// https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc#example-package-Insecure
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func setupTracing(ctx context.Context, serviceName, collectorAddr string) (tracesdk.SpanExporter, error) {

	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint("localhost:4317"),
		otlpgrpc.WithDialOption(grpc.WithBlock()), // useful for testing
	)

	exporter, err := otlp.NewExporter(ctx,
		driver)
	if err != nil {
		return nil, err
	}
	res := resource.NewWithAttributes(
		semconv.ServiceNameKey.String(serviceName),
	)
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSyncer(exporter),
		tracesdk.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	log.Println("exporter: ", exporter)

	return exporter, err
}

func main() {
	ctx := context.Background()
	exporter, err := setupTracing(ctx, "my-service", "localhost:4317")
	if err != nil {
		handleErr(err, err.Error())
	}
	defer exporter.Shutdown(ctx)

	tracer := otel.Tracer("component_name")
	ctx, span := tracer.Start(ctx, "say-hello", trace.WithNewRoot())
	println("helloStr")
	span.End()

	for i := 0; i < 20; i++ {
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
		log.Printf("Doing really hard work (%d / 10)\n", i+1)

		<-time.After(time.Second)
		iSpan.End()
	}

	log.Printf("Done!")

}

func handleErr(err error, message string) {
	if err != nil {
		log.Printf("************ Error ****************8")
		log.Fatalf("%s: %v", message, err)
	}
}
