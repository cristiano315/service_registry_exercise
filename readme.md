This app is a simple example of a distributed, containerized architecture, including:

- A service registry, for service discovery
- A client which connects to the service registry
- 2 Sample services, one with state and one stateless.

The app is fully containerized and uses docker and docker compose. To run it, make sure to have docker and docker compose installed on your system. If you have both installed, you can run the app by navigating to the root directory and executing the command "docker compose up". Make sure to add the flag --build for the first time, so that docker can create the containers.
All services use the integrated Go RPC to communicate.
registry: a simple registry, which holds the addresses and names of the registered services. It has a register, deregister, and a discover procedure.
weather-service: a simple dummy stateless service. It responds with a string (in this case a fake weather forecast)
counter-service: a simple stateful service. It uses Redis to keep a synchronized counter on an external database, which is increased on every call and is shared with all the instances of the service. It returns the updated counter value.
Both services register to the registry on startup. They shut down and deregister after 5 minutes, to simulate a deregistration.
client: a simple client to interact with the services. It has a list with the service names and their calls. On startup, it asks the registry for all the services' addresses and keeps them stored locally. Since the servers for the processes may deregister, it repeats this procedure every minute (the time may be changed). It uses a round robin local load balancer for every service, and keeps calling each one every second.
All the networking is done through docker's internal network, so the nodes must all run on the same machine.
You can use docker compose logs -f to see the logs while the containers are running.
