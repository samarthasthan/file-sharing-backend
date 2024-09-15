#!/bin/bash

kafka-topics.sh --bootstrap-server kafka:29092 --topic file-process-in --create --partitions 3 --replication-factor 1 --if-not-exists
kafka-topics.sh --bootstrap-server kafka:29092 --topic file-process-out --create --partitions 3 --replication-factor 1 --if-not-exists
kafka-topics.sh --bootstrap-server kafka:29092 --topic mail --create --partitions 3 --replication-factor 1 --if-not-exists