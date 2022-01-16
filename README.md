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
2. Is "sub" the correct way to identify a user?
3. What happens if user the logs in from Github/Gmail etc... are the identities merged?
4. Did I correctly handle / validate the profile in [middleware](middleware.go)?
5. How can I be confident the profile hasn't been manipulated?
6. Profile appears to be an interface, is there a proper struct for it?
7. The default UX with the browser is poor, Chrome generally offers [the wrong password](https://s.natalian.org/2022-01-16/wrong-password.png), why is that?
8. The Auth0 server is in the US, could it not be in Singapore?
