# NAU 
The side-project manager you didn't know you needed.


## Commands
NAU is built to be modular. Imagine a Makefile but for you computer. Currently is is aimed at managing your projects. Currently has these commands implemented

### Open
Opens a specific project. Currently calls VSCode.
```shell
nau open <project>
```
If `project` is specified attemps to open best match of all your projects. If it is not specified it prompts the user to choose which project to open. Projects are always listed according to most recently changed.

### Start
Loads a template and prompts the user for information in order to create a new instance.
```shell
nau new <template>
```
<p align="center">
<img alt="NAU demo" src="assets/new.gif" width="600" />
</p>

If `template` is specified will prompt the user for information in order to create a new project from the specified template. If it is not specified will prompt the user to choose which template to load.

### Archive
Cleans and compresses specific project. Moves to `Archives` directory.
```shell
nau archive <project>
```
If `project` is specified will run `make archive` before compressing and moving it to `Archives` directory. If it is not specified will prompt user to choose which one. Ordered in reverse order of last modified.

### Config
Shows current configuration.
```shell
nau config
```

