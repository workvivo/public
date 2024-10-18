/* 
This code will request CORS to the Workvivo server.
You will need to ensure that you have enabled CORS in your organisation.
See https://developer.workvivo.com
*/
const workvivoJWT = "{{ $jwt }}";
const workvivoURL = "{{ $url }}";

function fetchSSOConnect() {
    fetch(workvivoURL, {
        method: "GET",´
        headers: {
            "x-workvivo-jwt": workvivoJWT, // set the workvivo cookie via Header
        },
        credentials: "include" // required to access cookies
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        if (response.redirected) {
            window.location.replace(response.url);
        }
    })
    .catch(error => {
        console.error('There was a problem with the fetch operation:', error);
    });
}