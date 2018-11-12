
docker-compose up -d
docker exec feed_kafka_1 kafka/bin/kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 10 --topic instagram_cache

# docker rm -f `docker ps -aq`
# docker network ls -f "driver=bridge" | grep ' testnet ' > /dev/null || docker network create testnet
#
# docker run -d --name zookeeper --net testnet zookeeper:3.4.11
#
#
# docker build -t kafka kafka/
# docker run -d --name kafka --net testnet kafka
# docker exec kafka kafka/bin/kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 10 --topic instagram_cache
#
#
# docker run -d --name redis --net testnet redis:5.0.0
#
# docker build -t feed_api feed_api/
# docker run -d --name feedapi --net testnet feed_api
#
# docker build -t goconsumer cache_kafka_consumer/
# docker run -d --name go --net testnet goconsumer


# docker exec kafka kafka/bin/kafka-topics.sh --list --zookeeper zookeeper:2181
# docker exec kafka kafka/bin/kafka-run-class.sh kafka.tools.ConsumerOffsetChecker --broker-info --group test_group --topic instagram_cache --zookeeper kafka:2181
# docker exec kafka kafka/bin/kafka-consumer-groups.sh --bootstrap-server kafka:9092 --group instagram_cache --describe

# docker exec -it feed_kafka_1 /bin/bash
# bin/kafka-console-producer.sh --broker-list zookeeper:9092 --topic instagram_cache
