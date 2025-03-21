// Copyright 2025 Chainguard, Inc.
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

package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"

	"chainguard.dev/melange/pkg/build"
	"chainguard.dev/melange/pkg/config"
	"github.com/chainguard-dev/clog"
	"github.com/chainguard-dev/clog/slag"
)

// Simple wrapper for executing shell commands.
func runCommand(cmd []string) {
	proc := exec.Command(cmd[0], cmd[1:]...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	if err := proc.Run(); err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <path-to-melange-yaml> <destination-path>", os.Args[0])
	}

	filePath := os.Args[1]
	destDir := os.Args[2]

	// Now, create a new slag.Level object and set the level to 'error'.
	// This enables display of melange-level logs. Change the level to
	// info or debug if you want to see more logs.
	var level slag.Level
	level.Set("error")
	slog.SetDefault(slog.New(clog.NewHandler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: &level}))))
	log := clog.New(slog.Default().Handler())
	ctx := clog.WithLogger(context.Background(), log)

	// Parsing the configuration file. It's still not ready for 'consumption'
	// though!
	cfg, err := config.ParseConfiguration(ctx, filePath)
	if err != nil {
		log.Fatalf("Failed to parse melange config: %v", err)
	}

	// Prepare the substitution map and compile the pipelines, making sure that
	// the resulting pipeline run statements are all substituted with the
	// correct values and ready for execution.
	c := &build.Compiled{
		PipelineDirs: []string{},
	}
	sm, err := build.NewSubstitutionMap(cfg, "amd64", "gnu", nil)
	if err != nil {
		log.Fatalf("Failed to create substitution map: %v", err)
	}
	err = c.CompilePipelines(ctx, sm, cfg.Pipeline)
	if err != nil {
		log.Fatalf("Failed to compile pipelines: %v", err)
	}

	// Iterate over the pipeline steps and look for any source fetching steps.
	for _, step := range cfg.Pipeline {
		if step.Uses != "git-checkout" && step.Uses != "fetch" {
			continue
		}

		// TODO: Make this smarter. For now we assume the first source fetch
		//  operation is 'the one' and we just finish. Let's improve this
		//  in the near future.

		fmt.Printf("Found source fetching step: %s.\nFetching source to %s\n", step.Uses, destDir)
		os.MkdirAll(destDir, 0755)
		os.Chdir(destDir)
		cmd := []string{"/bin/sh", "-c", step.Pipeline[0].Runs}
		runCommand(cmd)

		break
	}
}
