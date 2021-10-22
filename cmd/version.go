/*
Copyright © 2021 SUSE LLC
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/rancher-sandbox/luet-cosign/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Args:  cobra.ExactArgs(0),
		Short: "Print the version",
		RunE: func(cmd *cobra.Command, args []string) error {
			v := version.Get()
			if cmd.Flag("long").Changed {
				fmt.Printf("%#v", v)
			} else {
				fmt.Printf("%s+g%s", v.Version, v.GitCommit[:7])
			}

			return nil
		},
	}
	f := cmd.Flags()
	f.Bool("long", false, "Show long version info")
	return cmd
}