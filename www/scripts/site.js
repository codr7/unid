function postJSON(url, body, onSuccess, onError) {
    fetch(url, {
	method: 'POST',
	headers: {'Content-Type': 'application/json'},
	body: JSON.stringify(body)
    }).then((response) => {
	if (response.ok) {
	    onSuccess(response);
	} else {
	    onError(response);
	}
    });
}
