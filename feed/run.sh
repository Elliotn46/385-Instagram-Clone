docker rm -f `docker ps -aq`
docker network ls -f "driver=bridge" | grep ' testnet ' > /dev/null || docker network create testnet

docker run -d --name zookeeper --net testnet zookeeper:3.4.11


docker build -t kafka kafka/
docker run -d --name kafka --net testnet kafka
docker exec kafka kafka/bin/kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 10 --topic instagram_cache


docker run -d --name redis --net testnet redis:5.0.0

docker build -t feed_api feed_api/
docker run -d --name feedapi --net testnet feed_api

docker build -t goconsumer cache_kafka_consumer/
docker run -d --name go --net testnet goconsumer
