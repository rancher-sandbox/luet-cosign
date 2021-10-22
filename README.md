# luet-cosign

luet-cosign is a plugin for [luet](https://luet-lab.github.io/docs/) to create and push signatures for containers using [cosign](https://github.com/sigstore/cosign).

### event parsing

When used as a plugin to luet (by calling luet with `--plugin luet-cosign`, see [plugin docs](https://luet-lab.github.io/docs/docs/concepts/plugins-and-extensions/)) luet emits events based on the actions being performed.
We take those events and payloads and execute an action depending on their contents.

Using it as a plugin for luet requires 2 environment variables in order to make cosign work properly:

 - COSIGN_PASSWORD: The password for the private key file
 - COSIGN_KEY_LOCATION: The location of the private key file

luet-cosing will use both those values to call cosign on the pushed container while publishing a repo and will create and push the signatures along the containers.

You can manually test those events by calling luet-cosign with no subcommands and 2 params. The first being the event emitted (see events emitted by luet [here](https://github.com/mudler/luet/blob/master/pkg/bus/events.go)) and the second a json payload, the contents depend on the type of event.

```bash
luet-cosign 'image.post.push' '{"data": {"ImageName": "quay.io/costoolkit/releases-opensuse:systemd-boot-live-26"}}'
```



## License

Copyright (c) 2021 [SUSE, LLC](http://suse.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.