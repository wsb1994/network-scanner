# network-scanner
A network scanner that determines if MySQL is running, and extracts data about it. It even checks for weird ports or other shenanigans. 


To run the service, 

1. make sure docker is installed on the machine
2. install golang if it is not already installed
3. run ```docker-compose up``` to start mysql
4. run ```cd service && go mod tidy```
5. run ```go run main.go```

this service defaults to the defaults of ```localhost:3306```

to change ```host``` use ```export HOST_PORT=5000``` as an example
to change ```ip``` use ```export HOST_IP=5000``` as an example

This was tested against mysql in the default configuration provided in the docker-compose provided. 

To run in 🐳docker🐳 do the following.
1. make sure you have an active docker or podman installation, and navigate to the base directory of the project and run the following.
2. ```cd service && docker build -t network-scanner .```
3. ```cd ..```
4. ```docker-compose up &```

5. use ls to clear the environment of any docker-compose stuff, or simply open a new terminal and run your now compiled docker container with the next command.
6. ```docker run --network="host" network-scanner```