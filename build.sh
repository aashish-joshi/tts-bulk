#!/usr/bin/env bash
# A simple bash script to build a Go package for multiple platforms
# Usage: ./build.sh <go-package-name>

# Read the package name from the supplied argument
package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi
package_split=(${package//\// })
package_name=${package_split[-1]}
# Remove the .go extension from the final package name
final_name=(${package_name//\.go/ })
build_dir='builds'
platforms=("linux/amd64" "windows/amd64" "darwin/amd64")

# Check if the build_dir is empty. If not, ask user if they want to clean it.
if [ -d $build_dir ]; then
    read -p "The build directory is not empty. Do you want to clean it? (yes/no): " clean
    if [[ $clean == "yes" ]]; then
        rm -rf $build_dir/*
    fi
else
    mkdir $build_dir
fi

# Generate builds for each platform
for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$build_dir'/'$GOOS'-'$GOARCH'/'$final_name
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi    

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
    if [ $? -ne 0 ]; then
           echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
