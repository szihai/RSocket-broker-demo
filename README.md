# RSocke broker demo

## Introduction
Some time ago, I showed another [demo](https://github.com/szihai/broker-flat) about one of the broker use cases. In this demo, we'll examine what's inside the broker. This broker is a prototype built on RSocket [golang SDK](https://github.com/rsocket/rsocket-go).
## What's different about RSocket
Other technologies use broker as well. For example, Kafka requires the broker deployment model. So, how is RSocket broker different from the others?   
First, RSocket doesn't mandate the broker model. It is for the simplicity of devop tasks.  
The current microservices orchastration tools make the networking unnecessarily complicated for the devleopers. The developers should not worry about which ip or port to use. All they should focus on is the application level networking: service A wants to talks to service B. With point to point connections, the burden is on the developers to figure out the IP and port. If, however, we let devops deploy the brokers as infrastructure and applications will only need to connect to a known IP:port, things will be much easier. Moreover, the broker not only shields the networking complecity from the developers, but it also provides a place outside of the applications to add policies or controls or logging features. I would say it enables the combination of control plane and data plane with similar functionality as service mesh.

In our example we use only one broker. All endpoints connect to the broker.  Both the service and the consumer are clients to the broker. The only difference is a service has to register itself to the broker.
![image](broker.png)

With this architecture, we can explore more of what RSocket enables:
### RSocket is connection oriented

In the broker, when an endpoint connects it will first send a setup frame to register. In our example, we use the `metadataRegister` to save the information. And when an endpoint drops the connection it will trigger the `OnClose` function. Thus, we don't need additional health check for endpoints. This is huge! Think about how much effort is spent on health checking especially in a large scale cluster. Also we will be able to save the 3rd party service registry and discovery. Yes, with the load balancing feature offered by the SDK and broker, we save a lot of efforts in connection pool management.

### RSocket is bi-directional

In our example, the title "consumer" and "service" don't mean anything. If the consumer allow, it will serve the service as well. And, pay attention to the code, they are using the same socket connection. What that means is not only are RSocket multiplex, but it is bi-directional as well. As a result, even for the broker, there will be way less connections to manage compared to other brokers for the same amount of endpoints.

### RSocket ensures back pressure end to end

Assume the client cannot process enough of the information. If the client now sends `request(N)` to the broker, the broker will relay the request to the service. Without this capability, the broker has to use buffer to hold the data. Now that the back pressure can reach the sender, the sender can simply slow down with the production. 

## How to use the demo
This demo is straightforward to try out. Just `git clone` the repo and run `go build` in each folder. Then you can start the broker first. After that, the provider and consumer. 
