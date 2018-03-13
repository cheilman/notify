# notify

A system for sending notifications to the user.  This contains a couple of components:

- A service that listens on a port, tracks notifications and history, etc.
    - This service should be running on a system that has a user sitting in front of it
    - The base implementation will tie to `notify-osd`
- A client that submits notifications to the service
    - Parameters should include:
        - A message
        - A category/sub-category (optional)
        - A severity (optional)
        - Some icons maybe? (optional)
    - The client should also possibly print the message out locally, in case this is running remotely (through ssh) or if there isn't a running service.

