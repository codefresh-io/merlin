# Merlin

```
A command line application for a Codefresh developer

Usage:
  merlin [command]

Available Commands:
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

* Create merlin config by running `merlin init`
 ```
  merlin init --help
 ```

 ## Example

Show all operators that an environment exposes
```
merlin describe environment
```

Show all params that a operator requires
```
merlin describe operator [NAME]
```

 # Merlin definitions
 * `config.yaml` - configuration file created by `merlin init` command - hold all the environments available
 * `environment.js` - JS file , defines the environment and all corresponding operators and components
 * operator - A unit (functin) that describes how to interact with an Environment or a Component. Each operator has a name
 * component - Logical part of the environment


 ## Multiple operator execution flows
Running a `merlin run [NAME]` command will search the all operators named: NAME on the environment level and executed all of them sequentially (`merlin list` can help you to understand which operators have multiple executions).
A few notes to have in mind:
* If a flag `--component COMPONENT` set, operators from environment level and component level will be executed, where the environment level is in priority.
* Next operator executed only when the previous one finished successfully
* Operators executions flow are sharing environment variables exist on the host process (that runs the `merlin run` command)
