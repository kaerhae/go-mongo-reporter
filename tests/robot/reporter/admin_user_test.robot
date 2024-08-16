*** Settings ***
Documentation        User API tests for go-mongo-reporter
Library        SeleniumLibrary
Library        RequestsLibrary
Library        JSONLibrary
Library        Collections
Library        ../ExtendedJsonLib.py

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
Test Successful Get Users
    [Documentation]    Test should return successfully all users
    Create Authenticated Reporter Session    ${TOKEN}
    ${response}=    GET On Session    auth-reporter-session    /user-management/users
    Status Should Be    200
    Assert Get User    ${response.json()}[0]

Test Successful Get Single User
    [Documentation]    Test should return successfully all users
    Create Authenticated Reporter Session    ${TOKEN}
    ${response}=    GET On Session    auth-reporter-session    /user-management/users
    Status Should Be    200
    Assert Get User    ${response.json()}[0]

Test Successful Report Post
    [Documentation]    Test user POST route handler. Should be able to create user and inform success message and code.
    Create Authenticated Reporter Session    ${TOKEN}
    Find and Delete User
    ${report_data}=    Create Dictionary    username=testuser    email=test@user.com    password=1234
    ${response}=    POST On Session    auth-reporter-session    /user-management/users    json=${report_data}    expected_status=any
    Status Should Be    201
    Dictionary Should Contain Key    ${response.json()}    message
    Should Contain    ${response.json()}[message]       was succesfully created
    

Test Unsuccessful Report Post
    [Documentation]    Test user POST route handler. Should be unsuccessful, due to malformatted body.
    Create Authenticated Reporter Session    ${TOKEN}
    ${report_data}=    Create Dictionary    username=testuser    password=1234
    ${response}=    POST On Session    auth-reporter-session    /user-management/users    json=${report_data}    expected_status=any
    Status Should Be    400
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]    Malformatted body

Test Unsuccessful Report Post Without Authentication
    [Documentation]    Test user POST route handler. Should be unsuccessful, due to malformatted body.
    Create Reporter Session
    ${report_data}=    Create Dictionary    username=testuser    password=1234
    ${response}=    POST On Session    reporter-session    /user-management/users    json=${report_data}    expected_status=any
    Status Should Be    401
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]       401 Unauthorized


Test Successful Report Update
    [Documentation]    Test user PUT route handler. Should be able to update user and inform success message and code.
    Create Authenticated Reporter Session    ${TOKEN}
    Ensure Created User
    ${report_data}=    Create Dictionary    username=testuser    email=new@email.com    password=1234
    ${get_response}=     GET On Session    auth-reporter-session    /user-management/users
    Status Should Be    200
    ${result}=    Filter By Value    ${get_response.json()}    username    testuser
    ${response}=    PUT On Session    auth-reporter-session    /user-management/users/${result}[0][id]    json=${report_data}    expected_status=any
    Status Should Be    200
    Dictionary Should Contain Key    ${response.json()}    message
    Should Contain    ${response.json()}[message]       was succesfully updated
    Find and Delete User


Test Unsuccessful Report Update Without Authentication
    [Documentation]    Test user PUT route handler. Should be unsuccessful.
    Create Reporter Session
    ${report_data}=    Create Dictionary    username=testuser    email=new@email.com    password=1234
    Ensure Created User
    ${report_data}=    Create Dictionary    username=testuser    email=new@email.com    password=1234
    ${get_response}=     GET On Session    auth-reporter-session    /user-management/users
    Status Should Be    200
    ${result}=    Filter By Value    ${get_response.json()}    username    testuser
    ${response}=    PUT On Session    reporter-session    /user-management/users/${result}[0][id]    json=${report_data}    expected_status=anything
    Status Should Be    401
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]       401 Unauthorized

    

Test Successful User UpdatePassword
    [Documentation]    Test user PUT route handler. Should be able to update user and inform success message and code.
    Create Authenticated Reporter Session    ${TOKEN}
    Ensure Created User
    ${get_response}=     GET On Session    auth-reporter-session    /user-management/users
    Status Should Be    200
    ${result}=    Filter By Value    ${get_response.json()}    username    testuser
    ${report_data}=    Create Dictionary    userID=${result}[0][id]    oldPassword=1234    newPassword=12345

    ${response}=    PUT On Session    auth-reporter-session    /change-password    json=${report_data}
    Status Should Be    200
    Dictionary Should Contain Key    ${response.json()}    message
    Should Contain    ${response.json()}[message]       Password updated successfully

    ${data}=        Create Dictionary        username=testuser        password=12345
    ${response}=        POST On Session        reporter-session        /login        json=${data}
    Status Should Be        200

    Find and Delete User


Test Unsuccessful User UpdatePassword Without Authentication
    [Documentation]    Test user PUT route handler. Should be unsuccessful.
    Create Reporter Session
    ${report_data}=    Create Dictionary    username=testuser    email=new@email.com    password=1234
    Ensure Created User
    ${get_response}=     GET On Session    auth-reporter-session    /user-management/users
    Status Should Be    200
    ${result}=    Filter By Value    ${get_response.json()}    username    testuser
    ${report_data}=    Create Dictionary    userID=${result}[0][id]    oldPassword=1234    newPassword=12345
    ${response}=    PUT On Session    reporter-session    /change-password    json=${report_data}    expected_status=any
    Status Should Be    401
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]       401 Unauthorized


Test Successful User UpdatePermissions
    [Documentation]    Test user PUT route handler. Should be able to update user and inform success message and code.
    Create Authenticated Reporter Session    ${TOKEN}
    Ensure Created User
    ${get_response}=     GET On Session    auth-reporter-session    /user-management/users
    Status Should Be    200
    ${result}=    Filter By Value    ${get_response.json()}    username    testuser
    ${permissions_data}=    Create Dictionary    admin=${True}    write=${True}    read=${False}
    ${report_data}=    Create Dictionary    userID=${result}[0][id]    permissions=${permissions_data}

    ${response}=    PUT On Session    auth-reporter-session    /user-management/users/change-permissions    json=${report_data}
    Status Should Be    200
    Dictionary Should Contain Key    ${response.json()}    message
    Should Contain    ${response.json()}[message]       permissions updated

    ${response}=    GET On Session    auth-reporter-session    /user-management/users/${result}[0][id]

    Should Be Equal    ${response.json()}[permission][admin]    ${True}
    Should Be Equal    ${response.json()}[permission][write]    ${True}
    Should Be Equal    ${response.json()}[permission][read]    ${False}
    Find and Delete User


Test Unsuccessful User UpdatePermissions Without Authentication
    [Documentation]    Test user PUT route handler. Should be unsuccessful.
    Create Reporter Session
    Ensure Created User
    ${get_response}=     GET On Session    auth-reporter-session    /user-management/users
    Status Should Be    200
    ${result}=    Filter By Value    ${get_response.json()}    username    testuser
    ${report_data}=    Create Dictionary    userID=${result}[0][id]    oldPassword=1234    newPassword=12345
    ${response}=    PUT On Session    reporter-session    /user-management/users/change-permissions    json=${report_data}    expected_status=anything
    Status Should Be    401
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]       401 Unauthorized

Test Successful Report Delete
    [Documentation]    Test user DELETE route handler. Should be able to delete user and inform success message and code.
    Create Authenticated Reporter Session    ${TOKEN}
    Ensure Created User
    ${get_response}=     GET On Session    auth-reporter-session    /user-management/users
    Status Should Be    200
    ${result}=    Filter By Value    ${get_response.json()}    username    testuser
    ${response}=    DELETE On Session    auth-reporter-session    /user-management/users/${result}[0][id]    expected_status=any
    Status Should Be    200
    Dictionary Should Contain Key    ${response.json()}    message
    Should Contain    ${response.json()}[message]       Deleted 1 user


Test Unsuccessful Report Delete Without Authentication
    [Documentation]    Test user DELETE route handler. Should be unsuccessful.
    Create Reporter Session
    Ensure Created User
    ${get_response}=     GET On Session    auth-reporter-session    /user-management/users
    Status Should Be    200
    ${result}=    Filter By Value    ${get_response.json()}    username    testuser
    ${response}=    DELETE On Session    reporter-session    /user-management/users/${result}[0][id]    expected_status=any
    Status Should Be    401
    Dictionary Should Contain Key    ${response.json()}    message
    Should Be Equal As Strings    ${response.json()}[message]       401 Unauthorized

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

Ensure Created User
    [Documentation]    Test user POST route handler. Should be able to create user and inform success message and code.
    Create Authenticated Reporter Session    ${TOKEN}
    ${report_data}=    Create Dictionary    username=testuser    email=test@user.com    password=1234
    ${response}=    POST On Session    auth-reporter-session    /user-management/users    json=${report_data}    expected_status=any

Find and Delete User 
    Create Authenticated Reporter Session    ${TOKEN}
    ${response}=    GET On Session    auth-reporter-session    /user-management/users
    Status Should Be    200
    ${result}=    Filter By Value    ${response.json()}    username    testuser
    ${length}=    Check List Length    ${result}
    IF  $length > 0
        ${res}=    DELETE On Session    auth-reporter-session    /user-management/users/${result}[0][id]   
        Status Should Be    200
    END

Assert Get User
    [Arguments]    ${RESPONSE_JSON}
    Dictionary Should Contain Key   ${RESPONSE_JSON}    id
    Dictionary Should Contain Key    ${RESPONSE_JSON}    username
    Dictionary Should Contain Key    ${RESPONSE_JSON}    email
    Dictionary Should Contain Key    ${RESPONSE_JSON}    createdAt
    Dictionary Should Contain Key    ${RESPONSE_JSON}    reports
    Dictionary Should Contain Key    ${RESPONSE_JSON}    permission
    Dictionary Should Not Contain Key    ${RESPONSE_JSON}    passwordHash

