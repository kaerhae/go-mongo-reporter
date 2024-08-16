# /api

This folder contains gin router and handlers for routes. App contains several routes for report and user-management. 

* GET /api/reports: Retrieves all reports from database. Requires admin rights or minimum read permissions.
* GET /api/reports/:{id} : Retrieves single report by ID from database. Requires admin rights or minimum read permissions.
* POST /api/reports: Insert single report to database. Requires Report model from body and requires admin rights or minimum write permissions
* PUT /api/reports/:{id} : Updates single report to database. Requires Report model from body and requires admin rights or minimum write permissions
* DELETE /api/reports/:{id} : Deletes single report by ID from database. Requires admin rights or minimum write permissions.


* GET /user-management/users: Retrieves all users from database. Requires admin rights.
* GET /user-management/users/:{id} : Retrieves single user by ID from database. Requires admin rights.
* POST /user-management/users: Insert single user to database. Requires CreateUser model from body and requires admin rights.
* PUT /user-management/users/:{id} : Updates single user to database. Requires User model from body and requires admin rights.
* PUT /user-management/users/change-permissions: Updates single user permissions to database. Requires UserPermissionUpdate model from body and requires admin rights.
* DELETE /user-management/users/:{id} : Deletes single user by ID from database. Requires admin rights.


* PUT /change-password: Updates user password to database. Requires UserPasswordChange model from body and either admin rights or correct user-to-be-updated session with read permissions
* POST /login: Performs a login to app. Requires LoginUser model from body.
* POST /signup: Performs a signup to app. Requires CreateGuestUser model from body.