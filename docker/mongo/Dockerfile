FROM mongo

COPY ./docker/mongo/robot-test-report-data.json .
CMD mongoimport --host mongo --username root --password example --authenticationDatabase=admin --db reporter --collection reports --type json --file /robot-test-report-data.json --jsonArray
