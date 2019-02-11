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

* Create merlin config by running `merlin init [name]`, provide all required flags
 ```
  merlin init --help
 ```

 ## Example

Show all commands that an environment provides
```
merlin list
```

 # Merlin definitions
 * `merlin.yaml` - The configuration file that is been created using `merlin init` command. This file represents one environment. The file contains information about how to talk to codefresh, how to talk to kubernetes cluster where the environment is set etc.
 * Environment - A set of components and components
 * Component - Logical part of the environment, has a set of operators.
 * Operator - Unit that describes how to interact with an Environment or a Component. Each operator has a name, in general, there is no uniqueness of a name across the Environment and nested Components

 ## Multiple operator execution flows
Running a `merlin run [NAME]` command will search the all operators named: NAME on the environment level and executed all of them sequentially (`merlin list` can help you to understand which operators have multiple executions).
A few notes to have in mind:
* If a flag `--component COMPONENT` set, operators from environment level and component level will be executed, where the environment level is in priority.
* Next operator executed only when the previous one finished successfully
* Operators executions flow are sharing environment variables exist on the host process (that runs the `merlin run` command)
