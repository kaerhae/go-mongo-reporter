# /tests/robot

This folder contains Robot Framework tests for migrate, "/", user handlers and report handlers.

To run these tests, `robotframework` must be installed with `robotframework-seleniumlibrary`, `robotframework-requests`, and `robotframework-jsonlibrary`. Also MongoDB has to be running for every test, and additionally reporter server must be running for `report_test.robot`, `user_test.robot`, and `smoke_test.robot`. Recommended way is to run docker-compose file on root of the project with: `docker-compose up -d`. 

<u>Note</u>: Check test suite variables, and configure them to match current environment.
