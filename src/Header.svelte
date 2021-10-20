<script lang="ts">
    import {user} from "./typescript/User";
    import {link, useNavigate, useLocation} from "svelte-navigator";
    import {server_url} from "./typescript/Server";
    import {report_fetch_error} from "./typescript/Error";
    import {ServerTab} from "./typescript/ServerTab";

    const location = useLocation();
    const navigate = useNavigate();

    function logout() {
        navigate("/");
        $user = null;
    }

    async function newTab() {
        const resp = await fetch(`${server_url}/tab/new`, {
            method: "POST",
            body: JSON.stringify({
                Token: $user.Token
            })
        });

        if (!resp.ok) {
            await report_fetch_error(resp);
        } else {
            const data: ServerTab = await resp.json();
            navigate(`/edit/${data.Id}`)
        }
    }
</script>

<header>
    <a href="/" use:link><h2>Home</h2></a>
    {#if $user === null}
        <a href="/login" use:link><h2>Login</h2></a>
    {/if}
    {#if $user !== null}
        <a href="/tabs" use:link><h2>Tabs</h2></a>
    {/if}

    {#if $location.pathname === "/tabs"}
        <a on:click={newTab}>
            <h2>New Tab</h2>
        </a>
    {/if}

    <div class="separator"></div>

    {#if $user !== null}
        <a on:click={logout}>
            Logout
        </a>
        <a href="/account" class="logged-in" use:link>
            Logged in as {$user.Name}
        </a>
    {/if}
</header>

<style lang="scss">
  @import "global";

  .logged-in {
    color: white;
    margin-right: 1em;
  }

  .separator {
    margin-left: auto;
    margin-right: auto;
  }

  header {
    background-color: $darkred;
    display: flex;
    flex-direction: row;
    align-items: center;

    a {
      text-decoration: none !important;

      margin-left: 1em;
      color: white;

      &:hover {
        cursor: pointer;
      }

      h2 {
        margin: 0;
        padding: 0;
      }
    }
  }
</style>