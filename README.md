[![][workflow-badge]][workflow-link]

# fluent-bit-rabbitmq-output

Fluent-Bit go RabbitMQ output plugin

# Run Fluent Bit with the new plugin

```
$ sudo docker-compose up -d rabbitmq
# wait for rabbitmq to come up
$ sudo docker-compose up --build fluentbit
```

# Configuration Parameters

| **Key**        | **Description**                                                 |
|----------------|-----------------------------------------------------------------|
| RabbitHost     | The hostname of the Rabbit-MQ server                            |
| RabbitPort     | The port under which the Rabbit-MQ is reachable                 |
| RabbitUser     | The user of the Rabbit-MQ host                                  |
| RabbitPassword | The username of the user which connects to the Rabbit-MQ server |
| ExchangeName   | The exchange to which fluent-bit send its logs                  |
| ExchangeType   | The exchange-type                                               |

# License

 * MIT

[workflow-link]:    https://github.com/Rayleigh865/fluent-bit-rabbitmq-output/actions?query=workflow%3AGo
[workflow-badge]:   https://github.com/Rayleigh865/fluent-bit-rabbitmq-output/workflows/Go/badge.svg?branch=master