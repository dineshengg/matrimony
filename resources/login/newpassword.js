document.addEventListener('DOMContentLoaded', function() {
    const form = document.querySelector('form');
    const passwordInput = document.getElementById('password');
    const confirmPasswordInput = document.getElementById('confirmpassword');
    const guidElement = document.getElementById('guid');
    const errorField = document.getElementById('error-field');
    const errorLabel = document.getElementById('error');

    // Get GUID from the hidden element
    const guid = guidElement ? guidElement.textContent.trim() : '';

    // Password validation regex: 8+ chars, includes A-Z, a-z, 0-9, and special characters
    const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@#$%^&*])[A-Za-z\d@#$%^&*]{8,}$/;

    function showError(message) {
        errorLabel.textContent = 'Error: ' + message;
        errorField.style.display = 'block';
        errorLabel.hidden = false;
    }

    function hideError() {
        errorField.style.display = 'none';
        errorLabel.hidden = true;
    }

    // Form submission handler
    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        hideError();

        const password = passwordInput.value;
        const confirmPassword = confirmPasswordInput.value;

        // Validation checks
        if (!guid) {
            showError('Invalid reset link. Please request a new password reset.');
            return;
        }

        if (password.length < 8) {
            showError('Password must be at least 8 characters long.');
            return;
        }

        if (!passwordRegex.test(password)) {
            showError('Password must include uppercase, lowercase, number, and special character (@#$%^&*).');
            return;
        }

        if (password !== confirmPassword) {
            showError('Passwords do not match.');
            return;
        }

        // Prepare form data
        const formData = new FormData();
        formData.append('guid', guid);
        formData.append('newpassword', password);
        formData.append('confirmpassword', confirmPassword);

        // Construct the URL with guid as query parameter
        const url = `/api/noauth/reset-password?guid=${encodeURIComponent(guid)}`;

        try {
            // Disable submit button to prevent double submission
            const submitButton = form.querySelector('button[type="submit"]');
            submitButton.disabled = true;
            submitButton.textContent = 'Submitting...';

            // Send POST request
            const response = await fetch(url, {
                method: 'POST',
                body: formData
            });

            const responseText = await response.text();

            if (response.ok) {
                // Success - show success message and redirect
                showError('Password has been reset successfully! Redirecting to login page...');
                window.location.href = '/api/noauth/login';
            } else {
                // Error response from server
                showError(responseText || 'Failed to reset password. Please try again.');
                submitButton.disabled = false;
                submitButton.textContent = 'Submit';
            }
        } catch (error) {
            console.error('Error resetting password:', error);
            showError('Network error. Please check your connection and try again.');
            
            // Re-enable submit button
            const submitButton = form.querySelector('button[type="submit"]');
            submitButton.disabled = false;
            submitButton.textContent = 'Submit';
        }
    });

    // Real-time password strength indicator (optional enhancement)
    passwordInput.addEventListener('input', function() {
        const password = this.value;
        if (password.length > 0 && !passwordRegex.test(password)) {
            this.style.borderColor = '#ff6b6b';
        } else if (password.length >= 8) {
            this.style.borderColor = '#51cf66';
        } else {
            this.style.borderColor = '';
        }
    });

    // Real-time password match indicator
    confirmPasswordInput.addEventListener('input', function() {
        const password = passwordInput.value;
        const confirmPassword = this.value;
        
        if (confirmPassword.length > 0 && password !== confirmPassword) {
            this.style.borderColor = '#ff6b6b';
        } else if (confirmPassword.length > 0 && password === confirmPassword) {
            this.style.borderColor = '#51cf66';
        } else {
            this.style.borderColor = '';
        }
    });
});
