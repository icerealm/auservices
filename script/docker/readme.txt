#access postgres container shell
docker exec -ti pg_dev bash
#access postgres sh container database
docker exec -ti pg_dev psql -U postgres