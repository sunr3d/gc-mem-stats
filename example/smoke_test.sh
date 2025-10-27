# 1000 запросов на сервер
for i in {1..1000}; do
  curl http://localhost:8080/metrics > /dev/null
done

xdg-open http://localhost:8080/metrics