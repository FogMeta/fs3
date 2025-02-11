// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"github.com/joho/godotenv"
	"github.com/minio/cli"
	sysconfig "github.com/minio/minio/config"
	"github.com/minio/minio/internal/config"
	"github.com/minio/minio/logs"
	"github.com/minio/minio/scheduler"
	"github.com/minio/pkg/console"
	"github.com/minio/pkg/trie"
	"github.com/minio/pkg/words"
	"os"
	"path/filepath"
	"sort"
)

// GlobalFlags - global flags for minio.
var GlobalFlags = []cli.Flag{
	// Deprecated flag, so its hidden now - existing deployments will keep working.
	cli.StringFlag{
		Name:   "config-dir, C",
		Value:  defaultConfigDir.Get(),
		Usage:  "[DEPRECATED] path to legacy configuration directory",
		Hidden: true,
	},
	cli.StringFlag{
		Name:  "certs-dir, S",
		Value: defaultCertsDir.Get(),
		Usage: "path to certs directory",
	},
	cli.BoolFlag{
		Name:  "quiet",
		Usage: "disable startup information",
	},
	cli.BoolFlag{
		Name:  "anonymous",
		Usage: "hide sensitive information from logging",
	},
	cli.BoolFlag{
		Name:  "json",
		Usage: "output server logs and startup information in json format",
	},
	// Deprecated flag, so its hidden now, existing deployments will keep working.
	cli.BoolFlag{
		Name:   "compat",
		Usage:  "enable strict S3 compatibility by turning off certain performance optimizations",
		Hidden: true,
	},
	// This flag is hidden and to be used only during certain performance testing.
	cli.BoolFlag{
		Name:   "no-compat",
		Usage:  "disable strict S3 compatibility by turning on certain performance optimizations",
		Hidden: true,
	},
}

// Help template for minio.
var minioHelpTemplate = `NAME:
  {{.Name}} - {{.Usage}}

DESCRIPTION:
  {{.Description}}

USAGE:
  {{.HelpName}} {{if .VisibleFlags}}[FLAGS] {{end}}COMMAND{{if .VisibleFlags}}{{end}} [ARGS...]

COMMANDS:
  {{range .VisibleCommands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
  {{end}}{{if .VisibleFlags}}
FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}
VERSION:
  {{.Version}}
`

func newApp(name string) *cli.App {
	// Collection of minio commands currently supported are.
	commands := []cli.Command{}

	// Collection of minio commands currently supported in a trie tree.
	commandsTree := trie.NewTrie()

	// registerCommand registers a cli command.
	registerCommand := func(command cli.Command) {
		commands = append(commands, command)
		commandsTree.Insert(command.Name)
	}

	findClosestCommands := func(command string) []string {
		var closestCommands []string
		closestCommands = append(closestCommands, commandsTree.PrefixMatch(command)...)

		sort.Strings(closestCommands)
		// Suggest other close commands - allow missed, wrongly added and
		// even transposed characters
		for _, value := range commandsTree.Walk(commandsTree.Root()) {
			if sort.SearchStrings(closestCommands, value) < len(closestCommands) {
				continue
			}
			// 2 is arbitrary and represents the max
			// allowed number of typed errors
			if words.DamerauLevenshteinDistance(command, value) < 2 {
				closestCommands = append(closestCommands, value)
			}
		}

		return closestCommands
	}

	// Register all commands.
	registerCommand(serverCmd)
	registerCommand(gatewayCmd)

	// Set up app.
	cli.HelpFlag = cli.BoolFlag{
		Name:  "help, h",
		Usage: "show help",
	}

	app := cli.NewApp()
	app.Name = name
	app.Author = "MinIO, Inc."
	app.Version = ReleaseTag
	app.Usage = "High Performance Object Storage"
	app.Description = `Build high performance data infrastructure for machine learning, analytics and application data workloads with MinIO`
	app.Flags = GlobalFlags
	app.HideHelpCommand = true // Hide `help, h` command, we already have `minio --help`.
	app.Commands = commands
	app.CustomAppHelpTemplate = minioHelpTemplate
	app.CommandNotFound = func(ctx *cli.Context, command string) {
		console.Printf("‘%s’ is not a minio sub-command. See ‘minio --help’.\n", command)
		closestCommands := findClosestCommands(command)
		if len(closestCommands) > 0 {
			console.Println()
			console.Println("Did you mean one of these?")
			for _, cmd := range closestCommands {
				console.Printf("\t‘%s’\n", cmd)
			}
		}

		os.Exit(1)
	}

	return app
}

// Main main for minio server.
func Main(args []string) {
	// Set the minio app name.
	appName := filepath.Base(args[0])

	initConfigAndLog()
	initUserConfig(sysconfig.GetSysConfig().StandAlone)
	scheduler.SendDealScheduler()
	scheduler.BackupScheduler()
	scheduler.RebuildScheduler()

	logs.GetLogger().Info("Your FS3 Server is running successfully. Please copy and paste the url below to open in a browser")

	// Run the app - exit on error.
	if err := newApp(appName).Run(args); err != nil {
		os.Exit(1)
	}
}

func initUserConfig(standAlone bool) {
	if standAlone {
		LoadEnv()
	}
	swanAddress := os.Getenv("SWAN_ADDRESS")
	fs3VolumeAddress := os.Getenv("FS3_VOLUME_ADDRESS")
	fs3WalletAddress := os.Getenv("FS3_WALLET_ADDRESS")
	carFileSize := os.Getenv("CAR_FILE_SIZE")
	ipfsApiAddress := os.Getenv("IPFS_API_ADDRESS")
	ipfsGateway := os.Getenv("IPFS_GATEWAY")
	swanToken := os.Getenv("SWAN_TOKEN")
	lotusClientApiUrl := os.Getenv("LOTUS_CLIENT_API_URL")
	lotusClientAccessToken := os.Getenv("LOTUS_CLIENT_ACCESS_TOKEN")
	volumeBackupAddress := os.Getenv("VOLUME_BACKUP_ADDRESS")
	psqlHost := os.Getenv("PSQL_HOST")
	psqlUser := os.Getenv("PSQL_USER")
	psqlPassword := os.Getenv("PSQL_PASSWORD")
	psqlDbname := os.Getenv("PSQL_DBNAME")
	psqlPort := os.Getenv("PSQL_PORT")
	//logs.GetLogger().Println(swanAddress, fs3VolumeAddress, fs3WalletAddress, carFileSize, ipfsApiAddress, ipfsGateway, swanToken, lotusClientApiUrl, lotusClientAccessToken, volumeBackupAddress)
	config.InitUserConfig(swanAddress, fs3VolumeAddress, fs3WalletAddress, carFileSize, ipfsApiAddress, ipfsGateway, swanToken, lotusClientApiUrl, lotusClientAccessToken, volumeBackupAddress, psqlHost, psqlUser, psqlPassword, psqlDbname, psqlPort)

}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		logs.GetLogger().Error(err)
	}
}

func initConfigAndLog() {
	logs.InitLogger()
	sysconfig.InitSysConfig("")
}
