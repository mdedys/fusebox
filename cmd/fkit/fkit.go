package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github/mdedys/fusebox/playground"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var (
	noop = func(context.Context, []string) error { return flag.ErrHelp }
)

func playgroundCmd() *ffcli.Command {

	generate := &ffcli.Command{
		Name:       "generate",
		ShortUsage: "fkit playground generate <num_nodes>",
		ShortHelp:  "generates docker-compose file based on described environment",
		LongHelp:   "Will generate a docker-compose environment based on the settings provided by the user",
		Exec: func(ctx context.Context, args []string) error {
			nodes, err := strconv.Atoi(args[0])
			if err != nil {
				nodes = 5
			}

			err = playground.Generate(nodes)
			if err != nil {
				return err
			}

			return nil
		},
	}

	up := &ffcli.Command{
		Name:       "up",
		ShortUsage: "fkit playground up",
		ShortHelp:  "starts playground",
		LongHelp:   "Starts playground from generate docker-compose",
		Exec: func(ctx context.Context, args []string) error {
			composeFile, err := playground.GetComposeFilepath()
			if err != nil {
				return err
			}
			cmd := "docker"
			cmdargs := []string{"compose", "-f", composeFile, "up", "-d"}
			_, err = exec.Command(cmd, cmdargs...).Output()
			if err != nil {
				return err
			}
			return nil
		},
	}

	down := &ffcli.Command{
		Name:       "down",
		ShortUsage: "fkit playground down",
		ShortHelp:  "stops playground",
		LongHelp:   "Stops playground from generate docker-compose",
		Exec: func(ctx context.Context, args []string) error {
			composeFile, err := playground.GetComposeFilepath()
			if err != nil {
				return err
			}
			cmd := "docker"
			cmdargs := []string{"compose", "-f", composeFile, "down"}
			_, err = exec.Command(cmd, cmdargs...).Output()
			if err != nil {
				return err
			}
			return nil
		},
	}

	return &ffcli.Command{
		Name:        "playground",
		ShortUsage:  "fkit playground <subcommand>",
		ShortHelp:   "commands for interacting with playground",
		LongHelp:    "Commands for interacting with playground",
		Subcommands: []*ffcli.Command{generate, up, down},
		Exec:        noop,
	}
}

func main() {

	fkit := &ffcli.Command{
		ShortUsage: "fkit [global flags] <subcommand> [subcommand flags] [subcommand args]",
		ShortHelp:  "cli for interacting with nostr",
		LongHelp: `fkit is a CLI client for interacting with nostr and the fkit playground environment.
Each subcommand (e.g. playground) contains a set of commands to interact with that resources.

Global flags can be set with shell environment variables for ease of use.
		`,
		Subcommands: []*ffcli.Command{playgroundCmd()},
		Options:     []ff.Option{ff.WithEnvVarPrefix("FUSEKIT")},
		Exec:        noop,
	}

	if err := fkit.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "fkit: %v\n", err)
		}
		os.Exit(1)
	}
}
