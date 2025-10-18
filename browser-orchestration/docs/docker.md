# SETUP DOCKER
## Make the DockerFile.Selinium and add this inside
```bash
FROM selenium/standalone-chrome-debug:latest

USER root

# Install ffmpeg and dependencies
RUN apt-get update && \
    apt-get install -y ffmpeg && \
    rm -rf /var/lib/apt/lists/*

USER seluser
```
Make the docker container using DockerFile.Selinium.Now, your docker container has chrome and ffmpeg inside

## Why we require it?

standalone-chrome-debug will roll up chrome for you. We will use ffmpeg for recording inside the docker container. Add your bind address inside chrome.go or any other browser.go which i have added. Your video will be added there