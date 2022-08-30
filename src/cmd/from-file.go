package cmd

/**
 * from-file.go - pull config from file and run
 *
 * @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
 */

import (
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/yyyar/gobetween/src/config"
	"github.com/yyyar/gobetween/src/info"
	"github.com/yyyar/gobetween/src/utils"
	"github.com/yyyar/gobetween/src/utils/codec"
)

/**
 * Add Root Command
 */
func init() {
	RootCmd.AddCommand(FromFileCmd)
}

/**
 * FromFile Command
 */
var FromFileCmd = &cobra.Command{
	Use:   "from-file <path>",
	Short: "Start using config from file",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		data, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatal(err)
		}

		var cfg config.Config

		datastr := string(data)
		if isConfigEnvVars {
			datastr = utils.SubstituteEnvVars(datastr)
		}

		if err = codec.Decode(datastr, &cfg, format); err != nil {
			log.Fatal(err)
		}

		info.Configuration = struct {
			Kind string `json:"kind"`
			Path string `json:"path"`
		}{"file", args[0]}

		start(&cfg)
	},
}
