*** Settings ***
Documentation        Smoke test for go-mongo-reporter
Library        SeleniumLibrary
Library        RequestsLibrary
Library        JSONLibrary
Library        Collections

*** Variables ***
${HOST}        localhost
${PORT}        8080
${URL}        http://${HOST}:${PORT}
${MESSAGE}        Server up and running!

*** Test Cases ***
Do a Smoke test to /
    [Documentation]    Example test.
    Create Session        reporter-session        ${URL}        verify=${True}
    ${response}=        GET On Session        reporter-session        /
    Status Should Be    200
    Assert Response    ${response.text}


*** Keywords ***
Assert Response
    [Arguments]    ${response}
    Should Be Equal As Strings    ${MESSAGE}    ${response}