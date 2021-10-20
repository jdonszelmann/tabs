<script lang="ts">
    import {report_error, report_fetch_error} from "../typescript/Error";
    import {server_url} from "../typescript/Server";
    import {User, user} from "../typescript/User";
    import {useNavigate} from "svelte-navigator";

    let password: string;
    let passwordRepeat: string;

    let addUserUsername: string;
    let addUserPassword: string;

    let users: User[] = []

    const navigate = useNavigate();

    async function handleChangePassword() {
        if (password !== passwordRepeat) {
            report_error("passwords don't match")
            return;
        }
        if (password.length < 8) {
            report_error("password too short (less than 8)")
            return;
        }

        const resp = await fetch(`${server_url}/password`, {
            method: "PUT",
            body: JSON.stringify({
                Token: $user.Token,
                Password: password,
            })
        });

        if (!resp.ok) {
            await report_fetch_error(resp);
        } else {
            password = ""
            passwordRepeat = ""

            $user = null
        }
    }

    async function getUsers() {
        const resp = await fetch(`${server_url}/user/get-all`, {
            method: "POST",
            body: JSON.stringify({
                Token: $user.Token,
            })
        });

        if (!resp.ok) {
            await report_fetch_error(resp);
        } else {
            users = await resp.json()
        }
    }

    async function createUser(): Promise<void> {
        const resp = await fetch(`${server_url}/register`, {
            method: "POST",
            body: JSON.stringify({
                Token: $user.Token,
                Username: addUserUsername,
                Password: addUserPassword,
            })
        });

        if (!resp.ok) {
            await report_fetch_error(resp);
        } else {
            const json = await resp.json();
            users.push(json);
            users = users;
        }
    }

    async function deleteUser(name: string): Promise<void> {
        if (name == $user.Name) {
            if (!confirm("Are you sure you want to delete your own account?")) {
                return
            }
        } else {
            if (!confirm("Are you sure you want to delete this account?")) {
                return
            }
        }

        const resp = await fetch(`${server_url}/user`, {
            method: "DELETE",
            body: JSON.stringify({
                Token: $user.Token,
                Name: name,
            })
        });

        if (!resp.ok) {
            await report_fetch_error(resp);
        } else {
            if (name == $user.Name) {
                $user = null;
                navigate("/");
            }

            if ($user.Admin) {
                users = users.filter(i => i.Name !== name);
            }
        }
    }


    async function setAdmin(name: string, value: boolean): Promise<void> {
        const resp = await fetch(`${server_url}/admin`, {
            method: "PUT",
            body: JSON.stringify({
                Token: $user.Token,
                Name: name,
                Admin: value,
            })
        });

        if (!resp.ok) {
            await report_fetch_error(resp);
        } else {
            if (name == $user.Name) {
                $user.Admin = value;
            }

            for (const i of users) {
                if (i.Name === name) {
                    i.Admin = value;
                }
            }

            users = users;
        }
    }
</script>

<div class="account">
    <div class="content">
        <h1>Change Password</h1>
        <form on:submit|preventDefault={handleChangePassword}>
            <label for="password">Password</label>
            <input type="password" name="password" id="password" bind:value={password}>

            <label for="password-repeat">Repeat password</label>
            <input type="password" name="password-repeat" id="password-repeat" bind:value={passwordRepeat}>

            <button>Submit</button>
        </form>

        <button class="delete-account" on:click={deleteUser($user.Name)}>Delete Account</button>

        {#if $user !== null && $user.Admin}
            <h1>Users</h1>
            <div class="users">
                <div>Name</div><div>Admin</div><div></div>
                {#await getUsers() then _}
                    {#each users as u}
                        <div>{u.Name}</div>
                        <div>
                            {#if u.Name === $user.Name}
                                <input
                                    type="checkbox"
                                    disabled
                                    checked
                                />
                            {:else}
                                {#if u.Admin}
                                    <input
                                        type="checkbox"
                                        checked
                                        on:change={setAdmin(u.Name, false)}
                                    />
                                {:else}
                                    <input
                                        type="checkbox"
                                        on:change={setAdmin(u.Name, true)}
                                    />
                                {/if}
                            {/if}
                        </div>
                        <div>
                            <button on:click={deleteUser(u.Name)}>Delete User</button>
                        </div>
                    {/each}
                {/await}
            </div>

            <h1>Add User</h1>
            <form on:submit|preventDefault={createUser}>
                <label for="username">Username</label>
                <input name="username" id="username" bind:value={addUserUsername}>

                <label for="addUserPassword">Password</label>
                <input autocomplete="new-password" type="password" name="password" id="addUserPassword" bind:value={addUserPassword}>

                <button>Submit</button>
            </form>
        {/if}
    </div>
</div>

<style lang="scss">
  @import "src/global";

  .account {
    width: 100%;
    height: 100%;

    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;

    .delete-account {
      margin-top: 2em;
    }

    button {
      color: white;
      background-color: $darkred !important;
      border: 1px solid black;
      margin: 0;
      padding: .5em;
      width: 100%;

      &:hover {
        cursor: pointer;
        background-color: $red !important;
      }
    }

    .content {
      width: 60%;
      background-color: $brown;
      padding: 5em 3em;
      box-sizing: border-box;

      margin-top: 3em;

      .users {
        display: grid;
        grid-template-columns: auto auto auto;
        grid-auto-rows: auto;


        margin: 3em 1em 1em;

        width: 100%;

        * {
          padding: 1em;
          text-align: center;
          display: flex;
          flex-direction: column;
          justify-content: center;

          &:nth-child(6n+1), &:nth-child(6n+2), &:nth-child(6n+3) {
            background-color: $brown-a1;
          }

          &:nth-child(6n+4), &:nth-child(6n+5), &:nth-child(6n+6) {
            background-color: $gray;
          }
        }
      }

      h1 {
        text-align: center;
      }

      form {
        display: grid;

        grid-template-columns: 1fr 2fr;
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