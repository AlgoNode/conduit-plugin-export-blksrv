#!/bin/bash

#static 3MB build, no distro, no shell :)
docker build . -t urtho/conduit-export-blksrv:latest
docker push urtho/conduit-export-blksrv:latest
