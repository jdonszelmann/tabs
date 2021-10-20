<script lang="ts">
	import TabEditor from "./editor/TabEditor.svelte";
	import {Route, Router} from "svelte-navigator";
	import Login from "./user/Login.svelte";
	import Account from "./user/Account.svelte";
	import PrivateRoute from "./private/PrivateRoute.svelte";
	import Tabs from "./Tabs.svelte";
	import Header from "./Header.svelte";
	import {server_url} from "./typescript/Server";
	import {user} from "./typescript/User";
	import {report_fetch_error} from "./typescript/Error";
	import {ServerTab} from "./typescript/ServerTab";
	import Tab from "./Tab.svelte";
	import PublicTabs from "./PublicTabs.svelte";
	import {error_message} from "./typescript/Error";

	async function getTab(id): Promise<ServerTab> {
		const resp = await fetch(`${server_url}/tab/get`, {
			method: "POST",
			body: JSON.stringify({
				Token: $user && $user.Token,
				Id: id,
			})
		});

		if (!resp.ok) {
			await report_fetch_error(resp);
		} else {
			return await resp.json();
		}
	}
</script>

<div class="root">
	<Router>
		<Header />

		<main>
			<Route path="/">
				<PublicTabs />
			</Route>

			<Route path="/login">
				<Login />
			</Route>

			<PrivateRoute path="/tabs">
				<Tabs />
			</PrivateRoute>

			<PrivateRoute path="/account">
				<Account />
			</PrivateRoute>

			<Route path="/tab/:id" let:params>
				{#await getTab(params.id) then tab}
					<Tab serverTab="{tab}"/>
				{/await}
			</Route>

			<PrivateRoute path="/edit/:id" let:params>
				{#await getTab(params.id) then tab}
					<TabEditor serverTab="{tab}"/>
				{/await}
			</PrivateRoute>
		</main>
	</Router>

	{#if $error_message !== null}
		<div class="error">
			{$error_message}
		</div>
	{/if}
</div>


<style lang="scss">
	@import "global";

	.error {
		position: fixed;
		padding: 2em;

		right: 2em;
		bottom: 1em;

		background-color: $darkred;
		color: white;
	}

	:global(body, html) {
		margin: 0;
		padding: 0;
	}

	.root {
		display: grid;
		grid-template-rows: 4em auto;
		grid-template-columns: auto;

		height: 100%;
	}

	main {
		background-color: $gray;

		padding: 0;
		margin: 0;

		height: 100%;
	}
</style>
