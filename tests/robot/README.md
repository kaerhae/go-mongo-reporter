# /tests/robot

This folder contains Robot Framework tests for migrate, "/", user handlers and report handlers.

To run these tests, `robotframework` must be installed with `robotframework-seleniumlibrary`, `robotframework-requests`, and `robotframework-jsonlibrary`. Also MongoDB has to be running for every test, and additionally reporter server must be running for `report_test.robot`, `user_test.robot`, and `smoke_test.robot`.

Install dependencies with:
```bash
pip install -r requirements.txt
```

Recommended way is to run dedicated docker-compose file on root of the project. To run docker-compose in test environment, run:

```bash
docker-compose -f docker-compose.robot-test.yml up -d
```

If docker-compose has launched multiple times, environment can be cleaned to make sure that testing environment is neutral. To clean environment, run:
```bash
docker-compose -f docker-compose.robot-test.yml down -v --rmi all --remove-orphans
```

<u>Note</u>: Check test suite variables, and configure them to match current environment.

