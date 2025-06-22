clear
docker build -t minecraft-web .
docker run --rm -p 8000:8000 minecraft-web