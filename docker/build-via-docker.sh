#!/bin/bash -e

GID=$(id -g)

cd "$(dirname "$0")"
source _version.sh
echo "Building into $(pwd)"

# Empty -> build & package for all operating systems.
# Non-empty -> build & package only for this OS.
TARGET="$1"
if [ ! -z "$TARGET" ]; then
    echo "Only building for $TARGET."
fi

if [ -z "$GOPATH" ]; then
    echo "You have to define \$GOPATH." >&2
    exit 2
fi

# Use Docker to get Go in a way that allows overwriting the
# standard library with statically linked versions.
docker run -i --rm \
    -v $(pwd):/docker \
    -v "${GOPATH}:/go-local" \
    --env GOPATH=/go-local \
     golang:1.9 /bin/bash -e << EOT
echo -n "Using "
go version
cd \${GOPATH}/src/github.com/armadillica/svn-manager

function build {
    export GOOS=\$1
    export GOARCH=\$2
    export SUFFIX=\$3

    # GOARCH is always the same, so don't include in filename.
    TARGET=/docker/svn-manager-\$GOOS\$SUFFIX

    echo "Building \$TARGET"
    go get -a -ldflags '-s'
    go build -o \$TARGET

    if [ \$GOOS == linux ]; then
        strip \$TARGET
    fi
    chown $UID:$GID \$TARGET
}

export CGO_ENABLED=0
if [ -z "$TARGET" -o "$TARGET" = "linux"   ]; then build linux  amd64       ; fi
if [ -z "$TARGET" -o "$TARGET" = "windows" ]; then build windows amd64 .exe ; fi
if [ -z "$TARGET" -o "$TARGET" = "darwin"  ]; then build darwin  amd64      ; fi
EOT

# Package together with the static files
PREFIX="svn-manager-$APP_VERSION"
if [ -d $PREFIX ]; then
    rm -rf $PREFIX
fi
mkdir $PREFIX

echo "Assembling files into $PREFIX/"
rsync ../ui ../json_schemas $PREFIX -a --delete-after
cp ../{README.md,LICENSE.txt,CHANGELOG.md} $PREFIX/

if [ -z "$TARGET" -o "$TARGET" = "linux" ]; then
    echo "Creating archive for Linux"
    cp svn-manager-linux $PREFIX/svn-manager
    cp ../svn-manager.service $PREFIX/
    tar zcf $PREFIX-linux.tar.gz $PREFIX/
    rm -rf $PREFIX/svn-manager{,.service}
fi

if [ -z "$TARGET" -o "$TARGET" = "windows" ]; then
    echo "Creating archive for Windows"
    cp svn-manager-windows.exe $PREFIX/svn-manager.exe
    rm -f $PREFIX-windows.zip
    cd $PREFIX
    zip -9 -r -q ../$PREFIX-windows.zip *
    cd -
    rm -rf $PREFIX/svn-manager.exe
fi

if [ -z "$TARGET" -o "$TARGET" = "darwin" ]; then
    echo "Creating archive for Darwin"
    cp svn-manager-darwin $PREFIX/svn-manager
    rm -f $PREFIX-darwin.zip
    zip -9 -r -q $PREFIX-darwin.zip $PREFIX/
    rm -rf $PREFIX/svn-manager
fi

# Clean up after ourselves
rm -rf $PREFIX/

# Create the SHA256 sum file.
sha256sum svn-manager-$APP_VERSION-* | tee svn-manager-$APP_VERSION.sha256

echo "Done building & packaging SVN Manager $APP_VERSION."
