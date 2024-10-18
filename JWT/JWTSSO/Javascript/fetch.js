const workvivoJWT = "{{ $jwt }}";
const workvivoURL = "{{ $url }}";

function fetchSSOConnect() {
    fetch(workvivoURL, {
        method: "GET",
        headers: {
            "x-workvivo-jwt": workvivoJWT,
        },
        credentials: "include" // Include cookies in the request
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