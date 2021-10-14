<script lang="ts">
	import Bar from "./Bar.svelte";

</script>

<main class="bg-gray">
	test
	<Bar />
</main>

<style global lang="postcss">
	@tailwind base;
	@tailwind components;
	@tailwind utilities;

	:global(body, html) {
		margin: 0;
		padding: 0;
	}

	main {
		padding: 0;
		margin: 0;

		width: 100%;
		height: 100%;
	}
</style>
