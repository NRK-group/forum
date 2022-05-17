# Forum

## Objectives

This project consists in creating a web forum that allows :

- communication between users.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.

## [SQLite](https://www.sqlite.org/index.html)

In order to store the data in your forum (like users, posts, comments, etc.) you will use the database library SQLite.

SQLite is a popular choice as an embedded database software for local/client storage in application software such as web browsers. It enables you to create a database as well as controlling it by using queries.

To structure your database and to achieve better performance we highly advise you to take a look at the entity relationship diagram and build one based on your own database.

- You must use at least one SELECT, one CREATE and one INSERT queries.
To know more about SQLite you can check the SQLite page.

## Authentication

In this segment the client must be able to register as a new user on the forum, by inputting their credentials. You also have to create a login session to access the forum and be able to add posts and comments.

You should use cookies to allow each user to have only one opened session. Each of this sessions must contain an expiration date. It is up to you to decide how long the cookie stays "alive". The use of UUID is a Bonus task.

Instructions for user registration:

- Must ask for email
  - When the email is already taken return an error response.
- Must ask for username
- Must ask for password
  - The password must be encrypted when stored (this is a Bonus task)

The forum must be able to check if the email provided is present in the database and if all credentials are correct. It will check if the password is the same with the one provided and, if the password is not the same, it will return an error response.

## Communication

In order for users to communicate between each other, they will have to be able to create posts and comments.

- Only registered users will be able to create posts and comments.
- When registered users are creating a post they can associate one or more categories to it.
  - The implementation and choice of the categories is up to you.
- The posts and comments should be visible to all users (registered or not).
- Non-registered users will only be able to see posts and comments.

## Likes and Dislikes

- Only registered users will be able to like or dislike posts and comments.

- The number of likes and dislikes should be visible by all users (registered or not).

## Filter

You need to implement a filter mechanism, that will allow users to filter the displayed posts by :

- categories
- created posts
- liked posts

You can look at filtering by categories as subforums. A subforum is a section of an online forum dedicated to a specific topic.

Note that the last two are only available for registered users and must refer to the logged in user.

## [Docker](https://docs.docker.com)

- For this project you must create at least :
  - one Dockerfile
  - one image
  - one container
- You must apply metadata to Docker objects.

-You have to take caution of unused object (often referred to as "garbage collection").

## Instructions

- You must use SQLite.
- You must handle website errors, HTTP status.
- You must handle all sort of technical errors.
- The code must respect the [good practices](https://git.learn.01founders.co/root/public/src/branch/master/subjects/good-practices/README.md).
- It is recommended to have test files for [unit testing](https://go.dev/doc/tutorial/add-a-test).

## Allowed packages

- All [standard Go](https://pkg.go.dev/std) packages are allowed.
- [sqlite3](https://github.com/mattn/go-sqlite3)
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [UUID](https://github.com/satori/go.uuid)

## Optionals

### Objectives of Image upload
>
> You must follow the same principles as the first subject.
>
> In forum image upload, registered users have the possibility to create a post containing an image as well as text.
>
> - When viewing the post, users and guests should see the image associated to it.
>
> There are several extensions for images like: JPEG, SVG, PNG, GIF, etc. In this project you have to handle at least JPEG, PNG and GIF types.
>
> The max size of the images to load should be 20 mb. If there is an attempt to load an image greater than 20mb, an error message should inform the user that the image is too big.
>
> ### Hints
>
> - Be cautious with the size of the images.
>
### Objectives of Authentication
>
> The goal of this project is to implement, into your forum, new ways of authentication. You have to be able to register and to login using at least Google and Github authentication tools.
>
> Some examples of authentication means are:
>
> - Facebook
> - GitHub
> - Google

### Objectives of Security

>
> For this project you must take into account the security of your forum.
>
> - You should implement a Hypertext Transfer Protocol Secure (HTTPS) protocol :
>
>   - Encrypted connection : for this you will have to generate an SSL certificate, you can think of this like a identity card for your website. You can create your certificates or use "Certificate Authorities"(CA's)
>
>   - We recommend you to take a look into cipher suites.
>
> - The implementation of Rate Limiting must be present on this project
>
> - You should encrypt at least the clients passwords. As a Bonus you can also encrypt the database, for this you will have to create a password for your database.
>
> Sessions and cookies were implemented in the previous project but not under-pressure (tested in an attack environment). So this time you must take this into account.
>
> - Clients session cookies should be unique. For instance, the session state is stored on the server and the session should present an unique identifier. This way the client has no direct access to it. Therefore, there is no way for attackers to read or tamper with session state.
>
> ### Hints
>
> - You can take a look at the openssl manual.
> - For the session cookies you can take a look at the Universal Unique Identifier (UUID)

### Objectives of Moderation

> The forum moderation will be based on a moderation system. It must present a moderator that, depending on the access level of a user or the forum set-up, approves posted messages before they become publicly visible.
>
> - The filtering can be done depending on the categories of the post being sorted by irrelevant, obscene, illegal or insulting.
> For this optional you should take into account all types of users that can exist in a forum and their levels.
>
> You should implement at least 4 types of users :
>
> ### Guests
>
> - These are unregistered-users that can neither post, comment, like or dislike a post. They only have the permission to see those posts, comments, likes or dislikes.
>
> ### Users
>
> - These are the users that will be able to create, comment, like or dislike posts.
>
> ### Moderators
>
> - Moderators, as explained above, are users that have a granted access to special functions :
>   -They should be able to monitor the content in the forum by deleting or reporting post to the admin
> -To create a moderator the user should request an admin for that role
>
> ### Administrators
>
> - Users that manage the technical details required for running the forum. This user must be able to :
>   - Promote or demote a normal user to, or from a moderator user.
>   - Receive reports from moderators. If the admin receives a report from a moderator, he can respond to that report
>   - Delete posts and comments
>   - Manage the categories, by being able to create and delete them.

### Objectives of Advanced Features

> - In forum advanced features, you will have to implement the following features :
>
> - You will have to create a way to notify users when their posts are :
>
>   - liked/disliked
>   - commented
> - You have to create an activity page that tracks the user own activity. In other words, a page that :
>
>   - Shows the user created posts
>   - Shows where the user left a like or a dislike
>   - Shows where and what the user has been commenting. For this, the comment will have to be shown, as well as the post commented
>   -You have to create a section where you will be able to Edit/Remove posts and comments.

## Audit

> [Forum](https://git.learn.01founders.co/root/public/src/branch/master/subjects/forum/audit)
>
> [Image Upload](https://git.learn.01founders.co/root/public/src/branch/master/subjects/forum/image-upload/audit.md)
>
> [Authentication](https://git.learn.01founders.co/root/public/src/branch/master/subjects/forum/authentication/audit.md)
>
> [Security](https://git.learn.01founders.co/root/public/src/branch/master/subjects/forum/audit)
>
> [Moderation](https://git.learn.01founders.co/root/public/src/branch/master/subjects/forum/moderation/audit.md)
>
> [Advanced Features](https://git.learn.01founders.co/root/public/src/branch/master/subjects/forum/advanced-features)

This project will help you learn about:

#### Forum

- The basics of web :
  - HTML
  - HTTP
  - Sessions and cookies
- Using and setting up Docker
  - Containerizing an application
  - Compatibility/Dependency
  - Creating images
- SQL language
  - Manipulation of databases
- The basics of encryption

#### Image Upload

- Image manipulation
- Image types

#### Authentication

- Sessions and cookies.
- Protecting routes.

#### Security

- HTTPS
- Cipher suites
- Goroutines
- Channels
- Rate Limiting
- Encryption
  - password
  - session/cookies
  - Universal Unique Identifier (UUID)

#### Moderation

- Moderation System
- User access levels
