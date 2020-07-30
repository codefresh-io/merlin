# Merlin

<img src="https://github.com/codefresh-io/merlin/blob/master/Merlin.png?raw=true" width="200" align="right">


## Installation
* MacOS:
  * `brew tap codefresh-io/merlin`
  * `brew install merlin`

* Prerequisite:
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
