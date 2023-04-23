FROM --platform=$TARGETPLATFORM ubuntu:22.04
SHELL ["/bin/bash", "-c"]

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

ENV DEBIAN_FRONTEND=noninteractive
RUN set -eux; \
    apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates \
      git-lfs git \
      unzip xz-utils \
      curl axel wget \
      vim \
      gosu \
    ; \
    rm -rf /var/lib/apt/lists/*

# set locale
RUN set -eux; \
	apt-get update && apt-get install -y --no-install-recommends \
      locales \
    ; \
    rm -rf /var/lib/apt/lists/*; \
    \
	localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
ENV LANG='en_US.UTF-8' \
    LANGUAGE='en_US:en' \
    LC_ALL='en_US.UTF-8'

VOLUME /var/run/seal
COPY /image/ /
COPY /build/agent-${TARGETOS}-${TARGETARCH} /usr/bin/seal
CMD ["seal", "--log-debug", "--log-verbosity=4"]
