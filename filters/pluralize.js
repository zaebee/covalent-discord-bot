import Vue from 'vue'

Vue.filter('pluralize', (amount, singular, plural = null) => {
  if (plural === null) {
    plural = `${singular}s`;
  }
  return (amount > 1 || amount === 0) ? plural : singular;
});
