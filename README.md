# remindme
Remindme is a tool that helps you to remember commands that you use frequently, it's like a bookmark for your commands.
You can add, remove, search commands that you have saved.

## Requirements
To benefit from the copy to clipboard functionality, you will need to have installed one of the following "clipboard manager" on Linux.
```sh
xsel, xclip, wl-clipboard or Termux:API add-on for termux-clipboard-get/set
```


## To install

Use the binary from the releases page or build it yourself:

```sh
$ make install
```

This installs the `rmm` binary into `$GOBIN` (or `$GOPATH/bin`).

## Configuration

On first run, `rmm` asks which storage to use and writes the answers to `~/.config/remindme/config.yaml`, or `$XDG_CONFIG_HOME/remindme/config.yaml` when `XDG_CONFIG_HOME` is set and no config exists in `~/.config/remindme` yet:

- `yaml` (default): notes are stored in a YAML file, `~/.config/remindme/data.yaml` unless you choose another name.
- `mongo`: notes are stored in a MongoDB collection.

```yaml
storageType: mongo
mongo:
  host: localhost
  port: 27017
  database: notes
  collection: notes
```

Show the configuration in use with:

```sh
$ rmm config
```

### Running MongoDB locally

If you don't have a MongoDB server, you can start one with Docker:

```sh
$ docker run -d --name remindme-mongo -p 127.0.0.1:27017:27017 -v remindme-mongo:/data/db mongo:8.2
```

## To use

Remindme is a command line tool, you can use it by typing `rmm` followed by the command you want.
For every command you can use the help flag to get more information about the command.


### Help

```sh
$ rmm --help
```

### Add a command

Type `rmm add` and follow the prompts to add a command. A command and at least one tag are required.

```sh
❯ rmm add
✔ Command: k get pods -n cilium-system -l app.kubernetes.io/name=cilium-agent --no-headers | awk '{print $1}' | xargs -I {} kubectl delete pod {} -n cilium-system█
Description: delete all cilium agent pods
Tags: k8s,cilium
Note added successfully
```

Or pass everything as flags:

```sh
$ rmm add --command "kubectl get pods -A" --description "list all pods" --tags k8s
```

### List all commands

```sh
$ rmm list
```

### List all commands with a specific tag

```sh
$ rmm list --tags k8s
```

`--tags` accepts comma separated values (`--tags k8s,cilium`) and can be repeated.

### Show a single command

```sh
$ rmm list --id <id>
```

The command is also copied to the clipboard.

### List all tags

```sh
$ rmm list tags
```

### Search commands

Searches commands, descriptions and tags. Use `-c`, `-d` or `-t` to search only in commands, descriptions or tags.

```sh
$ rmm search cilium
$ rmm search -t k8s
```

### Delete a command
```sh
$ rmm rm --id <id>
```

### Delete all commands with a specific tag
```sh
$ rmm rm --tags k8s
```
