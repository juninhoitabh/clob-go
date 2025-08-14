FROM golang:1.25.0-alpine3.22

ENV PATH="/go/bin:${PATH}"
ENV APP="/go/src/"
ENV GO111MODULE=on
ENV CGO_ENABLED=1

RUN apk update && apk upgrade --available && \
  apk --no-cache add zsh zsh-vcs git curl wget vim openssh postgresql-client

# RUN go install -v golang.org/x/tools/gopls@latest && \
#   go install -v golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest && \
#   go install -v honnef.co/go/tools/cmd/staticcheck@latest && \
#   go install -v github.com/go-delve/delve/cmd/dlv@latest && \
#   go install -v github.com/cweill/gotests/gotests@latest && \
#   go install -v github.com/fatih/gomodifytags@latest && \
#   go install -v github.com/josharian/impl@latest && \
#   go install -v github.com/haya14busa/goplay/cmd/goplay@latest && \
#   go install -v github.com/masahiro331/go-commitlinter@0.1.0 && \
#   go install -v github.com/automation-co/husky@latest && \
#   go install -v github.com/swaggo/swag/cmd/swag@latest

ENV HOME=/home/go
ENV ZSH_CUSTOM=${HOME}/.oh-my-zsh

WORKDIR ${HOME}

USER root

RUN sh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)" && \
  git clone https://github.com/denysdovhan/spaceship-prompt.git "${ZSH_CUSTOM}/themes/spaceship-prompt" && \
  ln -s "${ZSH_CUSTOM}/themes/spaceship-prompt/spaceship.zsh-theme" "${ZSH_CUSTOM}/themes/spaceship.zsh-theme" && \
  git clone https://github.com/zsh-users/zsh-syntax-highlighting "${ZSH_CUSTOM}/plugins/zsh-syntax-highlighting" && \
  git clone https://github.com/zsh-users/zsh-syntax-highlighting "${HOME}/zsh-syntax-highlighting" && \
  git clone https://github.com/zsh-users/zsh-completions "${ZSH_CUSTOM}/plugins/zsh-completions" && \
  git clone https://github.com/zsh-users/zsh-autosuggestions "${ZSH_CUSTOM}/plugins/zsh-autosuggestions"

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.49.0

COPY ./terminal/.zshrc ${HOME}/.zshrc

ENV DOCKERIZE_VERSION=v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
  && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
  && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

WORKDIR ${APP}

ARG GIT_USER_NAME
ARG GIT_USER_EMAIL

RUN git config --global --add safe.directory ${APP} \
  && git config --global user.email ${GIT_USER_EMAIL} \
  && git config --global user.name ${GIT_USER_NAME}

EXPOSE 3000
