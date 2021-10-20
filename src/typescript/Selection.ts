import type {TabData} from "./TabData";

type Direction = "up" | "down" | "left" | "right";

export function wrap(value: number, modulus: number): number {
    if (value < 0) {
        value = modulus - 1
    }  else if (value > modulus - 1) {
        value = 0;
    }

    return value;
}

export class Selection {
    selectedSection: number = 0;
    selectedMeasure: number = 0;
    selectedString: number = 0;
    selectedBeat: number = 0;

    constructor() {
        this.reset();
    }

    firstMeasure(_: TabData) {
        this.selectedMeasure = 0;
    }

    lastMeasure(tab: TabData) {
        this.selectedMeasure = tab.sections[this.selectedSection].measures.length - 1;
    }

    prevSection(tab: TabData) {
        this.selectedSection -= 1;
        if (this.selectedSection < 0) {
            this.selectedSection = tab.sections.length - 1;
        }
    }

    prevMeasure(tab: TabData) {
        this.selectedMeasure -= 1;
        if (this.selectedMeasure < 0) {
            this.prevSection(tab);
            this.selectedMeasure = tab.sections[this.selectedSection].measures.length - 1;
        }
    }

    nextSection(tab: TabData) {
        this.selectedSection += 1;
        if (this.selectedSection > tab.sections.length - 1) {
            this.selectedSection = 0;
        }
    }

    nextMeasure(tab: TabData) {
        this.selectedMeasure += 1;
        if (this.selectedMeasure > tab.sections[this.selectedSection].measures.length - 1) {
            this.nextSection(tab);
            this.selectedMeasure = 0;
        }
    }

    move(direction: Direction, tab: TabData) {
        switch (direction) {
            case "up":
                this.selectedString -= 1;
                this.selectedString = wrap(this.selectedString, tab.sections[this.selectedSection].measures[this.selectedMeasure].strings.length)
                break;
            case "down":
                this.selectedString += 1;
                this.selectedString = wrap(this.selectedString, tab.sections[this.selectedSection].measures[this.selectedMeasure].strings.length)
                break;
            case "left":
                this.selectedBeat -= 1;
                if (this.selectedBeat < 0) {
                    this.prevMeasure(tab);
                    this.selectedBeat = tab.sections[this.selectedSection].measures[this.selectedMeasure].beats - 1;
                }

                break;
            case "right":
                this.selectedBeat += 1;
                if (this.selectedBeat > tab.sections[this.selectedSection].measures[this.selectedMeasure].beats - 1) {
                    this.nextMeasure(tab);
                    this.selectedBeat = 0;
                }

                break;
        }
    }

    reset() {
        this.selectedSection = 0;
        this.selectedString = 0;
        this.selectedMeasure = 0;
        this.selectedBeat = 0;
    }
}