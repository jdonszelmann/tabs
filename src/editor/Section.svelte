<script type="ts">
    import Measure from "./Measure.svelte";
    import {SectionData} from "../typescript/SectionData";
    import {TabData} from "../typescript/TabData";
    import {Selection} from "../typescript/Selection";
    import { watchResize } from "svelte-watch-resize";
    import {range} from "../typescript/Range";

    export let sectionIndex: number;
    export let section: SectionData;
    export let tab: TabData;
    export let selection: Selection | null;

    let numLines: number = 1;

    function handleResize(element) {
        const height = element.offsetHeight;

        const names: HTMLElement = document.querySelector(".lines .names");
        if (names === null) {
            numLines = 1;
            return;
        }

        const lineHeight = names.offsetHeight;

        numLines = height / lineHeight;
    }

</script>

<div class="section">

    {#if section.name !== ""}
    <div class="header">
        <h2>{section.name}</h2>
    </div>
    {/if}

    <div class="lines">

        <div class="outer-names">
            {#each range(numLines) as _}
                <div class="names">
                    <div class="inner-names">
                        {#each section.stringNames as name}
                            <span class="name">{name}</span>
                        {/each}
                    </div>
                </div>
            {/each}
        </div>

        <div class="outer-measures">
            <div class="measures" use:watchResize={handleResize}>
                {#each section.measures as measure, measureIndex}
                    <Measure
                            bind:measure="{measure}"
                            bind:section="{section}"
                            bind:tab="{tab}"
                            bind:selection="{selection}"
                            sectionIndex="{sectionIndex}"
                            measureIndex="{measureIndex}"
                    />
                {/each}
            </div>
        </div>
    </div>

</div>

<style lang="scss">
    @import "src/global";

    .section {
      padding: 1em;
      display: flex;
      flex-direction: column;

      .header {
        width: 100%;

        h2 {
          color: $red;
          text-align: center;


          margin-top: .25em;
          padding-bottom: 0;
          margin-bottom: .25em;
        }
      }

      .lines {
        display: flex;
        flex-direction: row;

        .outer-measures {
          display: inline-block;

          .measures {
            display: flex;
            flex-direction: row;

            flex-wrap: wrap;
          }
        }
      }
    }

    .outer-names {
      display: flex;
      flex-direction: column;

      .names {
        display: inline-block;

        .inner-names {
          display: flex;
          flex-direction: column;

          border-right: 2px solid $black;
          padding-bottom: .5em;
          padding-top: .5em;
          padding-right: .5em;

          .name {
            height: 1.25em;
          }
        }
      }
    }
</style>