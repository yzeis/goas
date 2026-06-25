#!/usr/bin/env sh

export CUR="github.com/yzidev/goas" # example: github.com/user/old-lame-name
export NEW="github.com/yzeis/goas" # example: github.com/user/new-super-cool-name
go mod edit -module ${NEW}
find . -type f -name '*.go' -exec perl -pi -e 's/$ENV{CUR}/$ENV{NEW}/g' {} \;