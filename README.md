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
attribute to check for their staleness.

## **Flags**
The following flags are configurable when launching UCRS:
- **--hostname (string)=** hostname for ucrs (default "0.0.0.0")
- **--port (int)=** port for ucrs (default 8080)
- **--redis-hostname (string)=** hostname for redis cache (default "0.0.0.0")
- **--redis-port (int)=** port for redis cache (default "0.0.0.0")
- **--topic (string)=** fcm topic for registration token subscription (default "un")

## **Topic**
UCRS functions by attaching all devices to topic specified by the **--topic** flag. By default,
the topic will be "un" which stands for universal notification. Regardless, when sending
notifications, ensure that notifications are being sent to the correct topic.

## Setup
1. Install Golang(1.16) onto the machine.
2. Setup Firebase and Redis as specified [here](#config).
3. Build the application with `go build`.
4. Launch the application with the appropriate flags via the `ucrs` binary.
