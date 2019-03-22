FROM ubuntu

RUN apt-get update >/dev/null && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y make build-essential \
    libssl-dev zlib1g-dev libbz2-dev \
    libreadline-dev libsqlite3-dev \
    wget curl llvm libncurses5-dev libncursesw5-dev \
    xz-utils tk-dev libffi-dev liblzma-dev python-openssl git

RUN curl https://pyenv.run | bash

ENV PATH="/root/.pyenv/bin:$PATH"
RUN eval "$(pyenv init -)"
RUN eval "$(pyenv virtualenv-init -)"

RUN pyenv install 3.6.6
