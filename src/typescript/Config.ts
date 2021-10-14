
export default class Config {
    strings: number;

    constructor(strings: number) {
        this.strings = strings;
    }

    static default(): Config {
        return new Config(
            6,
        );
    }
}