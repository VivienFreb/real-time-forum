// document.addEventListener('DOMContentLoaded', function(){
// const register = document.getElementById('register')
// const login = document.getElementById('login')
// })

document.addEventListener('DOMContentLoaded', function(){
    function toggleForm(form){
        document.getElementById('register').classList.add('hidden')
        document.getElementById('login').classList.add('hidden')
        document.getElementById(form).classList.remove('hidden')
    }
    
    function validateRegistrationForm(){
        console.log("Entering here! :D");
        var username = document.getElementById('register-username').value;
        var email = document.getElementById('register-email').value;
        var password = document.getElementById('register-password').value;
        var confirmPassword = document.getElementById('register-confirm-password').value;
        console.log(username,email,password,confirmPassword);
    
        if ((username) === ''|| (email) === '' || (password) === '' ){
            alert('Certains champs sont incomplets.')
            return false
        }
    
        if (password !== confirmPassword){
            alert('Les mots de passe sont divergents.')
            return false
        }
    
        return true
    }
    
    function validateLoginForm(){
        var username = document.getElementById('register-username').value;
        var email = document.getElementById('register-email').value;
        
        if (username === ''|| password === '' ){
            alert('Certains champs sont incomplets.')
            return false
        }
        return true
    }
    
    document.getElementById('login').addEventListener('submit', function(event){
        // event.preventDefault();
        validateLoginForm();
    })

    // document.getElementById('register-form').addEventListener('submit', function(event){
    //     // event.preventDefault();
    //     toggleForm('login');
    // })

    function NOPE(){
        document.getElementById('register').classList.add('hidden')
    }

    // document.getElementById('register-form').addEventListener('submit', function(event) {
    //     console.log("Entering here.");
    //     event.preventDefault(); // Prevent the default form submission
    
    //     if (validateRegistrationForm()) { // Validate the registration form
    //         const formData = new FormData(this); // Create FormData object
    
    //         // Make a POST request to the /register endpoint
    //         fetch('/register', {
    //             method: 'POST',
    //             body: formData
    //         })
    //         .then(response => {
    //             if (response.ok) {
    //                 console.log("Registration successful");
    //                 // Redirect or show success message
    //             } else {
    //                 console.error("Failed to register");
    //                 // Handle error response
    //             }
    //         })
    //         .catch(error => {
    //             console.error("An error occurred:", error);
    //             // Handle network error
    //         });
    //     }
    // });
    
})