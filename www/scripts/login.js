function getFields() {
    return {user: document.getElementById('user').value,
	    password: document.getElementById('password').value}
}

function onLogin(response) {
    window.location.replace('rcs.html');
}

function onLoginError(response) {
    response.text().then(text => {
	alert(text);
	window.location.replace('login.html');
    });
}
