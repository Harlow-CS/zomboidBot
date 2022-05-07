FROM ubuntu:latest

# Ports necessary for the server
EXPOSE 8766:8766/udp
EXPOSE 16261:16261/udp 

# Ports necessary for 5 players
EXPOSE 16262:16262/tcp 
EXPOSE 16263:16263/tcp 
EXPOSE 16264:16264/tcp 
EXPOSE 16265:16265/tcp 
EXPOSE 16266:16266/tcp 

# Set environment variables
ENV USER root
ENV HOME /root

# Set working directory
WORKDIR $HOME

# Insert Steam prompt answers
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN echo steam steam/question select "I AGREE" | debconf-set-selections \
  && echo steam steam/license note '' | debconf-set-selections

# Update the repository and install SteamCMD
ARG DEBIAN_FRONTEND=noninteractive
RUN dpkg --add-architecture i386 \
  && apt-get update -y \
  && apt-get install -y --no-install-recommends ca-certificates \
  locales steamcmd expect build-essential golang git &&\
  rm -rf /var/lib/apt/lists/*

# Add unicode support
RUN locale-gen en_US.UTF-8
ENV LANG 'en_US.UTF-8'
ENV LANGUAGE 'en_US:en'

# Create symlink for executable
RUN ln -s /usr/games/steamcmd /usr/bin/steamcmd

# Update SteamCMD and verify latest version
RUN steamcmd +quit

# Create project zomboid user and switch to it
RUN adduser pzuser
RUN mkdir /opt/pzserver &&\
  chown pzuser:pzuser /opt/pzserver

# Copy steam install script
COPY ./fixtures/steam/update_zomboid.txt /opt/pzserver/update_zomboid.txt
COPY . /zomboidBot

# Install project zomboid server
RUN steamcmd +runscript /opt/pzserver/update_zomboid.txt

# Copy zomboidBot service into systemd
RUN cp /zomboidBot/fixtures/systemd/zomboidBot.service /etc/systemd/system/zomboidBot.service

# Build ZomBot
WORKDIR /zomboidBot
RUN make build

# Swap to pzuser
USER pzuser

# First time server run
RUN source /zomboidBot/.env.dev && /zomboidBot/fixtures/scripts/first-time-server-start.sh $server_admin_password

# Start Discord Bot
ENTRYPOINT source /zomboidBot/.env.dev && /zomboidBot/bin/zomboidBot
