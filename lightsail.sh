# Insert Steam prompt answers
sudo echo steam steam/question select "I AGREE" | sudo debconf-set-selections \
  && sudo echo steam steam/license note '' | sudo debconf-set-selections

# Update the repository and install SteamCMD
export DEBIAN_FRONTEND="noninteractive"
sudo dpkg --add-architecture i386 \
  && sudo apt-get update -y \
  && sudo apt-get install -y --no-install-recommends ca-certificates locales steamcmd \
  && sudo rm -rf /var/lib/apt/lists/*

# Add unicode support
sudo locale-gen en_US.UTF-8
export LANG="en_US.UTF-8"
export LANGUAGE="en_US:en"

# Create symlink for executable
sudo ln -s /usr/games/steamcmd /usr/bin/steamcmd

# Update SteamCMD and verify latest version
sudo steamcmd +quit

# Create project zomboid user and switch to it
sudo adduser pzuser
sudo mkdir /opt/pzserver &&\
  sudo chown pzuser:pzuser /opt/pzserver

# Copy steam install script
# sudo cp ./fixtures/update_zomboid.txt /opt/pzserver/update_zomboid.txt

# Install project zomboid server
sudo steamcmd +runscript ./fixtures/update_zomboid.txt

