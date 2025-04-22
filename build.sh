
#!/bin/bash
source version.sh

# docker stop lorikeet
# docker rm lorikeet

TAG="$DOCKER_REGISTRY_PROJECTS/lorikeet:$VERSION"

cd ./ui &&
  npm i &&
  npm run build &&
  cd - &&
  docker build --platform linux/amd64 -t "$TAG" . &&
  docker push "$TAG" 
#  echo "Running container..." &&
#  sh run_container.sh
