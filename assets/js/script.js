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
        let sendMessageBtn = document.getElementById('sendMessageBtn')
        const messageText = document.getElementById('messageText')


        let currentUser;
        let currentOther;



    
        document.getElementById('showLoginPage').addEventListener('click', function() {
            inscription_page.classList.add("hidden");
            login_page.classList.remove("hidden");
            dashboard_page.classList.add("hidden");
        });
        document.getElementById('showRegistrationPage').addEventListener('click', function() {
            inscription_page.classList.remove("hidden");
            login_page.classList.add("hidden");
            dashboard_page.classList.add("hidden");
        });
        // document.getElementById('logout').addEventListener('click', function() {
        //     inscription_page.classList.add("hidden");
        //     login_page.classList.remove("hidden");
        //     dashboard_page.classList.add("hidden");
        // });

        function initWebSocket() 
        {
            const socket = new WebSocket('ws://localhost:8080/ws');
        
            socket.addEventListener('open', function (event) {
                console.log('WebSocket connection established');
                rebootStatus()
            });
        
            return socket;
        }

    
        const socket = initWebSocket();

        socket.addEventListener('open', function (event) {
            console.log('Connexion WebSocket établie');
        });
        
        registerBtn.addEventListener('click', function (event) {
            event.preventDefault(); // Empêche la soumission normale du formulaire
            if (register_username.value.length <= 0 
            || register_email.value.length <= 0 
            || register_password.value.length <= 0){
                alert("L'un des champs requis est vide.")
                return
            }
            console.log(registerForm);
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
            document.getElementById('User-Verif').classList.add('hidden')
            document.getElementById('content_dashboard').classList.remove('hidden')
            // displayPosts()
            // socket.close()
            // initWebSocket()
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
            data.username = login_username.value

            // Envoyer les données via WebSocket
            socket.send(JSON.stringify(data));
            currentUser = login_username.value;

            const sidebar = document.querySelector('.sidenav')
            sidebar.classList.remove('hidden')

            login_username.value = '';
            login_password.value = '';
            getUsers()
            requestsPosts()
        });

        socket.addEventListener('message', function(event){
            const response = JSON.parse(event.data)
            console.log(response);
            if (response && response.Name){
                if (response.Name == "Login")
                     if (response.Success){
                        document.getElementById('User-Verif').classList.add('hidden')
                        document.getElementById('content_dashboard').classList.remove('hidden')
                        console.log("Coucou, je rentre dans LOGIN ");
                     }
                if (response.Name === "Friends"){
                    console.log("This is the friends list:\n",response.Ray);
                    const friendList = response.Users;
                    const me = document.createElement('h4')
                    me.textContent = currentUser
                    me.style = "color: blue"
                    document.getElementById('MyName').appendChild(me)
                    friendList.forEach(user =>{
                        const li = document.createElement('li')
                        li.textContent = user.username
                        li.style.color = "red"
                        document.getElementById('userList').appendChild(li)

                        li.addEventListener('click', function(){
                            console.log("You clicked!");
                            currentOther = li.textContent
                            getChat()
                            document.querySelector('.dashboard').classList.add('hidden')
                            document.querySelector('.container').classList.remove('hidden')
                        })
                    })
                }
                if (response.Name === "userStatus"){
                    // console.log(response);
                    updateUserStatus(response.Checks)
                } 
                if (response.Name === "chatHistory"){
                    generateChat(response.Chats)
                }
            } 
        })

        function displayPosts(posts) {
            const dashboardContent = document.getElementById('content_dashboard')
            const postsList = document.createElement('ul')
            posts.forEach(post =>{
                const listItem = document.createElement('li')
                listItem.textContent = `${post.Title} - ${post.User_id} - ${post.Description}`;
                postsList.appendChild(listItem)
            })
            dashboardContent.innerHTML = '';
            dashboardContent.appendChild(postsList)
        }

        function generateChat(discs){
            console.log(discs)
            const chatContainer = document.getElementById('chatContainer')
            // chatContainer.innerHTML = ''
            // createSendDiv()

            discs.forEach(message =>{
                const chatBubble = document.createElement('div')
                chatBubble.classList.add('container')
                if (message.Speaker === currentUser){
                    chatBubble.classList.add('darker')
                }
                chatBubble.innerHTML = `<p><strong>${message.Speaker}:</strong> ${message.Content}</p>`;
                chatContainer.appendChild(chatBubble)
            })
            chatContainer.scrollTop = chatContainer.scrollHeight
        }

        function updateUserStatus(userStatus){
            const userList = document.getElementById('userList')
            userList.innerHTML = ''
            const newline = document.createElement('br')
            userList.appendChild(newline)
            console.log(userStatus);

            userStatus.forEach(user =>{
                const li = document.createElement('li');
                li.textContent = `${user.Name}`;
                li.style.color = user.Status === 'active' ? 'green' : 'red';
                userList.appendChild(li);

                li.addEventListener('click', function(){
                    console.log("You clicked!");
                    currentOther = li.textContent
                    getChat()
                    document.querySelector('.dashboard').classList.add('hidden')
                    document.querySelector('.container').classList.remove('hidden')
                })
            })
        }

        function clearChatHistory(){
            const chatContainer = document.getElementById('chatContainer')
            const containers = chatContainer.getElementsByClassName('container')
            while (containers.length > 0) {
                containers[0].parentNode.removeChild(containers[0]);
            }
        }

        sendMessageBtn.addEventListener('click', function(){
            console.log('Etape 1')
            const messageContent = messageText.value.trim();
            if (messageContent !== ''){
                const message = {
                    FormName:"chatEnvoy",
                    "Username": `${currentUser}`,
                    "Other": `${currentOther}`,
                    Content: messageContent
                }
                socket.send(JSON.stringify(message))
            }
            messageText.value = ''
            clearChatHistory()
            getChat()
        })

        const requestsPosts = () => {message = {FormName: "posts"};socket.send(JSON.stringify(message))}
        const rebootStatus = () => {message = {FormName: "reset"};socket.send(JSON.stringify(message))}
        const getUsers = () => {message = {FormName:"usershunt", "Username":`${currentUser}`};socket.send(JSON.stringify(message))}
        const getStatus = () => {message = {FormName:"userStatus", "Username":`${currentUser}`};socket.send(JSON.stringify(message))}
        const getChat = () => {message = {FormName:"discussions", "Username":`${currentUser}`,"Other":`${currentOther}`};socket.send(JSON.stringify(message))}
        setInterval(getStatus, 10000)
        // getStatus()