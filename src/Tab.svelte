
<script lang="ts">
    import {ServerTab} from "./typescript/ServerTab";
    import {TabData} from "./typescript/TabData";
    import Section from "./editor/Section.svelte";
    import {user} from "./typescript/User";
    import {useNavigate} from "svelte-navigator";
    import {server_url} from "./typescript/Server";
    import {report_fetch_error} from "./typescript/Error";

    export let serverTab: ServerTab;

    let tab: TabData;
    if (serverTab.Contents === "") {
        tab = TabData.loadOrDefault(serverTab.Id);
    } else {
        tab = TabData.fromJSON(JSON.parse(serverTab.Contents));
    }

    const navigate = useNavigate();

    function editTab() {
        navigate(`/edit/${tab.id}`)
    }

    async function togglePublic() {
        serverTab.Public = !serverTab.Public;
        serverTab = serverTab;

        const resp = await fetch(`${server_url}/tab/public`, {
            method: "PUT",
            body: JSON.stringify({
                Token: $user.Token,
                Id: tab.id,
                Public: serverTab.Public,
            })
        });

        if (!resp.ok) {
            await report_fetch_error(resp);
        }
    }

    async function deleteTab() {
        if (!confirm("Are you sure you want to delete this tab?")) {
            return
        }

        const resp = await fetch(`${server_url}/tab/`, {
            method: "DELETE",
            body: JSON.stringify({
                Token: $user.Token,
                Id: tab.id,
            })
        });

        if (!resp.ok) {
            await report_fetch_error(resp);
        } else {
            navigate("/tabs")
        }
    }
</script>


{#if $user !== null && serverTab.Owner === $user.Name}
<div class="viewer">

    <div tabindex="0" class="controls">
        <div class="control tab">
            <h2>Tab</h2>

            <div class="items">
                <label>
                    <button on:click={editTab}>Edit Tab</button>
                </label>

                <label>
                    <button on:click={deleteTab}>Delete Tab</button>
                </label>

                <label>
                    <button on:click={togglePublic}>Make
                        {#if serverTab.Public}
                            private
                        {:else}
                            public
                        {/if}
                    </button>
                </label>
            </div>
        </div>
    </div>


    <div class="tab">
        <div class="header">
            <h1>{tab.name}</h1>
            {#if tab.capo > 0}
                <div>Capo on fret {tab.capo}</div>
            {/if}
        </div>

        {#each tab.sections as section, sectionIndex}
            <Section
                    bind:section="{section}"
                    bind:tab="{tab}"
                    selection={null}
                    sectionIndex="{sectionIndex}"
            />
        {/each}
    </div>
</div>
{:else}
    <div class="tab">
        <div class="header">
            <h1>{tab.name}</h1>
            {#if tab.capo > 0}
                <div>Capo on fret {tab.capo}</div>
            {/if}
        </div>

        {#each tab.sections as section, sectionIndex}
            <Section
                    bind:section="{section}"
                    bind:tab="{tab}"
                    selection={null}
                    sectionIndex="{sectionIndex}"
            />
        {/each}
    </div>
{/if}

<style lang="scss">
  @import "global";

  .tab {
    display: flex;
    flex-direction: column;
    overflow-y: auto;

    .header {
      h1 {
        color: $red;
      }
      text-align: center;
      padding-bottom: 0;
      margin-bottom: 0;
    }
  }

  .viewer {
    display: grid;
    grid-template-columns: minmax(auto, 10em) auto;
    grid-template-rows: auto;

    height: 100%;

    .controls {
      background-color: $darkred;
      display: flex;
      flex-direction: column;
      color: white;
      padding-left: .5em;
      padding-right: .5em;
      padding-bottom: .5em;
      text-align: center;

      .control {
        h2 {
          text-transform: uppercase;
        }

        .items {
          display: flex;
          flex-direction: column;

          label {
            display: flex;
            flex-direction: column;
          }

          input {
            width: 100%;
            background-color: $red;
            border: 1px solid black;
            color: white;
          }
        }
      }

      .location {
        margin-top: auto;

        display: flex;
        flex-direction: column;
      }

      button {
        color: white;
        background-color: $darkred;
        border: 1px solid black;
        margin: 0 0 .5em;

        &:hover {
          cursor: pointer;
          background-color: $red;
        }
      }
    }
  }


</style>