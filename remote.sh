#!/bin/bash

set -x

go build -gcflags "all=-N -l" .

ssh mick@macpro.local 'killall dlv'

rsync -HvaxP /Users/mick/Documents/GitHub/gb-launch/ /Volumes/mick/Documents/GitHub/gb-launch/
#rsync -HvaxP /Users/mick/Documents/GitHub/scribeHelpers/ /Volumes/mick/Documents/GitHub/scribeHelpers/
# rsync -HvaxP launch /Volumes/mick/.launch/bin/

# ssh mick@macpro.local 'PATH="/Users/mick/.launch/bin:$PATH"; cd /Users/mick/Documents/GitHub/gb-launch; /Users/mick/go/bin/dlv --listen=:2345 --headless=true --log=true --log-output=debugger,debuglineerr,gdbwire,lldbout,rpc: --api-version=2 --accept-multiclient exec ./launch'
# ssh mick@macpro.local 'PATH="/Users/mick/.launch/bin:$PATH"; cd /Users/mick/Documents/GitHub/gb-launch; /Users/mick/go/bin/dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./launch'
ssh mick@macpro.local "PATH=\"/Users/mick/.launch/bin:$PATH\"; cd /Users/mick/Documents/GitHub/gb-launch; /Users/mick/go/bin/dlv --listen=:2345 --headless=true --api-version=2 exec ./launch $@"

