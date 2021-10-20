
<script lang="ts">
    import {ServerTab} from "./typescript/ServerTab";
    import {TabData} from "./typescript/TabData";
    import {useNavigate} from "svelte-navigator";

    export let tab: ServerTab;

    let navigate = useNavigate();

    let tabData: TabData = tab.Contents === ""
        ? TabData.default(tab.Id)
        : TabData.fromJSON(JSON.parse(tab.Contents))

    function gotoTab() {
        navigate(`/tab/${tab.Id}`)
    }

</script>

<div class="preview" on:click={gotoTab}>
    <div class="left">
        <h1>{tabData.name}</h1>
        <span>Creator: {tab.Owner}</span>
        <small>{tab.Id}</small>

    </div>

    <div class="right">
        <span class="capo">
            {#if tabData.capo === 0}
                No capo
            {:else}
                Capo: {tabData.capo}
            {/if}
        </span>
        {#if !tab.Public}
            Private
        {/if}
    </div>
</div>

<style lang="scss">
    @import "global";

    .preview {
      width: 100%;
      padding: 1em;
      margin: 1em;
      border-radius: .5em;
      background-color: $brown;

      .capo {
        position: relative;
        right: 0;
      }

      display: grid;
      grid-template-columns: auto auto;


      .left {
        display: flex;
        flex-direction: column;
      }
      .right {
        display: flex;
        flex-direction: column-reverse;

        text-align: right;
      }


      &:hover {
        cursor: pointer;
      }
    }
</style>