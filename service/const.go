package service

const dashBoardPage = `<html><body>

 {{if .Username}}
    <p>logged as:<b>{{.Username}}</b>, welcome to epfiot PoC </p>
	<ul>
	  <li><a href="/dashboard">dashboard</a></li>
	  <li><a href="/logout">Logout!</a></li>
	</ul>
 {{else}}
         <p>Either your session has expired or you've logged out! <a href="/login">Login</a></p>
 {{end}}

 </body></html>`


const logUserPage = `<html><body>
 {{if .LoginError}}<p style="color:red">Either username or password is not in our record! Sign Up?</p>{{end}}

 <form method="post" action="/login">
         {{if .Username}}
                  <p><b>{{.Username}}</b>, you're already logged in! <a href="/logout">Logout!</a></p>
         {{else}}
                 <label>Username:</label>
                 <input type="text" name="Username"><br>

                 <label>Password:</label>
                 <input type="password" name="Password">

                 <span style="font-style:italic"> Enter: 'mynakedpassword'</span><br>
                 <input type="submit" name="Login" value="Let me in!">
         {{end}}
 </form>
 </body></html>`

const mainPage = `<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
        {{if .Username}}
		    <div id="graphiql" style="height: 100vh;">Loading...</div>
            <script>
                function graphQLFetcher(graphQLParams) {
                    return fetch("/query", {
                        method: "post",
                        body: JSON.stringify(graphQLParams),
                        credentials: "include",
                    }).then(function (response) {
                        return response.text();
                    }).then(function (responseBody) {
                        try {
                            return JSON.parse(responseBody);
                        } catch (error) {
                            return responseBody;
                        }
                    });
                }

                ReactDOM.render(
                    React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
                    document.getElementById("graphiql")
                );
            </script>
        {{else}}

         <p>Either your session has expired or you've logged out! <a href="/login">Login</a></p>

        {{end}}
	</body>
</html>
`
