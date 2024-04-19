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
        const sidebar = document.querySelector('.sidenav')
        let logoutBtn = document.getElementById('logout')
        let homeBtn = document.getElementById('homeBtn')
        let chatContainer = document.getElementById('chatContainer')
        let currentUser;
        let currentOther;

        logoutBtn.addEventListener('click', function(){
            dashboard_page.classList.add('hidden')
            sidebar.classList.add('hidden')
            document.getElementById('User-Verif').classList.remove('hidden')
            document.getElementById('chatContainer').classList.add('hidden')
            document.getElementById('MyName').innerHTML = ''
            Empty('userList')
            document.getElementById('userList').appendChild(document.createElement('br'))
            unlog()
        })

        homeBtn.addEventListener('click', function(){
            if (dashboard_page.classList.contains('hidden')){
                document.querySelector('.dashboard').classList.remove('hidden')
                document.querySelector('.container').classList.add('hidden')
            }
        })

        function Empty(arg){
            const Arg = document.getElementById(arg)
            while (Arg.firstChild) Arg.removeChild(Arg.firstChild)
        }

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
            event.preventDefault(); // => Empêche la soumission du formulaire.
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
            data.username = register_username.value
            currentUser = register_username.value;

            //Envoi des données dans le serveur interne
            socket.send(JSON.stringify(data));
            document.getElementById('register_username').value = null;
            document.getElementById('register_email').value = null;
            document.getElementById('register_password').value = null;
            document.getElementById('register_confirm_password').value = null;
            requestsPosts()
            getUsers()
            document.getElementById('User-Verif').classList.add('hidden')
            document.getElementById('content_dashboard').classList.remove('hidden')
            sidebar.classList.remove('hidden')
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
            requestsPosts()
            getUsers()
        });

        socket.addEventListener('message', function(event){
            const response = JSON.parse(event.data)
            console.log(response);
            if (response && response.Name){
                if (response.Name == "Login") {
                    if (response.Success){
                    document.getElementById('User-Verif').classList.add('hidden')
                    document.getElementById('content_dashboard').classList.remove('hidden')
                    console.log("Coucou, je rentre dans LOGIN ");
                     }
                }                
                if (response.Name === "Friends"){
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
                            console.log(li.textContent);
                            if (currentOther !== li.textContent){console.log("You clicked!");
                            if (chatContainer.classList.contains('hidden')) chatContainer.classList.remove('hidden')
                            currentOther = li.textContent
                            getChat()
                            document.querySelector('.dashboard').classList.add('hidden')
                            document.querySelector('.container').classList.remove('hidden')}
                        })
                    })
                }
                if (response.Name === "userStatus"){
                  updateUserStatus(response.Checks)
                } 
                if (response.Name === "chatHistory"){
                    generateChat(response.Chats)
                }
                if (response.Name === "Posts"){
                    console.log(response.Posts);
                    displayPosts(response.Posts)
                }
                if (response.Success){
                    document.getElementById('User-Verif').classList.add('hidden')
                    document.getElementById('content_dashboard').classList.remove('hidden')
                    console.log("Coucou, je rentre dans LOGIN ");
                 }
            } 
        })

        function displayPosts(posts) {
            const dashboardContent = document.getElementById('content_dashboard')
            const postsList = document.createElement('ul')
            posts.forEach(post => {
                const listItem = document.createElement('li');
                listItem.textContent = `${post.Title} - ${post.User_name} - ${post.Description}`;
                listItem.className = "post"

                if (post.Comments && post.Comments.length > 0) {
                    const commentsList = document.createElement('ul');
                    post.Comments.forEach(comment => {
                        console.log("Commentaire:",comment);
                        const commentItem = document.createElement('li');
                        commentItem.textContent = `${comment.Username} said: ${comment.Content}`;
                        commentsList.appendChild(commentItem);
                    });
                    listItem.appendChild(commentsList);
                }
            
                postsList.appendChild(listItem);
            });
            dashboardContent.innerHTML = '';
            dashboardContent.innerHTML += postCode
            dashboardContent.innerHTML += commentCode
            dashboardContent.appendChild(postsList)
            const comment = document.getElementById('quickComment')
            const reply = document.getElementById('quickReply')

            dashboardContent.addEventListener('click', function(event) {
                const target = event.target;
                if (target && target.classList.contains('post')) {
                    comment.classList.remove('hidden')
                }
            });

            const postBtn = document.getElementById('postBtn')
            postBtn.addEventListener('click',function(){
                const formData = {
                    formName: 'postCreation',
                    "Username": currentUser,
                    "Subject": document.querySelector('#quickReply input[name="name"]').value,
                    "Content": document.querySelector('#quickReply textarea[name="com"]').value
                }
                socket.send(JSON.stringify(formData))
                document.querySelector('#quickReply input[name="name"]').value = ''
                document.querySelector('#quickReply textarea[name="com"]').value = ''
                reply.classList.add('hidden')
                requestsPosts()
            })
            const commentBtn = document.getElementById('commentBtn')
            commentBtn.addEventListener('click', ()=>{
                const formData = {
                    formName: 'commentCreation',
                    "Username": currentUser,
                    "Subject": "Placeholder",
                    "Content": document.querySelector('#quickComment textarea[name="com"]').value
                }
                socket.send(JSON.stringify(formData))
                document.querySelector('#quickComment input[name="name"]').value = ''
                document.querySelector('#quickComment textarea[name="com"]').value = ''
                comment.classList.add('hidden')
                // requestsPosts()
            })
        }

        function generateChat(discs){
            console.log(discs)
            clearChatHistory()
            const chatContainer = document.getElementById('chatContainer')
            if (discs != null) {
                discs.forEach(message =>{
                const chatBubble = document.createElement('div')
                chatBubble.classList.add('container')
                if (message.Speaker === currentUser){
                    chatBubble.classList.add('darker')
                }
                chatBubble.innerHTML = `<p><strong>${message.Speaker}:</strong> ${message.Content}</p>`;
                chatContainer.appendChild(chatBubble)
                })
            }
            chatContainer.scrollTop = chatContainer.scrollHeight
        }

        function updateUserStatus(userStatus){
            const userList = document.getElementById('userList')
            userList.innerHTML = ''
            const newline = document.createElement('br')
            userList.appendChild(newline)
            console.log(userStatus);
            requestsPosts()

            userStatus.forEach(user =>{
                const li = document.createElement('li');
                li.textContent = `${user.Name}`;
                li.style.color = user.Status === 'active' ? 'green' : 'red';
                userList.appendChild(li);

                li.addEventListener('click', function(){
                    if (currentOther !== li.textContent){console.log("You clicked!",this.textContent);
                    currentOther = li.textContent
                    getChat()
                    document.querySelector('.dashboard').classList.add('hidden')
                    document.querySelector('.container').classList.remove('hidden')}
                    if (chatContainer.classList.contains('hidden')) chatContainer.classList.remove('hidden');
                    
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
        const unlog = () => {message = {FormName:"delog", "Username":`${currentUser}`};socket.send(JSON.stringify(message))}
        setInterval(CheckIfNeeded,5000)

        function CheckIfNeeded(){
            if (!sidebar.classList.contains('hidden')) getStatus()
        }
        
        const postCode = `
        <div id="quickReply" class="extPanel  hidden" style="right: 0px; top: 10%;">
            <div id="qrForm">
                <div><input name="name" type="text" placeholder="Subject"></div>
                <div><textarea name="com" cols="48" rows="4" placeholder="Comment"></textarea></div>
                <button id="postBtn">Post</button>
            </div>
        </div>
    `;
            
        const commentCode = `
            <div id="quickComment" class="extPanel reply hidden" style="right: 0px; top: 30%;">
                <div id="qrForm">
                    <div><textarea name="com" cols="48" rows="4" placeholder="Comment"></textarea></div>
                    <button id="commentBtn">Comment</button>
                </div>
            </div>
        `;