document.addEventListener('DOMContentLoaded', function() {
    const form = document.querySelector('form');
    // const spinnerContainer = document.getElementById('spinner-container');
    // const actualSpinner = document.querySelector('.actual-spinner');

    if (form) {
        form.addEventListener('submit', async function(event) {
            event.preventDefault();

            // Clear previous error messages
            document.querySelectorAll('.error-message').forEach(el => el.textContent = '');

            // Client-side validation
            let isValid = true;

            // Full Name validation
            const fullName = document.getElementById('fullName');
            if (fullName.value.trim() === '') {
                document.getElementById('fullNameError').textContent = 'Full Name is required';
                isValid = false;
            }

            // Location validation
            const location = document.getElementById('location');
            if (location.value.trim() === '') {
                document.getElementById('locationError').textContent = 'Location is required';
                isValid = false;
            }

            // Email validation
            const email = document.getElementById('email');
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (email.value.trim() === '') {
                document.getElementById('emailError').textContent = 'Email is required';
                isValid = false;
            } else if (!emailRegex.test(email.value)) {
                document.getElementById('emailError').textContent = 'Invalid email format';
                isValid = false;
            }

            // Photos validation
            const photos = document.getElementById('photos');
            if (photos.files.length === 0) {
                document.getElementById('photosError').textContent = 'Please upload at least one photo';
                isValid = false;
            }

            // Size validation
            const size = document.getElementById('size');
            if (size.value === '') {
                document.getElementById('sizeError').textContent = 'Please select a print size';
                isValid = false;
            }

            // Paper Type validation
            const paperType = document.getElementById('paperType');
            if (paperType.value === '') {
                document.getElementById('paperTypeError').textContent = 'Please select a paper type';
                isValid = false;
            }

            // If client-side validation fails, stop submission
            if (!isValid) {
                return;
            }

              // Display spinner on form submission
            //   spinnerContainer.classList.remove('spinner-hidden');
            //   spinnerContainer.style.display = 'flex';

            // Create FormData
            const formData = new FormData(this);

            try {
                const response = await fetch('/submit-order', {
                    method: 'POST',
                    body: formData
                });

                const data = await response.json();

                if (!response.ok) {
                    // Handle server-side validation errors
                    if (data.errorField) {
                        
                        const errorElement = document.getElementById(`${data.errorField.toLowerCase()}Error`);
                        if (errorElement) {
                            errorElement.textContent = data.error;
                        }
                    } else {
                        // Generic error handling
                        alert(data.error || 'An error occurred');
                    }
                    return;
                }

                // Successful submission
                alert('Order submitted successfully!');
                this.reset(); // Clear the form
            } catch (error) {
                console.error('Submission error:', error);
                alert('Network error. Please try again.');
            } 
            // finally {
            //     // Hide the spinner after processing
            //     // spinnerContainer.classList.add('spinner-hidden')
            // }
        });
    }
    // spinner.addEventListener("transitionend", () => {
    //     if (spinnerContainer.classList.contains('spinner-hidden')) {
    //         spinnerContainer.style.display = 'none';
    //     }
    // });
});

