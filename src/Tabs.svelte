<script lang="ts">
    import {report_fetch_error} from "./typescript/Error";
    import {user} from "./typescript/User";
    import {server_url} from "./typescript/Server";
    import {ServerTab} from "./typescript/ServerTab";
    import TabPreview from "./TabPreview.svelte";

    async function getTabs(): Promise<ServerTab[]> {
        const resp = await fetch(`${server_url}/tab/all-for-user`, {
            method: "POST",
            body: JSON.stringify({
                Token: $user.Token
            })
        });

        if (!resp.ok) {
            await report_fetch_error(resp);
        } else {
            const json = await resp.json();

            if (json === null || json === []) {
                return []
            } else {
                return json
            }
        }
    }
</script>

<div class="wrapper">
    <div class="tabs">
        {#await getTabs() then tabs}
            {#if tabs.length === 0}
                <span class="no-tabs">You haven't made any tabs yet</span>
            {:else}
                {#each tabs as tab}
                    <TabPreview tab="{tab}" />
                {/each}
            {/if}
        {/await}
    </div>
</div>

<style lang="scss">
    .wrapper {
      display: flex;
      flex-direction: column;
      align-items: center;

      overflow-y: scroll;
      margin-top: 2em;

      .tabs {
        display: flex;
        flex-direction: column;

        width: calc(min(50em, 100vw));

        .no-tabs {
          margin-top: 2em;
          text-align: center;
        }
      }
    }

</style>