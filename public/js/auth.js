const API_URL = 'http://localhost:3000/api/v1';
let token = localStorage.getItem('token');

// Notyf settings
const notyf = new Notyf({
    duration: 3000,
    position: { x: 'right', y: 'top' },
    types: [
        {
            type: 'success',
            background: '#2ecc71',
        },
        {
            type: 'error',
            background: '#e74c3c',
        }
    ]
});

// Check authentication status on load
document.addEventListener('DOMContentLoaded', () => {
    if (token) {
        checkAuth();
    }
});

async function handleLogin(event) {
    event.preventDefault();
    toggleLoading(true);

    try {
        const email = document.getElementById('login-email').value;
        const password = document.getElementById('login-password').value;

        const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
        });

        const data = await response.json();
        if (data.success) {
            handleAuthSuccess(data);
            notyf.success('Login successful!');
        } else {
            notyf.error(data.message || 'Login failed');
        }
    } catch (error) {
        notyf.error('Login failed. Please try again.');
    } finally {
        toggleLoading(false);
    }
}

async function handleRegister(event) {
    event.preventDefault();
    const name = document.getElementById('register-name').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;

    try {
        const response = await fetch(`${API_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name, email, password }),
        });
        const data = await response.json();
        if (data.success) {
            handleAuthSuccess(data);
        } else {
            alert(data.message);
        }
    } catch (error) {
        alert('Registration failed');
    }
}

function handleOAuth(provider) {
    window.location.href = `${API_URL}/auth/oauth/${provider}`;
}

function handleAuthSuccess(data) {
    token = data.data.token;
    localStorage.setItem('token', token);
    checkAuth();
}

async function checkAuth() {
    try {
        const response = await fetch(`${API_URL}/users/me`, {
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });
        const data = await response.json();
        if (data.success) {
            showProtectedContent(data.data);
        } else {
            showAuthForms();
        }
    } catch (error) {
        showAuthForms();
    }
}

function handleLogout() {
    localStorage.removeItem('token');
    token = null;
    showAuthForms();
}

function showProtectedContent(user) {
    document.getElementById('auth-forms').classList.add('hidden');
    document.getElementById('protected-content').classList.remove('hidden');
    document.getElementById('user-name').textContent = user.name;
    document.getElementById('user-info').innerHTML = `
        <p>Email: ${user.email}</p>
        <p>Role: ${user.role}</p>
    `;
}

function showAuthForms() {
    document.getElementById('auth-forms').classList.remove('hidden');
    document.getElementById('protected-content').classList.add('hidden');
}

function toggleForms() {
    document.getElementById('login-form').classList.toggle('hidden');
    document.getElementById('register-form').classList.toggle('hidden');
}