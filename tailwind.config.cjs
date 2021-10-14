const config = {
  mode: "jit",
  purge: ["./src/**/*.{html,js,svelte,ts}"],

  theme: {
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      black: {
        DEFAULT: "#0A100D"
      },
      brown: {
        DEFAULT: "#B9BAA3",
      },
      gray: {
        DEFAULT: "#D6D5C9",
      },
      red: {
        DEFAULT: "#A22C29",
        dark: "#902923",
      }
    }
  },

  plugins: [],
};

module.exports = config;
