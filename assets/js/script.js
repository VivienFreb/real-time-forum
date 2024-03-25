let inscription_page = document.getElementById("content_inscription");
        let login_page = document.getElementById("content_login");
        let dashboard_page = document.getElementById("content_dashboard");

        let registerBtn = document.getElementById('registrationBtn');
        let registerForm = document.getElementById('registrationForm');
        let loginBtn = document.getElementById('loginBtn');
        let loginForm = document.getElementById('loginForm');

        let register_username = document.getElementById('register_username')
        let register_email = document.getElementById('register_email')
        let register_password = document.getElementById('register_password')
        let register_confirm_password = document.getElementById('register_confirm_password')
        let login_username = document.getElementById('login_username');
        let login_password = document.getElementById('login_password')



    
        document.getElementById('showLoginPage').addEventListener('click', function() {
            inscription_page.classList.add("wrapper");
            login_page.classList.remove("wrapper");
            dashboard_page.classList.add("wrapper");
        });
        document.getElementById('showDashboard').addEventListener('click', function() {
            inscription_page.classList.add("wrapper");
            login_page.classList.add("wrapper");
            dashboard_page.classList.remove("wrapper");
        });
        document.getElementById('showRegistrationPage').addEventListener('click', function() {
            inscription_page.classList.remove("wrapper");
            login_page.classList.add("wrapper");
            dashboard_page.classList.add("wrapper");
        });
        document.getElementById('logout').addEventListener('click', function() {
            inscription_page.classList.add("wrapper");
            login_page.classList.remove("wrapper");
            dashboard_page.classList.add("wrapper");
        });
    
        const socket = new WebSocket('ws://localhost:8080/ws');

        socket.addEventListener('open', function (event) {
            console.log('Connexion WebSocket établie');
        });

        socket.addEventListener('message', function (event) {
            console.log('Message reçu du serveur:', event.data);
        });

        
        registerBtn.addEventListener('click', function (event) {
            event.preventDefault(); // Empêche la soumission normale du formulaire
            if (register_username.value.length <= 0 
            || register_email.value.length <= 0 
            || register_password.value.length <= 0){
                alert("L'un des champs requis est vide.")
                return
            }
            const formData = new FormData(registerForm);

            // Convertir FormData en objet JSON
            const data = {};
            formData.forEach((value, key) => {
                data[key] = value;
            });

            // Envoyer les données via WebSocket
            socket.send(JSON.stringify(data));
            document.getElementById('register_username').value = null;
            document.getElementById('register_email').value = null;
            document.getElementById('register_password').value = null;
            document.getElementById('register_confirm_password').value = null;
        });
        loginBtn.addEventListener('click', function (event) {
            event.preventDefault(); // Empêche la soumission normale du formulaire
            if (login_username.value.length <= 0 || login_password.value.length <= 0){
                alert("popup, userame or password not fill")
                return
            }
            const formData = new FormData(loginForm);

            // Convertir FormData en objet JSON
            const data = {};
            formData.forEach((value, key) => {
                data[key] = value;
            });

            // Envoyer les données via WebSocket
            socket.send(JSON.stringify(data));
            login_username.value = 'Username';
            login_password.value = 'Password';
        });