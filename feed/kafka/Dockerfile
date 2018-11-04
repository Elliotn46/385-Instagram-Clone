FROM openjdk:8-jre

ENV KAFKA_VERSION 1.0.0
ENV KAFKA_SCALA_VERSION 2.11
ENV KAFKA_ARCH "kafka_$KAFKA_SCALA_VERSION-$KAFKA_VERSION.tgz"
ENV KAFKA_HOME /opt/kafka

WORKDIR /opt

RUN apt-get update && \
apt-get install -y jq && \
wget -O - $(wget -qO- https://www.apache.org/dyn/closer.cgi\?as_json\=1\&path\=/kafka/$KAFKA_VERSION/$KAFKA_ARCH | jq --raw-output '.preferred')kafka/$KAFKA_VERSION/$KAFKA_ARCH | tar zxf - && \
mv /opt/kafka_$KAFKA_SCALA_VERSION-$KAFKA_VERSION $KAFKA_HOME && \
sed -i 's/zookeeper.connect=localhost:2181/zookeeper.connect=zookeeper:2181/g' /opt/kafka/config/server.properties && \
sed -i 's/broker.id=0/broker.id=-1/g' /opt/kafka/config/server.properties

CMD ["/opt/kafka/bin/kafka-server-start.sh", "/opt/kafka/config/server.properties"]
