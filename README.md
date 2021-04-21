# epiphany-wrapper-poc
PoC application to wrap containerised modules of epiphany

## example run

### general 

#### e help 

You can use `--help` switch anyway in a cli sub commands path. 

```shell
> e --help                          
E wrapper allows to interact with epiphany

...
```

### module sub-command

#### e module help

```shell
> e module --help
TODO

Usage:
  e module [command]

Available Commands:
  info        shows ifo of named module
  install     installs module into currently used environment
  search      searches for named module

Flags:
  -h, --help   help for module

Global Flags:
      --configDir string   config directory (default is .e)
      --logLevel string    log level (default is warn, values: [trace, debug, info, error, fatal])

Use "e module [command] --help" for more information about a command.
```

#### e module search 

```shell
> e module search azbi
epiphany-platform-modules/azbi:dev
```

#### e module info

```shell
> e module info epiphany-platform-modules/azbi:dev
  Component Version:
   Version: dev
   Image: docker.io/epiphanyplatform/azbi:dev
    Command:
     Name metadata
     Description meta
    Command:
     Name init
     Description init
    Command:
     Name plan
     Description plan
    Command:
     Name apply
     Description apply
    Command:
     Name plan-destroy
     Description plan destroy
    Command:
     Name destroy
     Description destroy
```

#### e module install

```shell
> e module install epiphany-platform-modules/azbi:dev
Installed module azbi:dev to environment 210416-1214
```

### environments sub-command

#### e environments help

```shell
> e environments --help                                       
TODO

Usage:
  e environments [command]

Aliases:
  environments, env

Available Commands:
  export      Exports an environment as a zip archive
  import      Imports a zip compressed environment
  info        Displays information about currently selected environment
  list        TODO
  new         Creates new environment
  run         Runs installed component command in environment
  use         Allows to select environment to be used

Flags:
  -h, --help   help for environments

Global Flags:
      --configDir string   config directory (default is .e)
      --logLevel string    log level (default is warn, values: [trace, debug, info, error, fatal])

Use "e environments [command] --help" for more information about a command.
```

#### e environments new

```shell
> e environments new e1
```

No output is expected.

#### e environments info

Here used after command `e components install c1`

```shell
> e environments info
Environment info:
 Name: e1
 UUID: 1dd02223-66ab-482f-8b0b-6c9e5d154c34
  Installed Component:
   Name: azbi
   Type: docker
   Version: dev
   Image: docker.io/epiphanyplatform/azbi:dev
    Command:
     Name metadata
     Description meta
    Command:
     Name init
     Description init
    Command:
     Name plan
     Description plan
    Command:
     Name apply
     Description apply
    Command:
     Name plan-destroy
     Description plan destroy
    Command:
     Name destroy
     Description destroy
```

#### e environments use

```shell 
> e environments use 
Use the arrow keys to navigate: ↓ ↑ → ← 
? Environments: 
  ▸ 210416-1214 (1dd02223-66ab-482f-8b0b-6c9e5d154c34, current)
    e1 (f1a114d9-9a38-4c36-8488-4272ed94f359)
```

#### e environments run

```shell
> e environments run azbi metadata
{"labels":{"kind":"infrastructure","name":"Azure Basic Infrastructure","provider":"azure","provides-pubips":true,"provides-vms":true,"short":"azbi","version":"dev"}}
```

### repos sub-command

#### e repos help

```shell
> e repos --help                                               
Commands related to repos management

Usage:
  e repos [command]

Available Commands:
  install     installs new repository
  list        Lists installed repositories

Flags:
  -h, --help   help for repos

Global Flags:
      --configDir string   config directory (default is .e)
      --logLevel string    log level (default is warn, values: [trace, debug, info, error, fatal])

Use "e repos [command] --help" for more information about a command.
```

#### e repos list

```shell
> e repos list  
Repository: epiphany-platform-modules
	Module: terraform:0.1.0
	Module: azbi:dev
```

#### e repos install

```shell
> e repos install mkyc/my-epiphany-repo
```

No output is expected but `list` output should get longer.

```shell
> e repos list                         
Repository: epiphany-platform-modules
	Module: terraform:0.1.0
	Module: azbi:dev
Repository: mkyc-my-epiphany-repo
	Module: azbi:0.1.0
	Module: azks:0.1.0
```

## configuration directory structure

After all command executed in previous section directory structure looks in similar way to: 

```shell
> tree ~/.e                                                   
/Users/mateusz/.e
├── config.yaml
├── environments
│   ├── 63fdee7b-cf31-46f9-be9b-61fad761b484
│   │   ├── azbi
│   │   │   └── dev
│   │   │       ├── mounts
│   │   │       └── runs
│   │   ├── config.yaml
│   │   └── shared
│   └── d94a07eb-0501-4beb-8bee-cc74e5b0d403
│       ├── config.yaml
│       └── shared
├── repos
│   └── epiphany-platform-modules.yaml
└── tmp

11 directories, 4 files
```

Main config file contains: 

```shell
> cat ~/.e/config.yaml 
version: v1
kind: Config
current-environment: 63fdee7b-cf31-46f9-be9b-61fad761b484
```

Used environment config file contains: 

```shell
> cat ~/.e/environments/63fdee7b-cf31-46f9-be9b-61fad761b484/config.yaml 
name: e1
uuid: 63fdee7b-cf31-46f9-be9b-61fad761b484
installed:
- environment_ref: 63fdee7b-cf31-46f9-be9b-61fad761b484
  name: azbi
  type: docker
  version: dev
  image: docker.io/epiphanyplatform/azbi:dev
  workdir: /shared
  mounts: []
  shared: /shared
  commands:
  - name: metadata
    description: meta
    command: metadata
    envs: {}
    args:
    - --json

...
```

## TODO

There is a lot TODO in a code which should be fixed