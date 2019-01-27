# Merlin

```
A command line application for a Codefresh developer

Usage:
  merlin [command]

Available Commands:
  create      A command line application for a Codefresh developer
  help        Help about any command
  init        Create config file
  run         Run command
  version     Print merlin version

Flags:
  -h, --help                  help for merlin
      --merlinconfig string   overwrite merlin default config path
      --verbose               get extra logs

Use "merlin [command] --help" for more information about a command.
```

## Installation
* Prerequisite:
    * [codefresh cli](http://cli.codefresh.io)
    * [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl)
    * [telepresene](https://github.com/telepresenceio/telepresence)
* Install latest release from [here](https://github.com/codefresh-io/merlin/releases)

* Create merlin config by running `merlin init [name]`, provide all required flags
 ```
  merlin init --help
 ```

 * Create your environment
 ```
  merlin create [NAME]
 ```

 ## Example
 Debug cfapi
 * Connect to cfapi
 ```
  merlin run connect --component cfapi
 ```

 * Start cfapi service locally
 ```
  merlin run start --component cfapi
 ```