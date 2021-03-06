# zomboidBot
Discord Bot for managing a Project Zomboid dedicated server

## Running locally

### Prerequisites

Have a `.env.dev` file in the root of the repository set up to source, something like:
```
bot_oauth="bot-auth-here"
guild_id="test-server-id"
server_admin_password="idk-who-cares"
server_channel_id="test-channel-id"
server_name="servertest"
zomboid_cli_path="/opt/pzserver"
server_config_files_path="/home/pzuser/Zomboid/Server"
whitelisted_read_settings="PVP,PublicName,SafetySystem,SpawnPoint,SpawnItems,Password,PingLimit"
whitelisted_write_settings="PVP,SafetySystem,SpawnPoint,SpawnItems,PingLimit"
```

### Running the Bot + Server with Docker
Run the following commands:
```
make docker-build
docker run -p 16261:16261/udp -p 8766:8766/tcp -p 16262-16266:16262-16266/tcp zomboid-bot-image
```

The bot should now be set up like the real thing. The server itself won't be accessible to the internet unless you manually set up port forwarding on your router.

## Running Discord Bot locally
Run the following commands:
```
make build
./bin/zomboidBot
```

You should now have the bot running. Note that since it isn't running in the docker container this way, it won't be able to access the zomboid server within the docker container.

## Deploying

### Manually Deploy

Create a user for the bot to run with
```
sudo adduser pzuser
```

Then, clone this repo into the user's homedir (may need to install git)
```
cd /home/pzuser
sudo git clone https://github.com/Harlow-CS/zomboidBot.git
```

Next, create a `.env` file in the root of the cloned repo, structured the same as the aforementioned `.env.dev`
```
cd /home/pzuser/zomboidBot
sudo vim .env
```

Finally, navigate to the scripts folder of the repo and run the `lightsail.sh` script
```
cd /home/pzuser/zomboidBot/fixtures/scripts
./lightsail.sh
```

The environment will then install all of the remaining necesarry tools and run the bot daemon.
Your bot is now deployed!

Ensure the listed ports from the Dockerfile's `EXPOSE` commands are portforwarded to have a connectable dedicated server.

### Automagic Deploy
TODO

## Configure Server Settings

See `/fixtures/server` files, as well as `/bot/commands.go` for slash commands related to editing the server settings.


## More Resources

- https://pzwiki.net/wiki/Dedicated_Server
