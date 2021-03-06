Goal: Develop a Go start auth starter [without Gin
framework](https://github.com/auth0-samples/auth0-golang-web-app/tree/master/01-Login)
and learn along the way ...

https://www.youtube.com/watch?v=bpa_GQT16uM&feature=youtu.be

https://twitter.com/kaihendry/status/1482526555402571776

https://manage.auth0.com/dashboard/us/dev-h3aod060/users/YXV0aDAlN0M2MWUzN2JiNDY5MGNkMTAwNjg2Zjg0ZTI

Start dynamodb server

    ./scripts/local-dynamodb.sh
    ./scripts/create-table.sh

Start Go Web server

    ./scripts/start-local-server.sh

If you like this, check out https://github.com/kaihendry/local-audio which builds on this.

# Questions for Auth0

https://s.natalian.org/2022-01-16/auth-forum.mp4

1. How does one login via Github / Google Gmail address?

You have to choose "Continue with Google", not 100% sure how to enable the others.

https://manage.auth0.com/dashboard/us/hendry/connections/social/create/github

2. Is "sub" the correct way to identify a user?

Yes, it's short for **subject**.

3. What happens if user the logs in from Github/Gmail etc... are the identities merged?

You should create some sort of uuid translation table. I.e. user opts in to map his/her Gmail/Github subject to your system's uuid.

4. Did I correctly handle / validate the profile in [middleware](middleware.go)?

Yes, but it's not optimal in the sense it cannot be fine grained. See https://youtu.be/cw0PjurUiDQ

5. How can I be confident the profile (cookie?) has not been manipulated?

The cookies are signed by the server, they are tamper proof

6. Auth0's **profile** appears to be an interface, is there a proper struct for it?
7. The profile also has a clumsy was to access it as can be seen in https://youtu.be/bpa_GQT16uM?t=1672

8. The default UX with the browser is poor, Chrome generally offers [the wrong password](https://s.natalian.org/2022-01-16/wrong-password.png), why is that?
9. The Auth0 server is in the US, could it not be in Singapore?
10. Logout doesn't clear the cookie, why?
11. Can I tell how many users are logged in at any one time, and who they are?

# TODO

Split out into a Lambda authorizer
