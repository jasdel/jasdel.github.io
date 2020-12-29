FROM public.ecr.aws/amazonlinux/amazonlinux:latest

RUN yum update -y && \
	yum install -y tar xz gzip git coreutils
#RUN curl -L -o nodejs.tar.xz https://nodejs.org/dist/v14.15.1/node-v14.15.1-linux-x64.tar.xz && \
#    mkdir -p /opt/nodejs && \
#    tar --strip-components=1 -xJf nodejs.tar.xz -C /opt/nodejs && rm -f nodejs.tar.xz

RUN curl -L -o hugo.tar.gz https://github.com/gohugoio/hugo/releases/download/v0.79.1/hugo_extended_0.79.1_Linux-64bit.tar.gz && \
    tar -xzf hugo.tar.gz hugo && \
    mv hugo /usr/local/bin && \
	rm -f hugo.tar.gz

ADD . /jasdel-docs

WORKDIR /jasdel-docs

#ENV PATH /opt/nodejs/bin:${PATH}
#RUN npm install && \

RUN HUGO_ENV=production hugo -d docs --gc
