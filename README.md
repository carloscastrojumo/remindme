# remindme
Reminde is a tool that helps you to remember commands that you use frequently, it's like a bookmark for your commands.
You can add, remove, search commands that you have saved.

## To install

Use the binary from the releases page or build it yourself.

## To use

Remmindme is a command line tool, you can use it by typing `rmm` followed by the command you want.
For every command you can use the help flag to get more information about the command.

### Help

```sh
$ rmm --help
```

### Add a command

Type `rmm add` and follow the prompts to add a command.

```sh
❯ rmm add
✔ Command: k get pods -n cilium-system -l app.kubernetes.io/name=cilium-agent --no-headers | awk '{print $1}' | xargs -I {} kubectl delete pod {} -n cilium-system█
Description: delete all cilium agent pods
Tags: k8s,cilium
Note added successfully
```
### List all commands

```sh
$ rmm list
```

### List all commands with a specific tag

```sh
$ rmm list -t k8s
```

### Delete a command
```sh
$ rmm rm <id>
```

### Delete all commands with a specific tag
```sh
$ rmm rm -t k8s
```
