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

package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mudler/go-pluggable"
	"github.com/mudler/luet/pkg/bus"
	"github.com/rancher-sandbox/luet-cosign/pkg/helpers"
	"github.com/rancher-sandbox/luet-cosign/pkg/log"
	"os"
	"os/exec"
)

type UnpackEvent struct {
	Name string `json:"name"`
	Data string `json:"data"`
	File string `json:"file"`
}

type ImagePostPushData struct {
	ImageName string `json:"ImageName"`
}

type LuetEvent struct {
	event   pluggable.EventType
	payload string
}


func NewEventDispatcherAction(event string, payload string) *LuetEvent {
	return &LuetEvent{event: pluggable.EventType(event), payload: payload}
}


func (event LuetEvent) Run() (map[string]string, error) {
	log.Log("Got event: %s\n", event.event)
	switch event.event {
	case bus.EventImagePrePush:
		// We need to fail on the pre-push, otherwise the image will be pushed but the signature not, so we better make as many checks here as we can
		pass := os.Getenv("COSIGN_PASSWORD")
		keyLocation := os.Getenv("COSIGN_KEY_LOCATION")
		if pass == "" || keyLocation == "" {
			return helpers.WrapErrorMap(errors.New("missing cosign env vars COSIGN_PASSWORD or COSIGN_KEY_LOCATION"))
		}
		_, err := unPackImagePostPushDataPayload(event.payload)
		if err != nil {
			return helpers.WrapErrorMap(err)
		}
		return helpers.WrapErrorMap(nil)
	case bus.EventImagePostPush:
		// Do the checks again in case something changed between the 2 events
		pass := os.Getenv("COSIGN_PASSWORD")
		keyLocation := os.Getenv("COSIGN_KEY_LOCATION")
		if pass == "" || keyLocation == "" {
			return helpers.WrapErrorMap(errors.New("missing cosign env vars COSIGN_PASSWORD or COSIGN_KEY_LOCATION"))
		}
		data, err := unPackImagePostPushDataPayload(event.payload)
		if err != nil {
			return helpers.WrapErrorMap(err)
		}
		log.Log("Signing image: %s", data.ImageName)

		args := fmt.Sprintf("echo -n '%s' | cosign sign -key %s %s", pass, keyLocation, data.ImageName)
		out, err := exec.Command("bash", "-c", args).CombinedOutput()

		if err != nil {
			return helpers.WrapErrorMap(err)
		} else {
			// enhance return values with the command output
			ret, _ := helpers.WrapErrorMap(err)
			ret["state"] = string(out)
			log.Log("Finished signing and pushing %s", data.ImageName)
			return ret, err
		}
	default:
		log.Log("No event that I can recognize")
		return helpers.WrapErrorMap(nil)
	}
}


func unPackImagePostPushDataPayload(payload string) (ImagePostPushData, error) {
	payloadTmp := pluggable.Event{}
	dataTmp := ImagePostPushData{}
	// unpack payload
	err := json.Unmarshal([]byte(payload), &payloadTmp)
	if err != nil {
		log.Log("Error while unmarshalling payload")
		log.Log("Payload: %s", payload)
		return dataTmp, err
	}
	// unpack data inside payload
	err = json.Unmarshal([]byte(payloadTmp.Data), &dataTmp)
	if err != nil {
		log.Log("Error while unmarshalling data from the payload")
		log.Log("Payload: %s", payloadTmp.Data)
		return dataTmp, err
	}

	if dataTmp.ImageName == "" {
		log.Log("Some fields are missing from the event, cannot continue")
		return dataTmp, errors.New("field ImageName missing from payload")
	}

	return dataTmp, nil
}