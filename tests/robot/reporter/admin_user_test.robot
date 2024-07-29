*** Settings ***
Documentation        User API tests for go-mongo-reporter
Library        SeleniumLibrary
Library        RequestsLibrary
Library        JSONLibrary
Library        Collections

*** Variables ***
${HOST}    localhost
${PORT}    8089
${URL}     http://${HOST}:${PORT}
${TOKEN}

*** Settings ***
Suite Setup   Login User and Get token

*** Test Cases ***


Test Successful Get Users
    [Documentation]    Test should return successfully all users
    Create Authenticated Reporter Session    ${TOKEN}
    ${response}=    GET On Session    auth-reporter-session    /user-management/users
    Log    ${response.json()}    console=yes
    Status Should Be    200
    Assert Get Report    ${response.json()}[0]


*** Keywords ***
Create Reporter Session
    Create Session    reporter-session    ${URL}    verify=${True}

Create Authenticated Reporter Session
    [Arguments]    ${TOKEN}
    ${headers}=    Create Dictionary    Authorization=${TOKEN}
    Create Session    auth-reporter-session    ${URL}    verify=${True}    headers=${headers}


Login User and Get token
    [Documentation]    Test login route handler. It should return successful code and token due to previous signup.
    Create Reporter Session
    ${data}=        Create Dictionary        username=root        password=example
    ${response}=        POST On Session        reporter-session        /login        json=${data}    expected_status=anything
    Status Should Be        200
    Dictionary Should Contain Key    ${response.json()}    message
    Dictionary Should Contain Key    ${response.json()}    token
    Should Be Equal As Strings    ${response.json()}[message]    Login succesful 
    Should Not Be Empty    ${response.json()}[token]
    Set Global Variable    ${TOKEN}    ${response.json()}[token]
    Log    TOKEN IS ${TOKEN}    console=yes
Assert Get Report
    [Arguments]    ${RESPONSE_JSON}
    Dictionary Should Contain Key   ${RESPONSE_JSON}    id
    Dictionary Should Contain Key    ${RESPONSE_JSON}    username
    Dictionary Should Contain Key    ${RESPONSE_JSON}    email
    Dictionary Should Contain Key    ${RESPONSE_JSON}    createdAt
    Dictionary Should Contain Key    ${RESPONSE_JSON}    reports
    Dictionary Should Contain Key    ${RESPONSE_JSON}    appRole
    Dictionary Should Not Contain Key    ${RESPONSE_JSON}    passwordHash