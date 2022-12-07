A basic Laravel / PHP sample to validate the JWT token. [See all the detail here](https://workvivo.notion.site/JWT-authenticator-81ab55ea36774cdfa3eac6ac26db3da9)

## Environment variables

These are the environment variables used in the sample code.


* WORKVIVO_APP_URL

The URL to redirect to, i.e., this is the Workvivo protected end point, e.g https://acme.workvivo.co/social


* WORKVIVO_APP_URL_PUBLIC_KEYS
  
The URL where Workvivo publishes their public keys, e.g. https://acme.workvivo.co/.well-known/public-keys

* RELAY_URL
  
The URL to return back to once Workvivo completes authentication, e.g, https://thirdpartyapp.workvivo.co/dashboard

* FORCE_SSL=TRUE

Let' force SSL across the board


* WORKVIVO_ISSUE_HOST

The host used to validate JWT ISS attribute, e.g. acme.workvivo.co
