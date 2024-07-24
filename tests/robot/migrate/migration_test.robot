*** Settings ***
Documentation        Integration tests for cmd/migrate/*
Library    OperatingSystem
Library    Builtin
Library    SeleniumLibrary
Library    RequestsLibrary
Library    JSONLibrary
Library    Collections
Library    String

*** Variables ***
${MONGO_USER}    root
${MONGO_PASS}    example
${MONGO_IP}    172.18.0.2
${MONGO_PORT}    27017
${REPORTER_ROOT_USER}    root
${REPORTER_ROOT_PASSWORD}    example


${ROBOT_FOLDER}    ${CURDIR}
${MIGRATE_CMD_FOLDER}    ${ROBOT_FOLDER}${/}..${/}..${/}..${/}cmd/migrate
${BIN_FOLDER}    ${ROBOT_FOLDER}${/}..${/}..${/}..${/}bin${/}migrate
${MIGRATE_BINARY}    ${BIN_FOLDER}/migrate


*** Settings ***
Test Setup    Build Binary
Task Timeout    10
Test Teardown    Remove Possible Admin User and Remove Binary    

*** Test Cases ***
Test without Parameters
    [Documentation]    Test should return exit and missing argument error output
    ${output}=    Run    ${MIGRATE_BINARY}
    Should Contain    ${output}    Missing argument. Set argument to up or down


Test with Invalid Parameter
    [Documentation]    Test should return exit and invalid parameter error output
    ${output}=    Run    ${MIGRATE_BINARY} blaah
    Should Contain    ${output}    Invalid option for --create-user-admin

Test with Case-insensitive Parameter
    [Documentation]    Test should understand correct parameters case-insensitive way
    ${output}=    Run    ${MIGRATE_BINARY} UP
    ${output1}=    Run    ${MIGRATE_BINARY} DOWN
    ${output2}=    Run    ${MIGRATE_BINARY} Up
    ${output3}=    Run    ${MIGRATE_BINARY} Down
    ${output4}=    Run    ${MIGRATE_BINARY} uP
    ${output5}=    Run    ${MIGRATE_BINARY} DoWn

    Should Contain    ${output}    Successfully created admin user
    Should Contain    ${output2}    Successfully created admin user
    Should Contain    ${output4}    Successfully created admin user

    Should Contain    ${output1}    Successfully deleted admin user
    Should Contain    ${output3}    Successfully deleted admin user
    Should Contain    ${output5}    Successfully deleted admin user

Test Migrate Up When Already Up
    [Documentation]    Test that migrate correctly fails if migration has already done
    ${output}=    Run    ${MIGRATE_BINARY} up
    ${output1}=    Run    ${MIGRATE_BINARY} up

    Should Contain    ${output}    Successfully created admin user
    Should Contain    ${output1}    migration already done and user exists

Test Migrate Down When Already Down
    [Documentation]    Test migrate correctly fails if migration is already down
    ${output}=    Run    ${MIGRATE_BINARY} down
    ${output1}=    Run    ${MIGRATE_BINARY} up
    ${output2}=    Run    ${MIGRATE_BINARY} down
    ${output3}=    Run    ${MIGRATE_BINARY} down

    Should Contain    ${output}    no user found
    Should Contain    ${output1}    Successfully created admin user
    Should Contain    ${output2}    Successfully deleted admin user
    Should Contain    ${output3}    no user found



*** Keywords ***
Build Binary
    [Documentation]    Build migrate
    Set Environment Variable    MONGO_USER    ${MONGO_USER}
    Set Environment Variable    MONGO_PASS    ${MONGO_PASS}
    Set Environment Variable    MONGO_IP    ${MONGO_IP}
    Set Environment Variable    MONGO_PORT    ${MONGO_PORT}
    Set Environment Variable    REPORTER_ROOT_USER    ${REPORTER_ROOT_USER}
    Set Environment Variable    REPORTER_ROOT_PASSWORD    ${REPORTER_ROOT_PASSWORD}
    Directory Should Exist    ${MIGRATE_CMD_FOLDER}
    Directory Should Exist    ${BIN_FOLDER}

    Run    go build -o ${MIGRATE_BINARY} ${MIGRATE_CMD_FOLDER}/main.go
    File Should Exist    ${MIGRATE_BINARY}

Remove Possible Admin User and Remove Binary
    [Documentation]    Teardown possible added user and remove Binary
    Run    ${MIGRATE_BINARY} down
    Run    rm ${MIGRATE_BINARY}
    File Should Not Exist    ${MIGRATE_BINARY}

    
    