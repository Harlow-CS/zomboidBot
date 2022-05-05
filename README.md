# factorioBot
Discord Bot for managing factorio dedicated servers (WIP)

## Running locally with Docker
Create a `.env` file that source the necessary environment variables.

Example:
```
#!/bin/bash

export bot_oauth="token-value"
export factorio_cli_path="/opt/factorio"
export server_channel_id="channe-id"
export factorio_username="Factorio-Username"
export factorio_token="Factorio-Token"
```

Then run the following commands:
```
make docker-build
docker run -p 34197:34197/udp -p 27015:27015/tcp factorio-bot-image
```

After this, a server will be made available on the public matchmaking list!<br/>
If running from MacOS, you can only connect via directly the IP even if ports are properly forwarded.

## Configure Server Settings

See `/fixtures/server-settings.json`


## More Resources

- https://wiki.factorio.com/Command_line_parameters
