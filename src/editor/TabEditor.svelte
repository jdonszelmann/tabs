
<script lang="ts">
    import {TabData} from "../typescript/TabData";
    import Section from "./Section.svelte";
    import {onMount} from "svelte";
    import {Selection} from "../typescript/Selection";
    import {MeasureData, NoteData} from "../typescript/MeasureData";
    import {Writable, writable} from "svelte/store";
    import {ServerTab} from "../typescript/ServerTab";
    import {user} from "../typescript/User";
    import {useNavigate} from "svelte-navigator";

    const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

    export let serverTab: ServerTab;

    let tab: Writable<TabData>;
    if (serverTab.Contents === "") {
        tab = writable(TabData.loadOrDefault(serverTab.Id));
    } else {
        tab = writable(TabData.fromJSON(JSON.parse(serverTab.Contents)));
    }

    tab.subscribe(async t => {
        if (await t.save()) {
            await t.saveToServer($user.Token)
        }
    })

    let selection: Writable<Selection> = writable(new Selection());
    let numNewBeats: number = $tab.config.startNotesPerMeasure;
    let numBeats: number;
    let numMeasures: number;
    let numSections: number;
    let capo: number;

    let measureClipboard: MeasureData = MeasureData.default($tab.config);

    selection.subscribe(s => {
        numBeats = $tab.sections[s.selectedSection].measures[s.selectedMeasure].beats;
        numMeasures = $tab.sections[s.selectedSection].measures.length;
        numSections = $tab.sections.length;
        capo = $tab.capo;
    })

    let editorFocused = false;

    function setCapo() {
        if (capo < 0) {
            capo = 0;
        }

        $tab.setCapo(capo);
        $tab = $tab;
    }


    function setNumNewBeats() {
        if (numNewBeats <= 0) {
            numNewBeats = 1;
        }

        $tab.config.startNotesPerMeasure = numNewBeats;
    }

    function setBeats() {
        $tab.sections[$selection.selectedSection].
            measures[$selection.selectedMeasure].
            setBeats(numBeats);
        $tab = $tab;
    }

    function setMeasures() {
        $tab.sections[$selection.selectedSection].setMeasures(numMeasures, $tab.config);
        numMeasures = $tab.sections[$selection.selectedSection].measures.length;
        $tab = $tab;
    }

    function setSections() {
        $tab.setSections(numSections, $tab.config);
        numSections = $tab.sections.length;
        $tab = $tab;
    }

    async function handleDeleteMeasure() {
        $tab.sections[$selection.selectedSection].deleteMeasure($selection.selectedMeasure);
        numMeasures -= 1;
        $tab = $tab;
    }

    async function handleKey(k: KeyboardEvent) {
        if (!editorFocused) {
            return
        }

        k.preventDefault();

        setTimeout(() => {
            const note: NoteData = $tab.sections[$selection.selectedSection].
                measures[$selection.selectedMeasure].
                strings[$selection.selectedString].
                notes[$selection.selectedBeat];

            switch (k.key) {
                case "0": note.addNumber(0, $tab.capo); $tab = $tab; break;
                case "1": note.addNumber(1, $tab.capo); $tab = $tab; break;
                case "2": note.addNumber(2, $tab.capo); $tab = $tab; break;
                case "3": note.addNumber(3, $tab.capo); $tab = $tab; break;
                case "4": note.addNumber(4, $tab.capo); $tab = $tab; break;
                case "5": note.addNumber(5, $tab.capo); $tab = $tab; break;
                case "6": note.addNumber(6, $tab.capo); $tab = $tab; break;
                case "7": note.addNumber(7, $tab.capo); $tab = $tab; break;
                case "8": note.addNumber(8, $tab.capo); $tab = $tab; break;
                case "9": note.addNumber(9, $tab.capo); $tab = $tab; break;
                case "-": note.decrement(); $tab = $tab; break;
                case "=": note.increment(); $tab = $tab; break;
                case "Backspace": note.setNull(); $tab = $tab; break;
                case "ArrowUp": $selection.move("up", $tab); $selection = $selection; break;
                case "ArrowDown": $selection.move("down", $tab); $selection = $selection; break;
                case "ArrowLeft": $selection.move("left", $tab); $selection = $selection; break;
                case "ArrowRight": $selection.move("right", $tab); $selection = $selection; break;
                case "Enter": $selection.move("right", $tab); $selection = $selection; break;
                case "Space": $selection.move("right", $tab); $selection = $selection; break;
                case "End": $selection.lastMeasure($tab); $selection = $selection; break;
                case "Home": $selection.firstMeasure($tab); $selection = $selection; break;
                case "Tab": $selection.nextMeasure($tab); $selection = $selection; break;
                case "m": $selection.nextMeasure($tab); $selection = $selection; break;
                case "M": $selection.prevMeasure($tab); $selection = $selection; break;
                case "s": $selection.nextSection($tab); $selection = $selection; break;
                case "S": $selection.prevSection($tab); $selection = $selection; break;

                default:
                    if (k.key == "c" && k.ctrlKey) {
                        measureClipboard = $tab.sections[$selection.selectedSection].measures[$selection.selectedMeasure]
                        return;
                    }
                    if (k.key == "x" && k.ctrlKey) {
                        measureClipboard = $tab.sections[$selection.selectedSection].measures[$selection.selectedMeasure].clone();
                        $tab.sections[$selection.selectedSection].measures[$selection.selectedMeasure] = MeasureData.default($tab.config);
                        return;
                    }
                    if (k.key == "d" && k.ctrlKey) {
                        $tab.sections[$selection.selectedSection].measures[$selection.selectedMeasure] = MeasureData.default($tab.config);
                        return;
                    }
                    if (k.key == "v" && k.ctrlKey) {
                        if (measureClipboard !== null) {
                            $tab.sections[$selection.selectedSection].measures[$selection.selectedMeasure] = measureClipboard.clone();
                        }
                        return;
                    }
                    if (k.key == "r" && k.ctrlKey) {
                        window.location.reload();
                        return;
                    }
                    if (k.key == "z" && k.ctrlKey) {
                        tab.update(t => {
                            return t.undo();
                        })
                        return;
                    }


                    if (alphabet.includes(k.key)) {
                        $tab.sections[$selection.selectedSection].selectStringWithName($selection, k.key)
                        $selection = $selection;
                    }
            }
        }, 10);
    }

    const navigate = useNavigate();

    function viewTab() {
        navigate(`/tab/${$tab.id}`)
    }

    onMount(() => {
        $selection.reset();
    })

</script>

<svelte:window on:keydown={handleKey}/>

<div class="editor">
    <div tabindex="0" class="controls">
        <div class="control tab">
            <h2>Tab</h2>

            <div class="items">
                <label>
                    <button on:click={viewTab}>View Tab</button>
                </label>
                <label>
                    Name
                    <input
                            bind:value={$tab.name}
                            placeholder="no name"
                    >
                </label>
                <label>
                    Sections
                    <input type="number" bind:value={numSections} on:change={setSections}>
                </label>
                <label>
                    Capo
                    <input type="number" bind:value={capo} on:change={setCapo}>
                </label>
            </div>
        </div>
        <div class="control section">
            <h2>Section</h2>
            <div class="items">
                <label>
                    Name
                    <input
                            bind:value={$tab.sections[$selection.selectedSection].name}
                            placeholder="no name"
                    >
                </label>

                <label>
                    Measures
                    <input type="number" bind:value={numMeasures} on:change={setMeasures}>
                </label>
                <label>
                    Beats for new measures
                    <input type="number" bind:value={numNewBeats} on:change={setNumNewBeats}>
                </label>
            </div>
        </div>
        <div class="control measure">
            <h2>Measure</h2>
            <div class="items">
                <label>
                    Beats
                    <input type="number" bind:value={numBeats} on:change={setBeats}>
                </label>
                <label>
                    <button on:click={handleDeleteMeasure}>Delete current Measure</button>
                </label>
            </div>
        </div>

        <div class="location">
            <span>section: {$tab.sections[$selection.selectedSection].name || $selection.selectedSection + 1}</span>
            <span>measure: {$selection.selectedMeasure + 1}</span>
            <span>beat: {$selection.selectedBeat + 1}</span>
            <span>string: {$tab.sections[$selection.selectedSection].stringNames[$selection.selectedString]}</span>
        </div>
    </div>

    <div tabindex="1" class="sections" on:focusin={() => editorFocused = true} on:focusout={() => editorFocused = false}>
        <div class="header">
            <h1>{$tab.name}</h1>
            {#if capo > 0}
                <div>Capo on fret {capo}</div>
            {/if}
        </div>

        {#each $tab.sections as section, sectionIndex}
            <Section
                    bind:section="{section}"
                    bind:tab="{$tab}"
                    bind:selection="{$selection}"
                    sectionIndex="{sectionIndex}"
            />
        {/each}
    </div>
</div>

<style lang="scss">
  @import "src/global";
  .editor {
    display: grid;
    grid-template-columns: minmax(auto, 10em) auto;
    grid-template-rows: auto;

    height: 100%;

    .sections {
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
        margin: 0;

        &:hover {
          cursor: pointer;
          background-color: $red;
        }
      }
    }
  }

</style>