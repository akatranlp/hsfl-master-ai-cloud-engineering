build_image() {
    docker build -t "akatranlp/$1:latest" -f $2 .
    docker push "akatranlp/$1:latest"
}

for d in ./src/*/ ; do
    folder=$(echo $d | sed -rE 's/\.\/src\/(.*)\//\1/g')
    dockerfile="${d}Dockerfile"

    build_image $folder $dockerfile &
done
