# dashport

Dashport makes it possible to clone dashboards between Dynatrace tenants.

### Building

To build the project:  
```
cd cmd/dashport
go build -o $HOME/bin/dashport
```

### Usage

You need to create api keys on both prod and dev tenants before you can use dashport. 

Dev link to create apitoken: ```https://<dynatrace dev tenant>.live.dynatrace.com/#settings/integration/apikeys;gf=all```  
Prod link to create apitoken: ```https://<dynatrace prod tenant>.live.dynatrace.com/#settings/integration/apikeys;gf=all```

Ensure the token has read and write configuration enabled.

Add the details to $HOME/.config/dashport/dashportcfg.json, like so:  
```
mkdir $HOME/.config/dashport
vim $HOME/.config/dashport/dashportcfg.json
{
"apiConfig": {
	"Tenants": [
		{
		"env": "dev",
		"token": "<the token you created>",
		"url": "https://<dynatrace dev tenant>.live.dynatrace.com/api/config/v1"
		},
		{
		"env": "prod",
		"token": "<the token you created>",
		"url": "https://<dynatrace prod tenant>.live.dynatrace.com/api/config/v1"
		}
		]
	}
	
}
```

Save the file.

#### Usage examples

To print all dashboards for prod:  

```
dashport -act printall -oenv prod 
```

Get specific dashboard id by name:

```
dashport -act printall -oenv dev | grep -B1 clone-test | awk -F'"id": "|"' '/id/ { print $2 }'
```

Using jq:  

```
dashport -act printall -oenv dev | 'jq '.[] | .[] | select (.name=="clone-test") | .id'
```

To print specific dashboard in prod:  

```
dashport -act print -oenv prod -id <id>
```

To clone specific dashboard from prod to dev:  

```
dashport -act clone -oenv prod -denv dev -id <id>
```

To update a specific dashboard from prod to dev:  

```
dashport -act update -oenv prod -id <id> -denv dev -did <destination id>
```

To delete specific dashboard on tenant env:  

```
dashport -act delete -oenv dev -id <id>
```

And a simple script to backup all dashboards:  

```
#!/bin/bash

set -eu -o pipefail

tenant='dev'

for dash in $(dashport -act printall -oenv "$tenant" | jq '.[] | .[] | .id'); do
  dashport dashport -act print -oenv "$tenant" -id "$dash" > "$dash".json
done
```

### Disclaimer

This is software is not well tested and bugs may occur. Save copies of dashboards before experimenting with cloning etc. Report bugs.
