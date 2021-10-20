
<script lang="ts">
    import {MeasureData, NoteData} from "../typescript/MeasureData";
    import Note from "./Note.svelte";
    import {range} from "../typescript/Range";
    import {SectionData} from "../typescript/SectionData";
    import {TabData} from "../typescript/TabData";
    import {Selection} from "../typescript/Selection";

    export let measureIndex: number;
    export let measure: MeasureData;
    export let sectionIndex: number;
    export let section: SectionData;
    export let tab: TabData;
    export let selection: Selection | null;

    function selected(): boolean {
        return selection !== null && selection.selectedMeasure == measureIndex && selection.selectedSection == sectionIndex;
    }

    let columns: string;
    $: columns = "auto ".repeat(measure.beats);

    let rows: string;
    $: rows = "auto ".repeat(measure.strings.length);
</script>

<div
        class="measure"
        style="grid-template-columns: {columns}; grid-template-rows: {rows}"
        class:selected={selection !== null && selection.selectedMeasure===measureIndex && selection.selectedSection === sectionIndex}
>
    {#each range(measure.strings.length) as string}
        {#each range(measure.beats) as beat}
            <Note
                    bind:selection={selection}
                    string="{string}"
                    beat="{beat}"
                    measureIndex="{measureIndex}"
                    sectionIndex="{sectionIndex}"
                    bind:note={measure.strings[string].notes[beat]}
                    bind:tab={tab}
            />
        {/each}
    {/each}
</div>

<style lang="scss">
    @import "src/global";

    .measure {
      display: grid;

      border-right: 2px solid $black;
      padding-bottom: .5em;
      padding-top: .5em;

      cursor: pointer;

      &:hover {
        background-color: $brown-a1;
      }
    }

    .selected {
      background-color: $brown-a1;
    }
</style>