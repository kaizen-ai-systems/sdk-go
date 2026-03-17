# sdk-go

Official Go SDK for [Kaizen AI Systems](https://www.kaizenaisystems.com).

## Installation

```bash
go get github.com/kaizen-ai-systems/sdk-go@latest
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    kaizen "github.com/kaizen-ai-systems/sdk-go/kaizen"
)

func main() {
    client := kaizen.NewClient(&kaizen.ClientConfig{APIKey: "your-api-key"})

    query, err := client.Akuma.Query(context.Background(), &kaizen.AkumaQueryRequest{
        Dialect: kaizen.DialectPostgres,
        Prompt:  "Top 10 customers by MRR last month",
        Mode:    kaizen.ModeSQLAndResults,
    })
    if err != nil {
        panic(err)
    }

    fmt.Println(query.SQL)
}
```

## Akuma (NL->SQL)

Persisted source APIs (`SetSchema`, `ListSources`, `CreateSource`, `SyncSource`, `DeleteSource`) require a dashboard-created DB-backed API key. Demo keys remain schema-less.

```go
resp, err := client.Akuma.Query(ctx, &kaizen.AkumaQueryRequest{
    Dialect: kaizen.DialectPostgres,
    Prompt:  "Show revenue by month for 2025",
    Mode:    kaizen.ModeSQLAndResults,
    Guardrails: &kaizen.Guardrails{
        ReadOnly:    true,
        AllowTables: []string{"orders", "customers"},
        MaxRows:     1000,
    },
})

_, err = client.Akuma.SetSchema(ctx, &kaizen.AkumaSchemaRequest{
    Name:    "Warehouse Manual Schema",
    Dialect: kaizen.DialectPostgres,
    Version: "2026-02-17",
    Tables: []kaizen.AkumaTable{
        {
            Name:        "orders",
            Description: "Customer orders",
            Columns: []kaizen.AkumaColumn{
                {Name: "id", Type: "uuid"},
                {Name: "customer_id", Type: "uuid"},
                {Name: "total_amount", Type: "numeric"},
            },
            PrimaryKey: []string{"id"},
        },
    },
})
if err != nil {
    panic(err)
}

_, err = client.Akuma.Query(ctx, &kaizen.AkumaQueryRequest{
    Dialect:  kaizen.DialectPostgres,
    Prompt:   "Show revenue by month for 2025",
    SourceID: "src_123",
})
if err != nil {
    panic(err)
}

_, err = client.Akuma.CreateSource(ctx, &kaizen.AkumaCreateSourceRequest{
    Name:             "Warehouse",
    Dialect:          kaizen.DialectPostgres,
    ConnectionString: "postgres://user:password@db.example.com:5432/app",
    TargetSchemas:    []string{"public"},
})
if err != nil {
    panic(err)
}

sources, err := client.Akuma.ListSources(ctx)
if err != nil {
    panic(err)
}
if len(sources.Sources) > 0 {
    _, _ = client.Akuma.SyncSource(ctx, sources.Sources[0].ID)
    _, _ = client.Akuma.DeleteSource(ctx, sources.Sources[0].ID)
}
```

## Enzan (GPU Cost)

```go
summary, err := client.Enzan.Summary(ctx, &kaizen.EnzanSummaryRequest{
    Window:  kaizen.Window24Hour,
    GroupBy: []kaizen.GroupByDimension{kaizen.GroupByProject, kaizen.GroupByModel},
})
if err != nil {
    panic(err)
}

fmt.Printf("Total: $%.2f\n", summary.Total.CostUSD)

byModel, err := client.Enzan.CostsByModel(ctx, &kaizen.EnzanModelCostRequest{
    Window: kaizen.Window30Day,
})
if err != nil {
    panic(err)
}
for _, row := range byModel.Rows {
    fmt.Printf("%s: $%.2f (%d queries)\n", row.Model, row.CostUSD, row.Queries)
}

pricing, err := client.Enzan.ListModelPricing(ctx)
if err != nil {
    panic(err)
}
for _, row := range pricing {
    fmt.Printf("%s %s: $%.5f / $%.5f per 1K tokens\n", row.Provider, row.Model, row.InputCostPer1KTokensUSD, row.OutputCostPer1KTokensUSD)
}
```

## Sozo (Synthetic Data)

```go
data, err := client.Sozo.Generate(ctx, &kaizen.SozoGenerateRequest{
    SchemaName: "saas_customers_v1",
    Records:    10000,
})
if err != nil {
    panic(err)
}

csvOutput, _ := data.ToCSV()
fmt.Println(csvOutput)
```

## Error Handling

```go
_, err := client.Akuma.Query(ctx, req)
if err != nil {
    switch e := err.(type) {
    case *kaizen.AuthError:
        fmt.Println("invalid API key", e.Status)
    case *kaizen.RateLimitError:
        fmt.Println("rate limited", e.RetryAfter)
    case *kaizen.KaizenError:
        fmt.Println("api error", e.Message)
    default:
        fmt.Println("unknown error", err)
    }
}
```
