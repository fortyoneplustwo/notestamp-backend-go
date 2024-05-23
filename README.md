## Notestamp back-end (Go)
A back-end server for [notestamp.com](notestamp.com) that allows users to register an account and save their projects online.

## Stack
Initially implemented using NodeJS and MongoDB, this version is built in Go and uses PostgreSQL as the database. The same S3 bucket is used to store notes and media files.

## Authorization
JWTs are used to maintain user sessions. This version implements a cleanup of the revoked tokens database that runs concurrently whenever the a user signs out of their account.

## Interfaces
Unlike the NodeJS version, this one takes advantage of Go interfaces to abstract stores such as databases and file storage so that the underlying implementations (currently PostgreSQL and S3) can be easily switched if needed.

## Testing
Integration tests written in bash using `curl`.

## API
`POST /auth/register`  
**Description**: Create a user account.
**Accepts**: Form with fields: `username` and `password`.

`POST /auth/login`   
**Description**: Log into a registered account.  
**Accepts**: Form with fields: `username` and `password`.

`POST /auth/logout`  
**Description**: Log out of user account.

`POST /project/save`  
**Description**: Save a project.  
**Accepts**: Form data with fields: `metadata` (json format), `notesFile`and `mediaFile`.  
**Produces**: JSON containing the updated list of projects after save.

`GET /project/list`  
**Description**: Get a list of the user's saved projects.
**Produces**: JSON containing a list of projects.

`GET /project/get/{title}`  
**Description**: Retrieve saved project with provided title.
**Produces**: JSON containing metadata and notes content.

`DELETE /project/delete/{title}`  
**Description**: Delete a project with provided title.
**Produces**: JSON containing the updated list of projects after deletion.

`GET /media/download/{title}`  
**Description**: Download the media related to the project with provided title.
**Produces**: A buffer containing the media.

`GET /media/stream/{title}`
**Description**: Stream the media related to the project with provided title.
**Produces**: JSON containing a url which can be used to stream the media directly from the media store.

