*** Settings ***
Documentation        User API tests for go-mongo-reporter
Library        SeleniumLibrary
Library        RequestsLibrary
Library        JSONLibrary
Library        Collections

*** Variables ***
${HOST}    localhost
${PORT}    8080
${URL}     http://${HOST}:${PORT}
${TOKEN}
${USER_ID}
${REPORT_ID}
${UPDATED_TOPIC}    Updated Robot report    
${UPDATED_DESCRIPTION}    Updated Robot test report    
${UPDATED_AUTHOR}    Updated Robot

*** Settings ***
Suite Setup   Login User and Get token

*** Test Cases ***
Test Successful Report Get
    [Documentation]    Test GET /api/reports with authenticated guest user. Should be successful.
    Create Authenticated Reporter Session    ${TOKEN}
    ${response}=        GET On Session        auth-reporter-session        /api/reports
    Status Should Be    200
    ${len}=    Get Length    ${response.json()}
    Should Not Be Equal    ${len}    0
    Assert Response    ${response.json()}[0]
    Assert Response    ${response.json()}[1]
    Assert Response    ${response.json()}[2]

Test Successful Report GetByID
    [Documentation]    Test report GET route handler. Should be able to retrieve existing reports as json, with suitable status code
    Create Authenticated Reporter Session    ${TOKEN}
     ${r}=        GET On Session        auth-reporter-session        /api/reports
    Status Should Be    200

    ${res}=    GET On Session    auth-reporter-session    /api/reports/${r.json()}[0][ID]
    Status Should Be    200
    Assert Response    ${res.json()}


Test Unsuccessful Report GetByID
    [Documentation]    Test report GET route handler. Should be able to retrieve existing report as json, with suitable status code
    Create Authenticated Reporter Session    ${TOKEN}
    ${response}=        GET On Session        auth-reporter-session        /api/reports/111122223333444455556666    expected_status=any
    Status Should Be    400
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]    Error: mongo: no documents in result

Test Successful Report Post
    [Documentation]    Test report POST route handler. Should be able to create report and inform success message and code.
    Create Authenticated Reporter Session    ${TOKEN}
    ${report_data}=    Create Dictionary    topic=Robot report    description=Robot test report    author=Robot    userID=${USER_ID}
    ${response}=    POST On Session    auth-reporter-session    /api/reports    json=${report_data}
    Status Should Be    200
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]    Report was succesfully created

Test Report Post Without UserID
    [Documentation]    Test report POST route handler. Should return error message and status code.
    Create Authenticated Reporter Session    ${TOKEN}
    ${report_data}=    Create Dictionary    topic=Robot report    description=Robot test report    author=Robot
    ${response}=    POST On Session    auth-reporter-session    /api/reports    json=${report_data}    expected_status=any
    Status Should Be    400
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]    No userID found on request


Test Report Put
    [Documentation]    Test report PUT route handler. Should be able to successfully update report
    Create Authenticated Reporter Session    ${TOKEN}
    ${report_data}=    Create Dictionary    topic=${UPDATED_TOPIC}    description=${UPDATED_DESCRIPTION}    author=${UPDATED_AUTHOR}
    ${response}=    PUT On Session    auth-reporter-session    /api/reports/${REPORT_ID}    json=${report_data}
    Status Should Be    200
    Dictionary Should Contain Key    ${response.json()}    message
    Should Contain    ${response.json()}[message]    was succesfully updated
    
    ${response}=        GET On Session        auth-reporter-session        /api/reports/${REPORT_ID}
    Status Should Be    200
    Assert Updated Report    ${response.json()}

Test Report Put with Non-existing report
    [Documentation]    Test report PUT route handler with non-existing report. Should return error message and status code
    Create Authenticated Reporter Session    ${TOKEN}
    ${report_data}=    Create Dictionary    topic=${UPDATED_TOPIC}    description=${UPDATED_DESCRIPTION}    author=${UPDATED_AUTHOR}
    ${response}=    PUT On Session    auth-reporter-session    /api/reports/111122223333444455556666    json=${report_data}    expected_status=any
    Status Should Be    404
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]    No report found
    


Test Report Delete
    [Documentation]    Test report DELETE route handler. Should be able to successfully delete report.
    Create Authenticated Reporter Session    ${TOKEN}
    ${response}=        DELETE On Session        auth-reporter-session        /api/reports/${REPORT_ID}
    Status Should Be    200
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]    Deleted 1 reports

    ${response}=        GET On Session        auth-reporter-session        /api/reports/${REPORT_ID}    expected_status=anything
    Status Should Be    400
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]    Error: mongo: no documents in result


Test Report Delete with Non-existing report
    [Documentation]    Test report DELETE route handler with non-existing report. Should return error message and status code.
    Create Authenticated Reporter Session    ${TOKEN}
    ${response}=    DELETE On Session    auth-reporter-session    /api/reports/123    expected_status=anything
    Status Should Be    404
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]    Error while validating request: no report found


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
    ${data}=        Create Dictionary        username=robot_admin_user        password=1234
    ${response}=        POST On Session        reporter-session        /login        json=${data}    expected_status=anything
    Status Should Be        200
    Dictionary Should Contain Key    ${response.json()}    message
    Dictionary Should Contain Key    ${response.json()}    token
    Should Be Equal As Strings    ${response.json()}[message]    Login succesful 
    Should Not Be Empty    ${response.json()}[token]
    Set Global Variable    ${TOKEN}    ${response.json()}[token]
    Set Global Variable    ${USER_ID}    ${response.json()}[userID]


Assert Response
    [Arguments]    ${RESPONSE_JSON}
    Should Not Be Empty    ${RESPONSE_JSON}[ID]
    Should Not Be Empty    ${RESPONSE_JSON}[UserID]

    Set Global Variable    ${REPORT_ID}    ${RESPONSE_JSON}[ID]

Assert Updated Report
    [Arguments]    ${RESPONSE_JSON}
    Should Be Equal As Strings    ${RESPONSE_JSON}[topic]    ${UPDATED_TOPIC}
    Should Be Equal As Strings    ${RESPONSE_JSON}[description]    ${UPDATED_DESCRIPTION}
    Should Be Equal As Strings    ${RESPONSE_JSON}[author]    ${UPDATED_AUTHOR}