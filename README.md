# Hello World published at Heroku

Link: https://mumabi.herokuapp.com

# Run kafka in docker container

## Setup

### Here is the image

```
---
version: '2'
services:
  zk1:
    image: confluentinc/cp-zookeeper:5.5.0
    container_name: zookeeper1
    hostname: zk1
    ports:
      - "22181:22181"
    environment:
      ZOOKEEPER_SERVER_ID: 1
      ZOOKEEPER_CLIENT_PORT: 22181
      ZOOKEEPER_TICK_TIME: 2000
      ZOOKEEPER_INIT_LIMIT: 5
      ZOOKEEPER_SYNC_LIMIT: 2
      ZOOKEEPER_SERVERS: zk1:22888:23888;zk2:32888:33888;zk3:42888:43888

  zk2:
    image: confluentinc/cp-zookeeper:5.5.0
    container_name: zookeeper2
    hostname: zk2
    ports:
      - "32181:32181"
    environment:
      ZOOKEEPER_SERVER_ID: 2
      ZOOKEEPER_CLIENT_PORT: 32181
      ZOOKEEPER_TICK_TIME: 2000
      ZOOKEEPER_INIT_LIMIT: 5
      ZOOKEEPER_SYNC_LIMIT: 2
      ZOOKEEPER_SERVERS: zk1:22888:23888;zk2:32888:33888;zk3:42888:43888

  zk3:
    image: confluentinc/cp-zookeeper:5.5.0
    container_name: zookeeper3
    hostname: zk3
    ports:
      - "42181:42181"
    environment:
      ZOOKEEPER_SERVER_ID: 3
      ZOOKEEPER_CLIENT_PORT: 42181
      ZOOKEEPER_TICK_TIME: 2000
      ZOOKEEPER_INIT_LIMIT: 5
      ZOOKEEPER_SYNC_LIMIT: 2
      ZOOKEEPER_SERVERS: zk1:22888:23888;zk2:32888:33888;zk3:42888:43888

  kafka-1:
    image: confluentinc/cp-kafka:5.5.0
    container_name: kafka1
    hostname: kafka-1
    ports:
      - "19092:19092"
    depends_on:
      - zk1
      - zk2
      - zk3
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zk1:22181,zk2:32181,zk3:42181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-1:19092

  kafka-2:
    image: confluentinc/cp-kafka:5.5.0
    container_name: kafka2
    hostname: kafka-2
    ports:
      - "29092:29092"
    depends_on:
      - zk1
      - zk2
      - zk3
    environment:
      KAFKA_BROKER_ID: 2
      KAFKA_ZOOKEEPER_CONNECT: zk1:22181,zk2:32181,zk3:42181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-2:29092

  kafka-3:
    image: confluentinc/cp-kafka:5.5.0
    container_name: kafka3
    hostname: kafka-3
    ports:
      - "39092:39092"
    depends_on:
      - zk1
      - zk2
      - zk3
    environment:
      KAFKA_BROKER_ID: 3
      KAFKA_ZOOKEEPER_CONNECT: zk1:22181,zk2:32181,zk3:42181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-3:39092

```

### run ZooKeeper and kafka cluster by executing the following command:

```console
docker-compose up -d
```

### You can run this command to verify that the services are up and running:

```console
docker-compose ps
```

### Update your /etc/host file

```console
127.0.0.1	localhost

# Your_Local_IP kafka-1 kafka-2 kafka-3
192.168.43.148 kafka-1 kafka-2 kafka-3


# The following lines are desirable for IPv6 capable hosts
::1     ip6-localhost ip6-loopback
```

### Additional file creation

Create a file credentials.env in the PubSub/consumer/ repo  
Contents of this file(write yours, without "")

```
EMAIL_PASS = your_Email_Password
smsAccountSid = your_sms_Account_Sid
smsAuthToken = your_sms_Auth_Token
smsFrom = your_mobile_number_(+countryCode + mobileNumbe)
```

### Install kafkacat for testing

```
sudo apt-get install kafkacat
```

### List all the brokers and topic in the cluster

```
kafkacat -L -b kafka-1:19092
```

## Usage

### Run producer service

```
cd producer
go run main.go
```

Running this program will create three topics "topic1","topic2","topic3" with 10 partition by default, you can change it in producer/config/config.json file

### Run consumer service (in separate terminls)

```
cd consumer

#for mail consumer group
go run main.go mail

#for sms consumer group
go run main.go sms

```

### Producing messages

Run the [postman](www.postman.com), in the `x-www-form-urlencoded` in `body` set the keys according to the given response format

```
  {
      "request_id": "1",
      "topic_name": "topic1",
      "message_body": "Transaction successful",
      "transaction_id": "90",
      "email": "your_email@host.com",
      "phone": "+919876543210",
      "customer_id": "1",
      "key": "key1"
  }
```

hit the send button, message has been appended to in the kafka queue.

### Consuming messages

Run the above given command to consume messages.

## Features

### Basic

-   Messages can be produced based on key or partition or in a least-byte manner. if the key is empty, then partition number allocation will be in action. if both empty, then the least-byte allocation will be served by the program.
-   ConsumerGroup can be subscribed to multiple topics, keeping this property in mind, we can specify more than one topic for consuming, and the subscription handle happens at the service bootup.
-   As soon as the new topic is typed in the JSON object, A new topic will be created with ten partitions by default.
-   SMS and Email both clients are in action with the proper fault-tolerant system. For Email, we used Gomail, and for SMS Gotwilio library is used.
-   If the client or the consumer service itself is down, no worries, messages are safe and will be consumed again when service or client is up.
-   Made extensive documentation for running these services on any system.
-   Email and SMS client credentials are safe and stored in the environment file.
-   Millions of messages are processed in a small instance of time. Thanks to Goroutines.
-   Writing and rotating all the logs and errors in the separate files, not in the terminal itself.

### Advance

-   Design principles and standard conventions of a typical Golang project are strictly followed.
-   Testing is easy to implement for such a well-designed code. Extensive reading of articles related to UT shared by the mentor was also done. If the time would have allowed us, this could add feather in the cup.
-   All the configuration constants are stored in the JSON file, and the Viper library is used to read the same.
-   Goroutines made our service tremendously fast and accurate that our service is able to process hundreds of thousands of messages in no time.

## Terminal Commands

### Create a topic

```
docker run \
  --net=host \
  --rm \
  confluentinc/cp-kafka:5.5.0 \
  kafka-topics --create --topic bar --partitions 3 --replication-factor 3 --if-not-exists --zookeeper localhost:32181

```

### Who is leader, who is follower

```
for i in 22181 32181 42181; do
   docker run --net=host --rm confluentinc/cp-zookeeper:5.0.0 bash -c "echo stat | nc localhost $i | grep Mode"
done
```

### Run producer side instance

Send message to the first node of the cluster

```
kafkacat -P -b kafka-1:19092 -t bar
```

### Run consumer side instance

Receive message from the third node of the cluster

```
kafkacat -C -b kafka-3:39092 -t bar
```

### To purge the topic

```
docker exec kafka1 kafka-configs --zookeeper zookeeper1:22181 --alter --entity-type topics --entity-name bar --add-config retention.ms=10

### dont forget set it back to normal once purging is done
docker exec kafka1 kafka-configs --zookeeper zookeeper1:22181 --alter --entity-type topics --entity-name bar --add-config retention.ms=86400000
```

### To delete the topic

```
docker run \
  --net=host \
  --rm \
  confluentinc/cp-kafka:5.5.0 \
  kafka-topics --delete --topic bar --zookeeper localhost:32181
```
