/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"path"
	"strings"

	"github.com/labring/sealvm/pkg/apply"
	fileutil "github.com/labring/sealvm/pkg/utils/file"
	"github.com/labring/sealvm/pkg/utils/maps"
	v1 "github.com/labring/sealvm/types/api/v1"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
func newRunCmd() *cobra.Command {
	vm := v1.VirtualMachine{}
	var nodes int
	var dev bool
	var src string
	var defaultMount = fmt.Sprintf("%s:%s", path.Join(os.Getenv("GOPATH"), "src"), "/root/go/src")
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			applier, err := apply.NewApplierFromArgs(&vm)
			if err != nil {
				return err
			}
			return applier.Apply()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if strings.Contains(vm.Name, "-") {
				return fmt.Errorf("your cluster name contains chart '-' ")
			}

			if err := checkInstall(vm.Spec.Type); err != nil {
				return err
			}
			if dev {
				if src == "" {
					return fmt.Errorf("src must be set")
				}
				mounts := maps.StringToMap(src, ",")

				vm.Spec.Hosts = append(vm.Spec.Hosts, v1.Host{
					Role:   v1.GOLANG,
					Count:  1,
					Mounts: mounts,
					Resources: map[string]int{
						v1.CPUKey:  2,
						v1.MEMKey:  4,
						v1.DISKKey: 50,
					},
				})
			}
			if nodes != 0 {
				vm.Spec.Hosts = append(vm.Spec.Hosts, v1.Host{
					Role:   v1.NODE,
					Count:  nodes,
					Mounts: map[string]string{},
					Resources: map[string]int{
						v1.CPUKey:  2,
						v1.MEMKey:  4,
						v1.DISKKey: 50,
					},
				})
			}
			return nil
		},
	}
	runCmd.Flags().StringVarP(&vm.Spec.SSH.PkFile, "pk", "i", path.Join(fileutil.GetHomeDir(), ".ssh", "id_rsa"), "selects a file from which the identity (private key) for public key authentication is read")
	runCmd.Flags().StringVar(&vm.Spec.SSH.PkPasswd, "pk-passwd", "", "passphrase for decrypting a PEM encoded private key")
	runCmd.Flags().StringVarP(&vm.Spec.SSH.PublicFile, "pub", "p", path.Join(fileutil.GetHomeDir(), ".ssh", "id_rsa.pub"), "selects a file from which the identity (public key) for public key authentication is read")
	runCmd.Flags().StringVarP(&vm.Spec.Type, "type", "t", v1.MultipassType, "choose a type of infra, multipass")
	runCmd.Flags().StringVarP(&vm.Name, "name", "n", "default", "name of cluster to applied init action")
	runCmd.Flags().IntVarP(&nodes, "nodes", "w", 0, "number of nodes")
	runCmd.Flags().BoolVarP(&dev, "dev", "d", false, "number of dev")
	runCmd.Flags().StringVarP(&src, "dev-mounts", "s", defaultMount, "gopath src dir")
	return runCmd
}

func init() {
	rootCmd.AddCommand(newRunCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
