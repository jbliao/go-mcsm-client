# go-mcsm-client
A cli tool for [MCSManager](https://github.com/MCSManager/MCSManager)

## Build
```
git clone https://github.com/jbliao/go-mcsm-client.git

go build -o mcsm cmd/mcsm/main.go

sudo install mcsm -o 755 /usr/local/bin
```

## Usage
### Setup Environment
```
export MCSM_URL=<Your MCSM panel url>
export MCSM_API_KEY=<Your API KEY>
export MCSM_SERVICE_UUID=<Your service uuid ("GUID" on the UI)>
export MCSM_INSTANCE_UUID=<Your instance uuid ("UUID" on the UI)>
```

### List file
```
$ mcsm instance listfile
...

$ mcsm instance listfile --path /bluemap
D       logs             4096
D       web              4096
F       minecraft-client-1.20.0.jar              23028278
F       pluginState.json                 245
F       resourceExtensions.zip           38205

$ mcsm instance listfile --path /bluemap --include .json
F       pluginState.json                 245
```

### Backup worlds
***Warning: this command will take several minutes depend on your worlds size***
```
$ mcsm instance backup
$ ls worlds*.zip
worlds_20231111.zip

$ mcsm instance backup -f worlds.zip
$ ls worlds.zip
worlds.zip

$ mcsm instance backup -f - > /tmp/worlds.zip
$ ls /tmp/worlds.zip
worlds.zip
```
