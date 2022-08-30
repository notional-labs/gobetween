package cmd

/**
 * from-consul.go - pull config from consul and run
 *
 * @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
 */

import (
	"log"

	consul "github.com/hashicorp/consul/api"
	"github.com/notional-labs/gobetween/src/config"
	"github.com/notional-labs/gobetween/src/info"
	"github.com/notional-labs/gobetween/src/utils"
	"github.com/notional-labs/gobetween/src/utils/codec"
	"github.com/spf13/cobra"
)

/* Parsed options */
var consulKey string
var consulConfig consul.Config = consul.Config{}

/**
 * Add command
 */
func init() {
	FromConsulCmd.Flags().StringVarP(&consulKey, "key", "k", "gobetween", "Consul Key to pull config from")
	FromConsulCmd.Flags().StringVarP(&consulConfig.Scheme, "scheme", "s", "http", "http or https")

	RootCmd.AddCommand(FromConsulCmd)
}

/**
 * FromConsul command
 */
var FromConsulCmd = &cobra.Command{
	Use:   "from-consul <host:port>",
	Short: "Start using config from Consul",
	Long:  `Start using config from the Consul key-value storage`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}
			return
		}

		consulConfig.Address = args[0]
		client, err := consul.NewClient(&consulConfig)
		if err != nil {
			log.Fatal(err)
		}

		pair, _, err := client.KV().Get(consulKey, nil)
		if err != nil {
			log.Fatal(err)
		}

		if pair == nil {
			log.Fatal("Empty value for key " + consulKey)
		}

		datastr := string(pair.Value)
		if isConfigEnvVars {
			datastr = utils.SubstituteEnvVars(datastr)
		}

		var cfg config.Config
		if err := codec.Decode(datastr, &cfg, format); err != nil {
			log.Fatal(err)
		}

		info.Configuration = struct {
			Kind string `json:"kind"`
			Host string `json:"host"`
			Key  string `json:"key"`
		}{"consul", consulConfig.Address, consulKey}

		start(&cfg)
	},
}
