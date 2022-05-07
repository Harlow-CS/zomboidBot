# Insert Steam prompt answers
sudo echo steam steam/question select "I AGREE" | sudo debconf-set-selections \
  && sudo echo steam steam/license note '' | sudo debconf-set-selections

# Update the repository and install SteamCMD
export DEBIAN_FRONTEND="noninteractive"
sudo dpkg --add-architecture i386 \
  && sudo apt-get update -y \
  && sudo apt-get install -y --no-install-recommends ca-certificates \
  locales steamcmd expect build-essential golang git \
  && sudo rm -rf /var/lib/apt/lists/*

# Add unicode support
sudo locale-gen en_US.UTF-8
export LANG="en_US.UTF-8"
export LANGUAGE="en_US:en"

# Create symlink for executable
sudo ln -s /usr/games/steamcmd /usr/bin/steamcmd

# Update SteamCMD and verify latest version
sudo steamcmd +quit

# Create project zomboid user
sudo adduser pzuser
sudo sudo usermod -a -G sudo pzuser
sudo mkdir /opt/pzserver &&\
  sudo chown pzuser:pzuser /opt/pzserver

# Install project zomboid server
sudo steamcmd +runscript /zomboidBot/fixtures/steam/update_zomboid.txt

# Create systemd service
sudo cp /zomboidBot/fixtures/systemd/zomboidBot.service /etc/systemd/system/zomboidBot.service
sudo systemctl daemon-reload

# Build the bot
cd /zomboidBot
make build
cd /

# Run the zomboid server for the first time
sudo su - pzuser
source /zomboidBot/.env && /zomboidBot/fixtures/scripts/first-time-server-start.sh $server_admin_password

# Run Bot daemon
sudo systemctl start zomboidBot
