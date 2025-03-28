#!/bin/bash
# Set Kafka broker address (replace with your actual address)
KAFKA_BROKER="broker1:9092"

# Check if Kafka broker is available
until kafka-topics.sh --bootstrap-server $KAFKA_BROKER --list &>/dev/null; do
  echo "Kafka broker not ready yet, waiting..."
  sleep 5
done


# Create topics using confluent-kafka-utils
kafka-topics --bootstrap-server $KAFKA_BROKER --from-file /path/to/your/topics.yaml --create

echo "Topics creation completed."
