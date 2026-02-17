package main

import (
	"context"
	"fmt"
	"os"

	kaizen "github.com/kaizen-ai-systems/sdk-go/kaizen"
)

func main() {
	client := kaizen.NewClient(&kaizen.ClientConfig{APIKey: os.Getenv("KAIZEN_API_KEY")})

	result, err := client.Akuma.Query(context.Background(), &kaizen.AkumaQueryRequest{
		Dialect: kaizen.DialectPostgres,
		Prompt:  "Top 10 customers by MRR last month",
		Mode:    kaizen.ModeSQLOnly,
	})
	if err != nil {
		panic(err)
	}
	if result.Error != "" {
		panic(result.Error)
	}

	fmt.Println(result.SQL)
}
