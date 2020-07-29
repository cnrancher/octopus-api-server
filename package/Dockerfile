FROM ubuntu:19.10
RUN apt-get update && apt-get install -y git curl unzip xz-utils tini && \
    useradd rancher && \
    mkdir -p /var/lib/rancher/octopus-ui  && \
    chown -R rancher /var/lib/rancher/octopus-ui /usr/local/bin

WORKDIR /var/lib/rancher/octopus-ui

ENV OCTOPUS_UI_VERSION v1.0.2
ENV OCTOPUS_UI_PATH /usr/share/rancher/octopus-ui
# Please update the api-ui-version in pkg/settings/settings.go when updating the version here.
ENV OCTOPUS_API_UI_VERSION 1.1.9

RUN mkdir -p /usr/share/rancher/octopus-ui && \
    cd /usr/share/rancher/octopus-ui && \
    curl -sL https://rancher-octopus.s3-accelerate.amazonaws.com/ui/dist/${OCTOPUS_UI_VERSION}.tar.gz | tar xvzf - --strip-components=2 && \
    mkdir -p /usr/share/rancher/octopus-ui/api-ui && \
    cd /usr/share/rancher/octopus-ui/api-ui && \
    curl -sL https://releases.rancher.com/api-ui/${OCTOPUS_API_UI_VERSION}.tar.gz | tar xvzf - --strip-components=1 && \
    cd /var/lib/rancher/octopus-ui

COPY package/entrypoint.sh bin/octopus-api-server /usr/bin/
RUN chmod +x /usr/bin/entrypoint.sh

VOLUME /var/lib/rancher/octopus-ui
ENTRYPOINT ["entrypoint.sh"]
