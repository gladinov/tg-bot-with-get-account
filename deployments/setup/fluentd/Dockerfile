FROM fluent/fluentd:v1.16.0-debian-1.0

USER root

RUN apt-get update
RUN apt-get install -y build-essential libgeoip-dev libmaxminddb-dev

# устанавливаем гемы по одному
RUN gem install "excon" 
RUN gem install "faraday" --no-document 
RUN gem install "faraday-net_http" --no-document 
RUN gem install "fluent-plugin-elasticsearch" --no-document 
RUN gem install "fluent-plugin-opensearch" --no-document 
RUN gem install "fluent-plugin-beats" --no-document
RUN gem install "fluent-plugin-docker" --no-document

USER fluent
