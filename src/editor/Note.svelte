<script lang="ts">
    import {Selection} from "../typescript/Selection";
    import {NoteData} from "../typescript/MeasureData";
    import {TabData} from "../typescript/TabData";

    export let note: NoteData;
    export let tab: TabData;
    export let measureIndex: number;
    export let sectionIndex: number;
    export let selection: Selection | null;
    export let string: number;
    export let beat: number;


    let noteSelected: boolean = false;
    let stringSelected: boolean = false;

    $: noteSelected = selection !== null && selection.selectedBeat === beat &&
        selection.selectedMeasure === measureIndex &&
        selection.selectedString === string &&
        selection.selectedSection === sectionIndex;

    $: stringSelected = selection !== null && selection.selectedString === string && selection.selectedSection === sectionIndex;

    function clickNote() {
        if (selection === null) {
            return;
        }

        selection.selectedBeat = beat;
        selection.selectedString = string;
        selection.selectedMeasure = measureIndex;
        selection.selectedSection = sectionIndex;
        selection = selection;
    }
</script>

<div
        class="note"
        on:click={clickNote}
        class:selected={noteSelected}
        class:stringSelected={stringSelected && !noteSelected}
>
    {#if note.fretNumber === null}
        <div
                class="line"
        >
            <div class="center"></div>
        </div>
    {:else}
        <div class="fret-number">{note.fretNumber - tab.capo}</div>
    {/if}
</div>

<style lang="scss">
    @import "src/global";

    .note {
      width: 1.25em;
      height: 1.25em;

      .fret-number {
        text-align: center;
        width: 100%;
      }

      .line {
        height: 100%;
        width: 100%;
        display: flex;
        flex-direction: column;
        justify-content: center;

        .center {
          border-bottom: 1px solid $black;
        }

        &:hover {
          background-color: $brown;
        }
      }

      &.selected {
        background-color: $brown;
      }

      &.stringSelected {
        background-color: $brown-a2;
      }
    }
</style>