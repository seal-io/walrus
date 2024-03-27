FROM --platform=$TARGETPLATFORM ubuntu:24.04
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

# get stats tool
RUN set -eux; \
    apt-get update && apt-get install -y --no-install-recommends \
      procps sysstat dstat \
      ifstat net-tools dnsutils \
    ; \
    rm -rf /var/lib/apt/lists/*

# get image tool
RUN set -exo pipefail; \
    case "${TARGETARCH}" in \
      amd64|x86_64) \
        CRANE_URL='https://github.com/google/go-containerregistry/releases/download/v0.19.0/go-containerregistry_Linux_x86_64.tar.gz'; \
        ;; \
      aarch64|arm64) \
        CRANE_URL='https://github.com/google/go-containerregistry/releases/download/v0.19.0/go-containerregistry_Linux_arm64.tar.gz'; \
        ;; \
      *) \
        echo "Unsupported arch: ${TARGETARCH}"; \
        exit 1; \
        ;; \
    esac; \
    \
    axel -n 16 -o /tmp/crane.tar.gz ${CRANE_URL} && \
      tar -xzf /tmp/crane.tar.gz --directory /usr/bin/ crane && \
      chmod a+x /usr/bin/crane && \
      rm -f /tmp/crane.tar.gz

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

# get k3s
VOLUME /var/lib/k3s
ENV K3S_DATA_DIR=/var/lib/k3s \
    K3S_RESOLV_CONF=/etc/resolv.conf \
    ETCDCTL_API=3
COPY --from=rancher/k3s:v1.27.11-k3s1 \
    /bin/blkid \
    /bin/cni \
    /bin/conntrack \
    /bin/containerd \
    /bin/containerd-shim-runc-v2 \
    /bin/ethtool \
    /bin/ip \
    /bin/ipset \
    /bin/k3s \
    /bin/losetup \
    /bin/pigz \
    /bin/runc \
    /bin/which \
    /bin/aux/xtables-legacy-multi \
    /usr/bin/
RUN set -eux; \
    ln -s /usr/bin/cni /usr/bin/bridge && \
    ln -s /usr/bin/cni /usr/bin/flannel && \
    ln -s /usr/bin/cni /usr/bin/host-local && \
    ln -s /usr/bin/cni /usr/bin/loopback && \
    ln -s /usr/bin/cni /usr/bin/portmap && \
    ln -s /usr/bin/k3s /usr/bin/crictl && \
    ln -s /usr/bin/k3s /usr/bin/ctr && \
    ln -s /usr/bin/k3s /usr/bin/k3s-agent && \
    ln -s /usr/bin/k3s /usr/bin/k3s-etcd-snapshot && \
    ln -s /usr/bin/k3s /usr/bin/k3s-server && \
    ln -s /usr/bin/k3s /usr/bin/kubectl && \
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
    IMAGE_URL="https://github.com/k3s-io/k3s/releases/download/v1.27.11%2Bk3s1/k3s-airgap-images-${TARGETARCH}.tar.zst"; \
    mkdir -p ${K3S_DATA_DIR}/agent/images && \
    axel -n 16 -o ${K3S_DATA_DIR}/agent/images/k3s-images.tar.zst ${IMAGE_URL}; \
    \
    ETCDCTL_VER="v3.5.11"; \
    ETCDCTL_URL="https://github.com/etcd-io/etcd/releases/download/${ETCDCTL_VER}/etcd-${ETCDCTL_VER}-linux-${TARGETARCH}.tar.gz"; \
    axel -n 16 -o /tmp/etcdctl.tar.gz ${ETCDCTL_URL} && \
    tar -xzvf /tmp/etcdctl.tar.gz --directory /usr/bin/ --strip-components=1 --no-same-owner etcd-${ETCDCTL_VER}-linux-${TARGETARCH}/etcdctl

ENV WALRUS_DATA_DIR="/var/run/walrus" \
    WALRUS_LID_DIR="/var/lib/walrus"

# get storage tools
## NB: please update the following files if changed.
## - hack/mirror/walrus-images.txt.
## - pkg/systemkubres/helmcharts.go.
ENV MINIO_IMAGE="bitnami/minio:2024.3.26-debian-12-r0" \
    MINIO_OS_SHELL_IMAGE="bitnami/os-shell:12-debian-12-r17"
RUN set -exo pipefail; \
    crane pull --platform=${TARGETPLATFORM} ${MINIO_IMAGE} ${K3S_DATA_DIR}/agent/images/minio.tar; \
    crane pull --platform=${TARGETPLATFORM} ${MINIO_OS_SHELL_IMAGE} ${K3S_DATA_DIR}/agent/images/minio-os-shell.tar
ADD https://charts.bitnami.com/bitnami/minio-14.1.3.tgz "${WALRUS_LID_DIR}/charts/minio.tgz"

# get deployer tools
## NB: please update the following files if changed.
## - hack/mirror/walrus-images.txt.
## - pkg/systemkubres/helmcharts.go.
ENV DEPLOYER_TERRAFORM_IMAGE="sealio/terraform-deployer:v1.5.7-seal.1" \
    DEPLOYER_OPENTOFU_IMAGE="sealio/opentofu-deployer:v1.6.0-seal.1" \
    DEPLOYER_TERRAFORM_NETWORK_MIRROR_IMAGE="sealio/hermitcrab:v0.1.3"
RUN set -exo pipefail; \
    crane pull --platform=${TARGETPLATFORM} ${DEPLOYER_TERRAFORM_IMAGE} ${K3S_DATA_DIR}/agent/images/deployer-terraform.tar; \
    crane pull --platform=${TARGETPLATFORM} ${DEPLOYER_OPENTOFU_IMAGE} ${K3S_DATA_DIR}/agent/images/deployer-opentofu.tar; \
    crane pull --platform=${TARGETPLATFORM} ${DEPLOYER_TERRAFORM_NETWORK_MIRROR_IMAGE} ${K3S_DATA_DIR}/agent/images/deployer-terraform-network-mirror.tar
ADD https://github.com/seal-io/helm-charts/releases/download/hermitcrab-0.1.3/hermitcrab-0.1.3.tgz "${WALRUS_LID_DIR}/charts/hermitcrab.tgz"

# get cli
COPY /build/walrus-cli* ${WALRUS_LID_DIR}/clis/

ARG SERVE_UI_URL="file://${WALRUS_LID_DIR}/ui"
ARG WALRUS_FILES_URL="file://${WALRUS_LID_DIR}/files"
ENV WALRUS_SETTING_DEPLOYER_IMAGE="${DEPLOYER_TERRAFORM_IMAGE}" \
    WALRUS_SETTING_SERVE_UI_URL="${SERVE_UI_URL}" \
    WALRUS_SETTING_WALRUS_FILES_URL="${WALRUS_FILES_URL}"
EXPOSE 80 443
VOLUME ${WALRUS_DATA_DIR}
COPY /image/ /
COPY /build/walrus-manager-${TARGETOS}-${TARGETARCH} /usr/bin/walrus-manager
COPY /build/walrus-server-${TARGETOS}-${TARGETARCH} /usr/bin/walrus-server
ENV _RUNNING_INSIDE_CONTAINER_="true"
CMD ["walrus-server"]
