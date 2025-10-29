# pdal-kafka-project
Kafka based project for a Product Design and Implementation course to learn how to scale Kafka while keeping data in order.

# Application
The application consists of three main parts; producer, consumer + API and frontend. The Kafka parts of the application are all going to be written in Go, but having Kafka as a broker means that the project could very well be language agnostic. I could always make another producer in Python, Rust or C++. The frontend is going to be written in JS/TS, with some minimal frotend framework.

## Producer
The producer acts as a data generation machine. I create a worker pool inside the same application and the workers will take a job as a data generator. Each generator acts as a delivery vehicle, that output their coordinates to Kafka every X amount of seconds. In a real life scenario the vehicles would probably be given unique names whenever assigned, but in our case we're just randomly generating a Vehicle Identifier, which is a vehicle-UUIDv7 string. This way we could create more producers easily without having clashing names and have potentially thousands of fake vehicles across different locations, all producing data to the same Kafka cluster.

## Consumer + API
The Kafka consumer and API are going to be written in Golang as well. I will be using Websockets as the protocol for this side, since it allows persistent real-time communication for the clients. I'd hope that the data generation is realistic enough that we could actually see where our vehicles are located on a global map. There's really no need for a dedicated database at this stage, because there won't be user handling for the frontend and Kafka can retain data for a pretty long time.

## Frontend
The frontend is most likely going to be a Vite application, using Shadcn components to keep things simple. I don't really care for the frontend as long as I can visualize some data. The primary goal of this application is to learn Kafka.
