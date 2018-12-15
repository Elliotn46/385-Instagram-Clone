#! /usr/bin/env bash

repo=$(dirname $0)/..
project=$(gcloud config get-value project)

for cfg in $(find $repo -name '*.yaml')
do
	tmpcfg="/tmp/$(basename $cfg)"
	sed "s|ssuuuu-222721|${project}|g" <$cfg >$tmpcfg
	kubectl create -f $tmpcfg
done

CASS_POD_ID=`kubectl get pods | awk '/cassandra/ { print $1 }'`

until kubectl exec $CASS_POD_ID -- cqlsh -e "show version";
do
        echo "Cassandra not ready......"
	sleep 5
done

CASS_COMMANDS=(
	"CREATE KEYSPACE IF NOT EXISTS instagram WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor' : 3}"
	"CREATE TABLE IF NOT EXISTS instagram.user(user_id timeuuid, username text, email text, password text, PRIMARY KEY (user_id));"
	"CREATE TABLE IF NOT EXISTS instagram.subscription(user_id timeuuid, tag text, subscribe_date timestamp, PRIMARY KEY (tag, subscribe_date)) WITH CLUSTERING ORDER BY (subscribe_date DESC);"
	"CREATE TABLE IF NOT EXISTS instagram.follows(user_id timeuuid, follows timeuuid, follow_date timestamp, PRIMARY KEY (user_id, follow_date)) WITH CLUSTERING ORDER BY (follow_date DESC);"
	"CREATE TABLE IF NOT EXISTS instagram.followed_by(user_id timeuuid, followed_by timeuuid, follow_date timestamp, PRIMARY KEY (user_id, follow_date)) WITH CLUSTERING ORDER BY (follow_date DESC);"
	"CREATE TABLE IF NOT EXISTS instagram.user_post(user_id timeuuid, post_id timeuuid, tag text, caption text, PRIMARY KEY (user_id, post_id)) WITH CLUSTERING ORDER BY (post_id DESC);"
	"CREATE TABLE IF NOT EXISTS instagram.user_post_tag(user_id timeuuid, post_id timeuuid, tag text, caption text, PRIMARY KEY(tag, post_id)) WITH CLUSTERING ORDER BY (post_id DESC);"
	"CREATE TABLE IF NOT EXISTS instagram.user_post_comment(user_id timeuuid, post_id timeuuid, comment text, PRIMARY KEY (user_id, post_id)) WITH CLUSTERING ORDER BY (post_id DESC);"
	"CREATE TABLE IF NOT EXISTS instagram.user_post_timeline(user_id timeuuid, post_id timeuuid, monthyear text, caption text, PRIMARY KEY ((user_id, monthyear), post_id)) WITH CLUSTERING ORDER BY (post_id DESC);"
	)

for element in "${CASS_COMMANDS[@]}"
do
	kubectl exec $CASS_POD_ID -- cqlsh -e "$element"
done

echo "Exposing Deployment Search && Auth"
kubectl expose deployment accountauth accountsearch --port 3301 --type LoadBalancer

echo "Setting up autoscaler on auth && search"
kubectl autoscale deployment accountauth accountsearch --min=1 --max=5 --cpu-percent=50

kubectl get hpa
