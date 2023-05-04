# NAU 
The side-project manager you didn't know you needed.
<p align="center">
<img alt="NAU demo" src="assets/home.png" width="600" />
</p>

# Installation

```bash
# macOS or Linux
brew tap antonio-leitao/taps
brew install nau
```
# Getting Started 
NAU allows you to store permanent variables such as your name, email and website for automatically adding in projects.
It also requires you to supply a directory where to store your Projects, where your Templates are and where you want to store your Archives.
You can setup all variables by running
```shell
nau config
```
<p align="center">
<img alt="NAU demo" src="assets/config.gif" width="600" />
</p>

Otherwise you can simply change specific fields aat any time by running:

```shell
nau config <field> <value>
```
Below are all available fields with their explanations:
- `AUTHOR`: Your name
- `WEBSITE`: Your website 
- `EMAIL`: Your email 
- `REMOTE`: Url for your remote repositories. This is accessible if you want to imediatly add a remote directory to new projects on initialization.
- `BASE_COLOR`: Hex value for NAU's ui
- `EDITOR`: command to run for starting your preffered editor (code, nvim, vim...)
- `PROJECTS_PATH`: Path to your projects folder. Your root path is appended at the begining of the string you supply.
- `TEMPLATES_PATH`: Path to where your templates reside (see more in the Templates section)
- `ARCHIVES_PATH`: Where should NAU place archived projects (see more on the `archive` command)
You can also manually configure NAU through the file `.naurc` in your root directory
To view your current configuration just run

```shell
nau
```
# Commands

NAU is built to be modular. Imagine a Makefile but for you computer. Is is aimed at managing your projects. Currently has these commands implemented

### Open 
Opens a specific project using specified `editor`, default is `neovim`.

> **Note**
> `nau` will run `make open` command on the directory. You can add/edit this makefile command to enable further customization. 

```shell
nau open <project>
```
If `project` is specified attemps to open best match of all your projects. If it is not specified it prompts the user to choose which project to open. Projects are always listed according to most recently changed.
<p align="center">
<img alt="NAU demo" src="assets/open_project.gif" width="600" />
</p>

### New
Creates a new project either empty or from a prebuilt `nau` template. 
If `template` is specified will prompt the user for information in order to create a new project from the specified template. If it is not specified will prompt the user to choose which template to load.

> **Note**
> `nau` will run `make new` command on the directtory after the template is collapsed. You can add/edit this makefile command to enable further customization. 

```shell
nau new <template>
```
<p align="center">
<img alt="NAU demo" src="assets/new.gif" width="600" />
</p>

### Memo
Stores and organizes miscellaneous memos.

```shell
nau memo 
```
Creates a new memo.

```shell
nau memos <query>
```
Allows you to see iteractivelly all your memos. You can also select, filter and delete.
If a `query` starts the list already with the filter applied.

### Archive
Cleans and compresses specific project. Moves to `Archives` directory.

> **Note**
> `nau` will run `make archive` command on the directory before the project is archived. You can add/edit this makefile command to enable further customization. 

```shell
nau archive <project>
```
If `project` is specified will run `make archive` before compressing and moving it to `Archives` directory. If it is not specified will prompt user to choose which one. Ordered in reverse order of last modified.

# Templates
Nau relies on understanding what type are your projects. Each project either comes from a template or it doesnt.
### Template Direcory
Check the [`nau-templates`] repository for an example of a direcotry of templates. The template directory is stored in `config.Templates_path` and should look like this:
```text
templates
│   Python_#3776AB
│   Rust_#B2292D
│   Web_#F7DF1E
│   PascalCase_#000000
```
The supplied colors are going to be used by NAU thorughout the commands.This will result in the following projects directory:
```text
projects
└───Python
│   │   FRS_FirstProject
│   │   MPR_MyProject
└───Rust
│   │   IDX_RustProject
```
### `.nau` file
Each template should have a now file that specified which files are templated and have to be collapsed.

### `make init`
After collapsing the template `nau` will attempt to run the Makefile command `make init` this is usefull to initiate any template specific environments such as node or conda.
