# FolderActions
A simple tool for linux which allows to run scripts on fs changes

# Instalation
`go get -u`<br>
`go build -o main`<br>
`sudo mv main /bin/fslistener`<br>
`sudo chmod +x /bin/fslistener`<br>

### Autocompletion
##### ZSH
fslistener --completion-script-zsh >> ~/.zshrc 


# Usage
```
usage: fslistener [<flags>]

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -d, --dir=DIR ...        The list of directories to watch
  -v, --verbose            Increase debug message
  -q, --quiet              Don't output messages
  -a, --allow-unsafe       Allow special characters in path and filename
  -s, --scripdir=SCRIPDIR  The folder containing the scripts
```

Example: 
```bash
fslistener --scripdir `pwd`/scripts --dir /dev/ --dir /tmp/ -v
```
If a file inside /dev/ or /tmp/ gets created or deleted, a script inside the ./scripts folder will be called.<br>
Let's assume the file /dev/sdc gets created (because you inserted a USB drive), the following script will be called:<br>
`./scripts/_dev_create.sh` with the argument `/dev/sdc`

### Env
`FA_SCRIPT_PATH` - sets the default path containing the scripts<br>
-s/--scriptdir overwrites this value
<br>
### Scripts
If a file in a folder changes (gets created or removed) a script is called by this tool. The path/name of this script is:<br>
`$FA_SCRIPT_PATH`/_real_path_of_changed_file_from_root_MODE.sh<br><br>
For example if a file named sda inside the dir /dev/ was created, following script would be called:<br>
`$FA_SCRIPT_PATH`/_dev_create.sh<br>
The first argument is the name of the changed file
