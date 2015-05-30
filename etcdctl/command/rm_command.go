// Copyright 2015 CoreOS, Inc.
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

package command

import (
	"errors"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/spf13/cobra"
)

// NewRemoveCommand returns the CLI command for "rm".
func NewRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm",
		Short: "remove a key",
		Run: func(cmd *cobra.Command, args []string) {
			handleAll(cmd, args, removeCommandFunc)
		},
	}
	cmd.Flags().Bool("dir", false, "removes the key if it is an empty directory or a key-value pair")
	cmd.Flags().Bool("recursive", false, "removes the key and all child keys(if it is a directory)")
	cmd.Flags().String("with-value", "", "previous value")
	cmd.Flags().Uint64("with-index", 0, "previous index")
	return cmd
}

// removeCommandFunc executes the "rm" command.
func removeCommandFunc(cmd *cobra.Command, args []string, client *etcd.Client) (*etcd.Response, error) {
	if len(args) == 0 {
		return nil, errors.New("Key required")
	}
	key := args[0]
	recursive, _ := cmd.Flags().GetBool("recursive")
	dir, _ := cmd.Flags().GetBool("dir")

	// TODO: distinguish with flag is not set and empty flag
	// the cli pkg need to provide this feature
	prevValue, _ := cmd.Flags().GetString("with-value")
	prevIndex, _ := cmd.Flags().GetUint64("with-index")

	if prevValue != "" || prevIndex != 0 {
		return client.CompareAndDelete(key, prevValue, prevIndex)
	}

	if recursive || !dir {
		return client.Delete(key, recursive)
	}

	return client.DeleteDir(key)
}
