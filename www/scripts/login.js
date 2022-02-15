document.addEventListener("DOMContentLoaded", () => {
    document.getElementById('user').addEventListener('keyup', event => {
	if (event.key == 'Enter') {
	    login();
	}
    });

    document.getElementById('password').addEventListener('keyup', event => {
	if (event.key == 'Enter') {
	    login();
	}
    });
});

function login() {
    postJSON('login', getFields(), onLogin, onLoginError);
}

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
