FROM ubuntu:latest

MAINTAINER Alex Levinson <alexl@uber.com>

ENV INITRD No
ENV LANGUAGE en_US.UTF-8
ENV LC_ALL en_US.UTF-8
ENV LANG en_US.UTF-8
ENV GOVERSION 1.10.3
ENV GOROOT /opt/go
ENV GOPATH /root/go
ENV GSCRIPT_REVISION v1

RUN DEBIAN_FRONTEND=noninteractive \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
      file build-essential vim nano wget curl sudo net-tools pwgen locales \
      git-core logrotate software-properties-common && \
    locale-gen en_US en_US.UTF-8 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN cd /opt && \
    wget https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/bin/go /usr/bin/ && \
    mkdir $GOPATH

RUN go get -d github.com/gen0cide/gscript/... && \
    cd $GOPATH/src/github.com/gen0cide/gscript && \
    git checkout $GSCRIPT_REVISION && \
    go get ./... && \
    cd cmd/gscript && \
    go install -i -a && \
    cd /root

ADD ps1.sh /etc/profile.d/Z1_PS1.sh

RUN chmod +x /etc/profile.d/Z1_PS1.sh && \
    echo "" >> /root/.bashrc && \
    echo "export GOPATH=/root/go" >> /root/.bashrc && \
    echo "export GOROOT=/opt/go" >> /root/.bashrc && \
    echo "GSCRIPT_REVISION=$GSCRIPT_REVISION" >> /root/.bashrc && \
    echo "source /etc/profile.d/Z1_PS1.sh" >> /root/.bashrc

RUN git clone https://github.com/scopatz/nanorc.git /usr/share/nano-syntax-highlighting/

ADD nanorc /etc/nanorc

RUN git clone --depth=1 https://github.com/amix/vimrc.git /opt/vim

ADD vimrc /root/.vimrc

RUN echo "autocmd BufNewFile,BufRead *.gs set syntax=javascript" >> /opt/vim/my_configs.vim
RUN echo "set tabstop=4" >> /opt/vim/my_configs.vim
RUN echo "set shiftwidth=4" >> /opt/vim/my_configs.vim
RUN echo "set expandtab" >> /opt/vim/my_configs.vim

VOLUME /root/share

ENV HOME /root
WORKDIR /root
