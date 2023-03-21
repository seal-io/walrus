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

# set user/group IDs
RUN set -eux; \
	groupadd -r seal --gid=1000 && \
	useradd -r -g seal --uid=1000 --home-dir=/var/lib/seal --shell=/bin/bash seal; \
    \
	mkdir -p /var/lib/seal && \
	chown -R seal:seal /var/lib/seal

# get postgres
VOLUME /var/lib/postgresql
RUN set -eux; \
	groupadd -r postgres --gid=999; \
	useradd -r -g postgres --uid=999 --home-dir=/var/lib/postgresql --shell=/bin/bash postgres; \
	mkdir -p /var/lib/postgresql; \
	chown -R postgres:postgres /var/lib/postgresql
RUN set -eux; \
    apt-get update && apt-get install -y --no-install-recommends \
      postgresql-14 \
    ; \
    rm -rf /var/lib/apt/lists/*; \
    \
    ln -s /usr/lib/postgresql/14/bin /usr/lib/postgresql/bin && \
    ln -s /etc/postgresql/14/main /etc/postgresql/main && \
    ln -s /var/lib/postgresql/14/main /var/lib/postgresql/main; \
    \
    ls -1 /usr/lib/postgresql/bin/ | xargs -I @ ln -sf /usr/lib/postgresql/bin/@ /usr/bin/@; \
    \
    gosu postgres /etc/init.d/postgresql start && \
      gosu postgres psql --command "CREATE USER root WITH SUPERUSER PASSWORD 'Root123';" && \
      gosu postgres createdb -O root seal && \
      gosu postgres echo "listen_addresses='*'" >> /etc/postgresql/main/postgresql.conf && \
      gosu postgres echo "local all  postgres            peer" > /etc/postgresql/main/pg_hba.conf && \
      gosu postgres echo "host  all  root  127.0.0.1/32  trust" >> /etc/postgresql/main/pg_hba.conf && \
      gosu postgres echo "host  all  root  ::1/128       trust" >> /etc/postgresql/main/pg_hba.conf && \
      gosu postgres echo "host  all  all   0.0.0.0/0     scram-sha-256" >> /etc/postgresql/main/pg_hba.conf && \
    gosu postgres /etc/init.d/postgresql stop

# get redis
VOLUME /var/lib/redis
RUN set -eux; \
	groupadd -r redis --gid=998; \
	useradd -r -g redis --uid=998 --home-dir=/var/lib/redis --shell=/bin/bash redis; \
	mkdir -p /var/lib/redis; \
	chown -R redis:redis /var/lib/redis
RUN set -eux; \
    apt-get update && apt-get install -y --no-install-recommends \
      redis \
    ; \
    rm -rf /var/lib/apt/lists/*

# get k3s
ENV K3S_IMAGE_DIR=/var/lib/seal/k3s/agent/images \
    K3S_MANIFEST_DIR=/var/lib/seal/k3s/server/manifests \
    K3S_RESOLV_CONF=/etc/resolv.conf
COPY --from=rancher/k3s:v1.22.7-k3s1 \
    /bin/blkid \
    /bin/charon \
    /bin/cni \
    /bin/conntrack \
    /bin/containerd \
    /bin/containerd-shim-runc-v2 \
    /bin/ethtool \
    /bin/ip \
    /bin/ipset \
    /bin/losetup \
    /bin/pigz \
    /bin/runc \
    /bin/swanctl \
    /bin/which \
    /bin/aux/xtables-legacy-multi \
    /usr/bin/
RUN set -eux; \
    ln -s /usr/bin/cni /usr/bin/bridge && \
    ln -s /usr/bin/cni /usr/bin/flannel && \
    ln -s /usr/bin/cni /usr/bin/host-local && \
    ln -s /usr/bin/cni /usr/bin/loopback && \
    ln -s /usr/bin/cni /usr/bin/portmap && \
    ln -s /usr/bin/containerd /usr/bin/crictl && \
    ln -s /usr/bin/containerd /usr/bin/ctr && \
    ln -s /usr/bin/containerd /usr/bin/k3s && \
    ln -s /usr/bin/containerd /usr/bin/k3s-agent && \
    ln -s /usr/bin/containerd /usr/bin/k3s-etcd-snapshot && \
    ln -s /usr/bin/containerd /usr/bin/k3s-server && \
    ln -s /usr/bin/containerd /usr/bin/kubectl && \
    ln -s /usr/bin/pigz /usr/bin/unpigz && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/iptables && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/iptables-save && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/iptables-restore && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/iptables-translate && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/ip6tables && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/ip6tables-save && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/ip6tables-restore && \
    ln -s /usr/bin/xtables-legacy-multi /usr/bin/ip6tables-translate; \
    \
    case "${TARGETARCH}" in \
      amd64|x86_64) \
        IMAGE_URL='https://github.com/k3s-io/k3s/releases/download/v1.22.7%2Bk3s1/k3s-airgap-images-amd64.tar.zst'; \
        IMAGE_SUM='51b342d888bc7ece1a56909435a25a67b377286717773fb363f2c88fcfe4e3d7'; \
        ;; \
      aarch64|arm64) \
        IMAGE_URL='https://github.com/k3s-io/k3s/releases/download/v1.22.7%2Bk3s1/k3s-airgap-images-arm64.tar.zst'; \
        IMAGE_SUM='8f948ed9635ecb420daae47afc246f578409b58e11bb19f55c84973681e6104e'; \
        ;; \
    esac; \
    \
    mkdir -p ${K3S_IMAGE_DIR} && \
    axel -n 16 -o ${K3S_IMAGE_DIR}/k3s-images.tar.zst ${IMAGE_URL} && \
    echo "${IMAGE_SUM} *${K3S_IMAGE_DIR}/k3s-images.tar.zst" | sha256sum -c -

# get casdoor
ENV BEEGO_CONFIG_PATH="/etc/casdoor/app.conf"
COPY --from=sealio/casdoor:v1.197.0-seal.3 \
    /casdoor \
    /usr/bin/

ARG SERVE_UI_INDEX="file:///var/lib/seal/ui"
ENV SERVER_SETTING_SERVE_URL="" \
    SERVER_SETTING_SERVE_UI_INDEX="${SERVE_UI_INDEX}"
EXPOSE 80 443
VOLUME /var/lib/seal
COPY /image/ /
COPY /build/server-${TARGETOS}-${TARGETARCH} /usr/bin/seal
CMD ["/entrypoint.sh", "seal", "--log-debug", "--log-verbosity=4"]
