FROM ubuntu:latest

EXPOSE 8766:8766/udp
EXPOSE 16261:16261/udp 

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
  && apt-get install -y --no-install-recommends ca-certificates locales steamcmd \
  && rm -rf /var/lib/apt/lists/*

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
COPY ./fixtures/update_zomboid.txt /opt/pzserver/update_zomboid.txt

# Install project zomboid server
RUN steamcmd +runscript /opt/pzserver/update_zomboid.txt

USER pzuser

ENTRYPOINT ["/opt/pzserver/start-server.sh"]
