package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kcarretto/paragon/graphql"
	"github.com/kcarretto/paragon/graphql/models"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

// Schema defines the structure of JSON data used to initialize the Teamserver.
//
// Example:
//
//    {
//     "targets": [
//         {
//             "name": "Target-10",
//             "ip": "10.1.0.10",
//             "tags": ["A", "B", "C"]
//         },
//         {
//             "name": "Target-11",
//             "ip": "10.1.0.11",
//             "tags": ["A", "B"]
//         }
//         {
//             "name": "Target-12",
//             "ip": "10.1.0.12"
//         }
//     ]
//    }
type Schema struct {
	Targets []struct {
		Name      string   `json:"name"`
		PrimaryIP string   `json:"ip"`
		OS        string   `json: "os"`
		Tags      []string `json:"tags"`
	} `json:"targets"`
}

func main() {
	ctx := context.Background()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = logger.Named("teamserver.init")

	teamserverURL := "http://127.0.0.1:80"
	if url := os.Getenv("TEAMSERVER_URL"); url != "" {
		teamserverURL = url
	}
	graph := &graphql.Client{
		URL:     fmt.Sprintf("%s/%s", teamserverURL, "graphql"),
		Service: "pg-teamserver-init",
	}

	var schemaPath string

	app := &cli.App{
		Name:  "init",
		Usage: "Initialize the Teamserver with data generated from a schema.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "schema",
				Usage:       "Path to the JSON schema-file to use for initialization.",
				Destination: &schemaPath,
			},
		},
		Action: func(c *cli.Context) error {
			logger = logger.With(zap.String("schema", schemaPath))
			logger.Debug("Populating Teamserver")

			// Get schema file
			schemaData, err := ioutil.ReadFile(schemaPath)
			if err != nil {
				return fmt.Errorf("failed to read schema file: %w", err)
			}

			// Parse schema JSON
			var schema Schema
			if err := json.Unmarshal(schemaData, &schema); err != nil {
				return fmt.Errorf("failed to parse schema file: %w", err)
			}

			// Assert at least one target provided
			if len(schema.Targets) < 1 {
				return fmt.Errorf("Must specify at least one target in the schema")
			}

			// Create all tags defined in the schema
			tagNames := make([]string, 0, len(schema.Targets))
			for _, target := range schema.Targets {
				if target.Tags == nil {
					continue
				}

				tagNames = append(tagNames, target.Tags...)
			}
			tagMap, err := graph.CreateTags(ctx, tagNames...)
			if err != nil {
				return err
			}
			logger.Info("Successfully created all tags", zap.Strings("tags", tagNames))

			// Create all targets defined in the schema
			requests := make([]models.CreateTargetRequest, 0, len(schema.Targets))
			for _, target := range schema.Targets {
				// Resolve all tag ids for the target
				tagIDs := make([]int, 0, len(target.Tags))
				for _, name := range target.Tags {
					tag, ok := tagMap[name]
					if !ok {
						return fmt.Errorf("failed to retrieve tag: %q", name)
					}
					tagIDs = append(tagIDs, tag.ID)
				}

				// Append request to the array
				requests = append(requests, models.CreateTargetRequest{
					Name:      target.Name,
					PrimaryIP: target.PrimaryIP,
					Os:        target.OS,
					Tags:      tagIDs,
				})
			}

			logger.Debug("Populating Teamserver Targets", zap.Reflect("requests", requests))

			for _, req := range requests {
				if _, err := graph.CreateTarget(ctx, req); err != nil {
					logger.Error("Failed to create target",
						zap.Error(err),
						zap.Reflect("target", req),
					)
				}
			}

			logger.Info("Finished Populating Teamserver")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal("failed to populate teamserver", zap.Error(err))
	}
}
