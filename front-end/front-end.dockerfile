FROM ubuntu:latest
LABEL authors="vovav"

ENTRYPOINT ["top", "-b"]