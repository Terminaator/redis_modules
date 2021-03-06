FROM redis:5.0.7
RUN apt-get update && apt-get install -y \
    apt-utils \
    build-essential \
    manpages-dev \
    iputils-ping
RUN gcc --version
RUN mkdir /modules
COPY modules /modules
WORKDIR /modules
RUN make