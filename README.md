# **UCRS**
UCRS (universal client registration service), is a microservice that handles 
subscription and lifecycle events for registration tokens as part of FCM. UCRS
checks and removes stale tokens, which ensures that only active devices are configured
to recieve notifications.

## **Config**

### **Firebase**
To setup UCRS, one must first setup FCM service via Firebase. Create a new Firebase project,
and download the service account file as needed by the firebase andmin SDK. Then, the path
to this JSON file must be passed to the **GOOGLE_APPLICATION_CREDENTIALS** environment variable.

### **Redis**
Redis must be setup and running as it will be used for caching registration tokens along with TTL
attribute to check for their staleness. Tokens that are found to be stale will be automatically removed
from the specified FCM topic.

### **Environemnt Setup**
The following environment variables are configurable when launching UCRS:
- **UCRS_HOSTNAME (string)=** hostname for ucrs (default "127.0.0.1")
- **UCRS_PORT (int)=** port for ucrs (default 8080)
- **UCRS_REDIS_HOSTNAME (string)=** hostname for redis cache (default "127.0.0.1")
- **UCRS_REDIS_PORT (int)=** port for redis cache (default 6379)
- **UCRS_FCM_TOPIC (string)=** fcm topic for registration token subscription (default "un")

## **Topic**
UCRS functions by attaching all devices to topic specified by the **UCRS_FCM_TOPIC** env var. By default,
the topic will be "un" which stands for universal notification. Regardless, when sending
notifications, ensure that notifications are being sent to the correct topic.

## **Endpoints**
The following REST endpoints are available:
- `/ping` **[GET]** Ping pong styled endpoint for checking online status of UCRS.
- `/token/register` **[POST]** Register a FCM Registration token with UCRS. Token is given a
TTL decay value and is put into the FCM topic.
- `/token/{rt}/heartbeat` **[PUT]** - Heartbeat event takes place given the registration token passed via
the placeholder `{rt}`. This resets the TTL of the token within the registry.
- `/token/{rt}/update-rt` **[PUT]** - Update event takes place given the registration token passed via
the placeholder `{rt}`. The token specified by `{rt}` is removed from the registry and the FCM topic, and
the `new_token` is instead put into the registry and added to the FCM topic.
 

## **Setup**

### **Bare-Metal**
1. Install Golang(1.16) onto the machine.
2. Setup Firebase and Redis as specified [here](#config).
3. Build the application with `go build`.
4. Launch the application with the appropriate flags via the `ucrs` binary.

### **Docker**
1. Build the container `docker build -t ucrs:1.0 .` . Note, the [environment setup](#environemnt-setup)
options can be passed here as build args to make them built into
the iamge. This can be useful for [GOOGLE_APPLICATION_CREDENTIALS](#firebase) since removes
the need to pass it the firebase project config to all containers.
2. If build-time setup isnt used, the image exposes the `/fcm` volume for passing
firebase service file into the container. Likewise, environment variables for
[environment setup](#environemnt-setup) would make use of standard docker `-e` flag.
Example: 
```bash
docker run --name ucrs-1 --rm -v $PWD:/fcm -e UCRS_HOSTNAME=0.0.0.0 -e UCRS_PORT=8083 -e UCRS_REDIS_HOSTNAME=10.0.1.10 -e UCRS_REDIS_PORT=6389 -e GOOGLE_APPLICATION_CREDENTIALS=/fcm/upn-service-.json -p 8083:8083 -d ucrs:1.0
```