<script lang="ts">
    import {user} from "../typescript/User"
    import {report_fetch_error} from "../typescript/Error";
    import {useNavigate, useLocation, navigate} from "svelte-navigator";
    import {server_url} from "../typescript/Server";

    let username: string;
    let password: string;

    const navigate = useNavigate();
    const location = useLocation();

    async function handleLogin() {
        const resp = await fetch(`${server_url}/login`, {
            method: "POST",
            body: JSON.stringify({
                "Username": username,
                "Password": password
            })
        });

        if (!resp.ok) {
            await report_fetch_error(resp, false);
        } else {
            const json = await resp.json();

            $user = await json;

            const from = ($location.state && $location.state.from) || "/";
            navigate(from, { replace: true });
        }
    }
</script>

<div class="login">
    <div class="prompt">
        <h1>Login</h1>

        <form on:submit|preventDefault={handleLogin}>
            <label for="username">Username</label>
            <input name="username" id="username" bind:value={username}>

            <label for="password">Password</label>
            <input type="password" name="password" id="password" bind:value={password}>

            <button>Submit</button>
        </form>
    </div>
</div>

<style lang="scss">
    @import "src/global";

    .login {
      width: 100%;
      height: 100%;

      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;

      .prompt {
        background-color: $darkred;
        color: white;

        padding: 2em;
        h1 {
          text-align: center;
        }

        form {
          display: grid;

          grid-template-columns: auto auto;
          grid-auto-rows: auto;

          gap: 1em;
          * {
            height: 100%;
          }
          label {
            display: flex;
            flex-direction: column;
            justify-content: center;
          }

          input {
            width: 100%;
            background-color: $red;
            border: 1px solid black;
            color: white;
          }

          button {
            grid-column: 1/3;

            color: white;
            background-color: $darkred;
            border: 1px solid black;
            margin: 0;

            width: 100%;

            &:hover {
              cursor: pointer;
              background-color: $red;
            }
          }
        }
      }
    }
</style>