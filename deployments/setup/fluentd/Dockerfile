FROM fluent/fluentd:v1.16.1-debian-1.0

USER root

RUN apt-get update \
    && apt-get install -y build-essential libgeoip-dev libmaxminddb-dev \
    && rm -rf /var/lib/apt/lists/*

# устанавливаем гемы по одному
RUN gem install "excon:0.109.0" \
    && gem install "faraday:2.8.1" --no-document \
    && gem install "faraday-net_http:3.0.2" --no-document \
    && gem install "fluent-plugin-elasticsearch:6.0.0" --no-document \
    && gem install "fluent-plugin-opensearch:1.1.5" --no-document \
    && gem install "fluent-plugin-beats:1.1.0" --no-document \
    && gem install "fluent-plugin-docker" --no-document

USER fluent

# Fluentd слушает HTTP и forward
EXPOSE 24224 24224/udp 9880

# HEALTHCHECK --interval=10s --timeout=5s --start-period=30s --retries=5 \
#   CMD curl -f http://localhost:9880/health || exit 1