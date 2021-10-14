const preprocess = require("svelte-preprocess");

const config = {
  preprocess: preprocess({
    defaults: {
      style: 'postcss',
    },
    postcss: true,
  }),
};

module.exports = config;
