FROM ghcr.io/linuxserver/baseimage-kasmvnc:ubuntujammy-version-4b14a364

# set version label
ARG BUILD_DATE
ARG VERSION
LABEL build_version="Metatrader Docker:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="gdalia"

ENV TITLE=Metatrader5
ENV WINEARCH=win64
ENV WINEPREFIX="/config/.wine"
ENV DISPLAY=:0

#RUN mkdir -p /config/.wine && \
#    chown -R abc:abc /config/.wine && \
#    chmod -R 755 /config/.wine

# Update package lists and upgrade packages
RUN apt-get update && apt-get upgrade -y

# Install required packages
RUN apt-get install -y \
    python3-pip \
    wget \
    dos2unix \
    python3-pyxdg \
    netcat \
    && pip3 install --upgrade pip

# Add i386 architecture and update package lists
RUN dpkg --add-architecture i386 \
    && apt-get update

# Copy the scripts directory and convert start.sh to Unix format
COPY Metatrader /Metatrader
RUN dos2unix /Metatrader/*.sh && \
    chmod +x /Metatrader/*.sh

RUN /Metatrader/install.sh
#ENTRYPOINT /Metatrader/start.sh

COPY /root /
RUN touch /var/log/mt5_setup.log && \
    chown abc:abc /var/log/mt5_setup.log && \
    chmod 644 /var/log/mt5_setup.log

RUN mkdir -p /config/.wine && \
    chown -R abc:abc /config/.wine && \
    chmod -R 755 /config/.wine


EXPOSE 3000 8001
#VOLUME /config