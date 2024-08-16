*** Settings ***
Documentation        Guest User API tests for go-mongo-reporter
Library        SeleniumLibrary
Library        RequestsLibrary
Library        JSONLibrary
Library        Collections

*** Variables ***
${HOST}    localhost
${PORT}    8080
${URL}     http://${HOST}:${PORT}
${TOKEN}

*** Settings ***
Test Setup    Setup Tests
Task Timeout    10


*** Test Cases ***

Test Successful Report Get
    [Documentation]    Test GET /api/reports with authenticated guest user. Should be successful.
    Create Reporter Session
    ${response}=        GET On Session        reporter-session        /api/reports
    Status Should Be    200
    ${len}=    Get Length    ${response.json()}
    Should Not Be Equal    ${len}    0
    Assert Response    ${response.json()}[0]
    Assert Response    ${response.json()}[1]
    Assert Response    ${response.json()}[2]

Test Unsuccessful Report Get
    [Documentation]    Test GET /api/reports with guest user authentication. Should be unsuccessful
    Create Unauthorized Reporter Session
    ${response}=        GET On Session        unauth-reporter-session        /api/reports       expected_status=anything 
    Status Should Be    401
    Should Not Be Empty    ${response.json()}[message]
    Should Be Equal As Strings    ${response.json()}[message]   401 Unauthorized

Test Successful GetByID
    [Documentation]    Test GET /api/reports/:id with guest user authentication. Should be successful. Since report IDs are generated first fetching all reports and use first ID
    Create Reporter Session
     ${r}=        GET On Session        reporter-session        /api/reports
    Status Should Be    200

    ${res}=    GET On Session    reporter-session    /api/reports/${r.json()}[0][ID]
    Status Should Be    200
    Assert Response    ${res.json()}

Test Unsuccessful Post Without Guest User Authentication
    [Documentation]    Test POST /api/reports without guest user authentication. Should be unsuccessful.
    Create Unauthorized Reporter Session
    ${report_data}=    Create Dictionary    topic=Robot report    description=Robot test report    author=Robot
    ${response}    POST On Session    reporter-session    /api/reports    expected_status=any    
    Status Should Be    401
    Should Be Equal As Strings    ${response.json()}[message]   401 Unauthorized

Test Unsuccessful Post With Guest User Authentication
    [Documentation]    Test POST /api/reports with guest user authentication. Should be unsuccessful.
    Create Reporter Session
    ${report_data}=    Create Dictionary    topic=Robot report    description=Robot test report    author=Robot
    ${response}    POST On Session    reporter-session    /api/reports    expected_status=any    
    Status Should Be    401
    Should Be Equal As Strings    ${response.json()}[message]   401 Unauthorized

Test Unsuccessful Update Without Guest User Authentication
    [Documentation]    Test PUT /api/reports/:id without guest user authentication. Should be unsuccessful.
    Create Unauthorized Reporter Session
    ${report_data}=    Create Dictionary    topic=Robot report    description=Robot test report    author=Robot
    ${response}    PUT On Session    reporter-session    /api/reports/1234    expected_status=any    
    Status Should Be    401
    Should Be Equal As Strings    ${response.json()}[message]   401 Unauthorized

Test Unsuccessful Update With Guest User Authentication
    [Documentation]    Test PUT /api/reports/:id with guest user authentication. Should be unsuccessful.
    Create Reporter Session
    ${report_data}=    Create Dictionary    topic=Robot report    description=Robot test report    author=Robot
    ${response}    PUT On Session    reporter-session    /api/reports/1234    expected_status=any    
    Status Should Be    401
    Should Be Equal As Strings    ${response.json()}[message]   401 Unauthorized

Test Unsuccessful Delete Without Guest User Authentication
    [Documentation]    Test DELETE /api/reports without guest user authentication. Should be unsuccessful.
    Create Unauthorized Reporter Session
    ${response}    DELETE On Session    reporter-session    /api/reports/1234    expected_status=any    
    Status Should Be    401
    Should Be Equal As Strings    ${response.json()}[message]   401 Unauthorized

Test Unsuccessful Delete With Guest User Authentication
    [Documentation]    Test DELETE /api/reports with guest user authentication. Should be unsuccessful.
    Create Reporter Session
    ${response}    DELETE On Session    reporter-session    /api/reports/1234    expected_status=any    
    Status Should Be    401
    Should Be Equal As Strings    ${response.json()}[message]   401 Unauthorized
    
*** Keywords ***
Create Unauthorized Reporter Session
    Create Session    unauth-reporter-session    ${URL}
Create Reporter Session
    ${headers}=    Create Dictionary    Authorization=${TOKEN}
    Create Session    reporter-session    ${URL}    headers=${headers}
Ensure That Guest User Exists
    [Documentation]    Create robot user as guest
    ${data}=    Create Dictionary    username=robot    email=robot@test.com    password=1234
    POST On Session    reporter-session    /signup    json=${data}    expected_status=anything
    
Login User and Get token
    [Documentation]    Login and get token
    ${data}=        Create Dictionary        username=robot        password=1234
    ${response}=        POST On Session        reporter-session        /login        json=${data}    expected_status=anything
    Status Should Be        200
    Dictionary Should Contain Key    ${response.json()}    message
    Dictionary Should Contain Key    ${response.json()}    token
    Should Be Equal As Strings    ${response.json()}[message]    Login succesful 
    Should Not Be Empty    ${response.json()}[token]
    Set Global Variable    ${TOKEN}    ${response.json()}[token]

Setup Tests
    [Documentation]    Create session, ensure that user exists, and finally set token
    Create Reporter Session
    Ensure That Guest User Exists
    Login User and Get token
    
Assert Response
    [Arguments]    ${RESPONSE_JSON}
    Should Not Be Empty    ${RESPONSE_JSON}[ID]
    Should Not Be Empty    ${RESPONSE_JSON}[topic]
    Should Not Be Empty    ${RESPONSE_JSON}[author]
    Should Not Be Empty    $${RESPONSE_JSON}[description]