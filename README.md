# luet-cosign

luet-cosign is a plugin for [luet](https://luet-lab.github.io/docs/) to create and push signatures for containers using [cosign](https://github.com/sigstore/cosign).

## Event parsing

When used as a plugin to luet (by calling luet with `--plugin luet-cosign`, see [plugin docs](https://luet-lab.github.io/docs/docs/concepts/plugins-and-extensions/)) luet emits events based on the actions being performed.
We take those events and payloads and execute an action depending on their contents.


### Signing

This action uses the event `EventImagePostPush`

Using it as a plugin for luet requires 2 environment variables in order to make cosign work properly:

 - COSIGN_PASSWORD: The password for the private key file
 - COSIGN_KEY_LOCATION: The location of the private key file

luet-cosign will use both those values to call cosign on the pushed container while publishing a repo and will create and push the signatures along the containers.

There is an optional parameter `COSIGN_FULCIO_URL` to override the fulcio url passed to cosign. Defaults to `https://fulcio.sigstore.dev`


### Verifying

This action uses the event `EventImagePreUnPack`

Using it as a plugin for luet requires 2 environment variables in order to make cosign work properly:

- COSIGN_PUBLIC_KEY_LOCATION: The location of the public key file (can be a file, url, KMS URI or Kubernetes Secret)

luet-cosign will use that key to call cosign on the pulled artifact and verify the signature.


### Skip verifying some artifacts

While verifying you can set the `COSIGN_SKIP` env var to a list (space separated) of regex values used to skip verification of some artifacts, or even the full repo.

For example, we want to skip verifying on luet-cosign packages older than `0.0.6` (included) because we didn't start signing those packages until `0.0.7`:

```bash
$ export COSIGN_SKIP=".*luet-cosign.*toolchain.*0.0.[0-6].*"
$ luet --plugin luet-cosign install toolchain/luet-cosign 
 INFO   🍭 Enabled plugins:
 INFO   	➡  luet-cosign (at /usr/local/bin/luet-cosign)
           
  Install  
           
 INFO   Downloading quay.io/costoolkit/releases-green:repository.yaml                                     
 SUCCESS   🍭  Plugin luet-cosign at /usr/local/bin/luet-cosign succeded, state reported:                 
 INFO   quay.io/costoolkit/releases-green:repository.yaml verified. See luet-cosign logs for full info.   
 INFO   Pulled: sha256:de493f32c460cfbdd97258ecf1d7c5ccd151feb6b8fc4e44b2c098f222528fa5                   
 INFO   Size: 628B                                                                                        
 INFO   Downloading quay.io/costoolkit/releases-green:tree.tar.zst                                        
 SUCCESS   🍭  Plugin luet-cosign at /usr/local/bin/luet-cosign succeded, state reported:                 
 INFO   quay.io/costoolkit/releases-green:tree.tar.zst verified. See luet-cosign logs for full info.      
 INFO   Pulled: sha256:f5561ccf05338a574e36273ed0e37f07109fcd25d9527747ebcae3d54ed6c155                   
 INFO   Size: 5.205KiB                                                                                    
 INFO   Downloading quay.io/costoolkit/releases-green:repository.meta.yaml.tar.zst                        
 SUCCESS   🍭  Plugin luet-cosign at /usr/local/bin/luet-cosign succeded, state reported:                 
 INFO   quay.io/costoolkit/releases-green:repository.meta.yaml.tar.zst verified. See luet-cosign logs for full info.
 INFO   Pulled: sha256:11684b8d37dd5d53a6806e5390bbe1bb70c1f27a6643ce54707ab1e81344b567                   
 INFO   Size: 193.3KiB                                                                                    
 INFO   🏠  Repository cOS revision: 137 (2021-11-02 04:53:22 +0100 CET)                                  
 INFO   ℹ  Repository: cos Priority: 1 Type: docker                                                       
 INFO   Packages that are going to be installed in the system:                                            

Program Name          | Version | License | Repository
toolchain/luet-cosign | 0.0.6-1 |         | cos       

 INFO   By going forward, you are also accepting the licenses of the packages that you are going to install in your system.
 INFO   Do you want to continue with this operation? [y/N]: 
y
Downloading packages [0/1] █                                                                       0% | 0s
 INFO   Downloading image quay.io/costoolkit/releases-green:luet-cosign-toolchain-0.0.6-1                 
 SUCCESS   🍭  Plugin luet-cosign at /usr/local/bin/luet-cosign succeded, state reported:                 
 INFO   Image quay.io/costoolkit/releases-green:luet-cosign-toolchain-0.0.6-1 found in skip list (.*luet-cosign.*toolchain.*0.0.[0-6].*)
 INFO   Pulled: sha256:9377596ed7f98c084a86a96f0536ddbfacf7a3e90567e8b577745b602d4d6d3d                   
 INFO   Size: 3.09MiB                                                                                     

Downloading packages [1/1] █████████████████████████████████████████████████████████████████████ 100% | 1s
 INFO   Checking for file conflicts..
 INFO   📦  Package  toolchain/luet-cosign-0.0.6-1 installed 

```

As you can see above, the `luet-cosign-toolchain-0.0.6-1` artifact was matched and skipped, while the rest of the artifacts were verified (tree, meta, repository.yaml)

As soon as version `0.0.7` is published, the same skip list won't match and the package will be verified.

And if this package pulled some dependencies, those would also be verified properly, so `COSIGN_SKIP` is meant to be used only on special occasions when there is no possibility of verifying the artifact and only to skip specific packages so the rest of the chain is verified fully.



### Keyless signing/verify

You can use the experimental keyless verify/sign process by setting the `COSIGN_EXPERIMENTAL=1` env var.
Please see the [upstream docs for cosign](https://github.com/sigstore/cosign/blob/main/KEYLESS.md)
This is only possible to do in an CI environment (github) if using luet-cosign as a plugin as it requires an OIDC token to be available.

### Manual use of luet-cosign

You can manually test those events by calling luet-cosign with no subcommands and 2 params. The first being the event emitted (see events emitted by luet [here](https://github.com/mudler/luet/blob/master/pkg/api/core/bus/events.go)) and the second a json payload, the contents depend on the type of event.

For example, for signing:
```bash
export COSIGN_PASSWORD=whatever
export COSIGN_KEY_LOCATION=/tmp/mykey.key
luet-cosign 'image.post.push' '{"data": {"ImageName": "quay.io/costoolkit/releases-opensuse:systemd-boot-live-26"}}'
```

Or verifying:
```bash
export COSIGN_PUBLIC_KEY_LOCATION=http://some.host/public.key.pub
luet-cosign 'image.pre.unpack' '{"data": {"Image": "quay.io/costoolkit/releases-opensuse:systemd-boot-live-26"}}'
```

### Debug

You can enable debug of cosign by setting `COSIGN_DEBUG=1`


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