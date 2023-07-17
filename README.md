# gin-demo

Sample project to show how to build go based applications with GitHub action & docker and used to personal test.

## Build docker image

We have two different Dockerfile for different purposes:

- Dockerfile
  - Build with `scratch`, which will build a small sized docker image, which is default in the demo
- Dockerfile.distoless

  - Build with distroless based Docker which will offer a minimal attack surface.

    ```bash
    docker build -t gin-demo:test -f Dockerfile.distroless .
    ```
