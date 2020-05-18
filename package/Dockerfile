FROM ubuntu:19.10
RUN apt-get update && apt-get install -y git curl unzip xz-utils tini && \
    useradd rancher && \
    mkdir -p /var/lib/rancher/edge-ui  && \
    chown -R rancher /var/lib/rancher/edge-ui /usr/local/bin

WORKDIR /var/lib/rancher/edge-ui

ENV EDGE_UI_VERSION 0.1.0
ENV EDGE_UI_PATH /usr/share/rancher/edge-ui
# Please update the api-ui-version in pkg/settings/settings.go when updating the version here.
ENV EDGE_API_UI_VERSION 1.1.9

RUN mkdir -p /usr/share/rancher/edge-ui && \
    cd /usr/share/rancher/edge-ui && \
    curl -sL https://rancheredge.s3-ap-northeast-2.amazonaws.com/ui/${EDGE_UI_VERSION}.tar.gz | tar xvzf - --strip-components=1 && \
    mkdir -p /usr/share/rancher/edge-ui/api-ui && \
    cd /usr/share/rancher/edge-ui/api-ui && \
    curl -sL https://releases.rancher.com/api-ui/${EDGE_API_UI_VERSION}.tar.gz | tar xvzf - --strip-components=1 && \
    cd /var/lib/rancher/edge-ui

COPY package/entrypoint.sh bin/edge-api-server /usr/bin/
RUN chmod +x /usr/bin/entrypoint.sh

VOLUME /var/lib/rancher/edge-ui
ENTRYPOINT ["entrypoint.sh"]
