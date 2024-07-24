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

*** Test Cases ***
Test Signup
    [Documentation]    Test signup route handler
    Create Reporter Session
    ${data}=    Create Dictionary    username=robot    email=robot@test.com    password=1234    appRole=admin
    ${response}=    POST On Session    reporter-session    /signup    json=${data}    expected_status=anything
    IF    ${response.status_code} == 400
        Assert Signup With Existing Username    ${response.json()}
    ELSE
        Assert Signup With New User    ${response.json()}
    END

Test Signup with Missing Username
    [Documentation]    Test signup route handler. Should return error code and message.
    Create Reporter Session
    ${data}=    Create Dictionary    email=robot@test.com    password=1234    appRole=admin
    ${response}=    POST On Session    reporter-session    /signup    json=${data}    expected_status=anything
    Log    RESPONSE IS : ${response.json()}    console=yes
    Status Should Be    400
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]    Malformatted body   


Test Signup with Malformatted AppRole
    [Documentation]    Test signup route handler. Should return error code and message.
    Create Reporter Session
    ${data}=    Create Dictionary    username=robot-role-tester    password=blaah    email=robot@test.com   appRole=wrong
    ${response}=    POST On Session    reporter-session    /signup    json=${data}    expected_status=any
    Status Should Be    400
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]    Malformatted role   

Test Successful Login
    [Documentation]    Test login route handler. It should return successful code and token due to previous signup.
    Create Reporter Session
    ${data}=        Create Dictionary        username=robot        password=1234
    ${response}=        POST On Session        reporter-session        /login        json=${data}    expected_status=anything
    Status Should Be        200
    Dictionary Should Contain Key    ${response.json()}    message
    Dictionary Should Contain Key    ${response.json()}    token
    Should Be Equal As Strings    ${response.json()}[message]    Login succesful 
    Should Not Be Empty    ${response.json()}[token]
    

Test Unsuccessful Login
    [Documentation]    Test login route handler. It should return unsuccessful code and message with bad credentials.
    Create Reporter Session
    ${data}=    Create Dictionary    username=definitely    password=wrong
    ${response}=    POST On Session    reporter-session    /login    json=${data}    expected_status=anything
    Status Should Be    401
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    No user found    ${response.json()}[message]

Test Unsuccessful Login with Empty Password
    [Documentation]    Test login route handler. Test should return unsuccessful code and message, because empty password is not allowed
    Create Reporter Session
    ${data}=    Create Dictionary    username=    password=
    ${response}=    POST On Session    reporter-session    /login    json=${data}    expected_status=anything
    Status Should Be    400
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal    Malformatted body    ${response.json()}[message]

*** Keywords ***
Create Reporter Session
    Create Session    reporter-session    ${URL}    verify=${True}

Assert Signup With Existing Username
    [Arguments]    ${RESPONSE_JSON}
    Log    Asserting signup with existing username    console=yes
    Dictionary Should Contain Key    ${RESPONSE_JSON}    message
    Should Be Equal As Strings    ${RESPONSE_JSON}[message]    Username already exists 
    
Assert Signup With New User
    [Arguments]    ${RESPONSE_JSON}
    Log    Asserting signup with new user    console=yes
    Status Should Be    200
    Dictionary Should Contain Key    ${RESPONSE_JSON}    message
    Should Contain    ${RESPONSE_JSON}[message]    New user
    Should Contain    ${RESPONSE_JSON}[message]    was succesfully created