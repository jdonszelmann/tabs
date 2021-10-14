<script type="ts">
    import Measure from "./Measure.svelte";
    import MeasureData from "./typescript/MeasureData";
    import Config from "./typescript/Config";

    let config: Config = Config.default();
    let measures: MeasureData[] = [];


</script>

<div class="flex flex-row">
    {#each measures as measure}
        <div class="bg-red-100">
            <Measure measure="{measure}" config="{config}"/>
        </div>
    {/each}

    <div class="controls bg-gray-100">

    </div>
</div>