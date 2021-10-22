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
)

// Basic stub to call the action in the package, does nothing really
func newEventCmd(args []string) error {
	event := args[0]
	payload := args[1]
	fmt.Printf("event: %s\npayload%s\n", event, payload)


	// As this is part of being a luet plugin we need to output to console ONLY the results in json formatting so luet
	// can parse it.


	// Let the root cmd be the one to set the exit status as success/failure
	return nil
}